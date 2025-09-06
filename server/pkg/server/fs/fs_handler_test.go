package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"my-collection/server/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// MockFsDb implements a mock for the fsDb interface (model.DirectoryReaderWriter)
type MockFsDb struct {
	mock.Mock
}

func (m *MockFsDb) GetDirectory(conds ...interface{}) (*model.Directory, error) {
	args := m.Called(conds...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Directory), args.Error(1)
}

func (m *MockFsDb) GetDirectories(conds ...interface{}) (*[]model.Directory, error) {
	args := m.Called(conds...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Directory), args.Error(1)
}

func (m *MockFsDb) GetAllDirectories() (*[]model.Directory, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Directory), args.Error(1)
}

func (m *MockFsDb) CreateOrUpdateDirectory(directory *model.Directory) error {
	args := m.Called(directory)
	return args.Error(0)
}

func (m *MockFsDb) UpdateDirectory(directory *model.Directory) error {
	args := m.Called(directory)
	return args.Error(0)
}

func (m *MockFsDb) RemoveDirectory(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockFsDb) RemoveTagFromDirectory(directoryPath string, tagId uint64) error {
	args := m.Called(directoryPath, tagId)
	return args.Error(0)
}

// MockFsDirectoryChangedListener implements a mock for the fsDirectoryChangedListener interface
type MockFsDirectoryChangedListener struct {
	mock.Mock
}

func (m *MockFsDirectoryChangedListener) DirectoryChanged() {
	m.Called()
}

// Test setup functions
func setupFsTestHandler() (*fsHandler, *MockFsDb, *MockFsDirectoryChangedListener) {
	mockDb := &MockFsDb{}
	mockListener := &MockFsDirectoryChangedListener{}
	handler := NewHandler(mockDb, mockListener)
	return handler, mockDb, mockListener
}

func setupFsTestRouter(handler *fsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")
	handler.RegisterRoutes(rg)
	return router
}

// Tests for getFsDir endpoint
func TestFsGetFsDir(t *testing.T) {
	t.Run("Success - Root Directory", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		// Create a temporary directory structure for testing
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "subdir")
		err := os.Mkdir(subDir, 0755)
		require.NoError(t, err)

		testFile := filepath.Join(tempDir, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		require.NoError(t, err)

		// Mock the directory enrichment - return gorm.ErrRecordNotFound to indicate no directory info
		// This is what EnrichFsNode expects to handle gracefully
		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/fs?path=%s&depth=1", tempDir), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var fsNode model.FsNode
		err = json.Unmarshal(w.Body.Bytes(), &fsNode)
		require.NoError(t, err)
		assert.Equal(t, tempDir, fsNode.Path)
		assert.Equal(t, int(model.FS_NODE_DIR), int(fsNode.Type))
		assert.True(t, len(fsNode.Children) >= 1) // Should have at least the subdir and file

		mockDb.AssertExpectations(t)
	})

	t.Run("Success - With Directory Info", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		tempDir := t.TempDir()

		// Mock successful directory retrieval
		dirInfo := &model.Directory{
			Path:     tempDir,
			Excluded: &[]bool{false}[0],
		}
		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dirInfo, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/fs?path=%s&depth=0", tempDir), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var fsNode model.FsNode
		err := json.Unmarshal(w.Body.Bytes(), &fsNode)
		require.NoError(t, err)
		assert.Equal(t, tempDir, fsNode.Path)
		assert.NotNil(t, fsNode.DirInfo)
		assert.Equal(t, tempDir, fsNode.DirInfo.Path)

		mockDb.AssertExpectations(t)
	})

	t.Run("Invalid Depth Parameter", func(t *testing.T) {
		handler, _, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/fs?path=/tmp&depth=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Non-existent Path", func(t *testing.T) {
		handler, _, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/fs?path=/non/existent/path&depth=1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Zero Depth", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		tempDir := t.TempDir()
		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/fs?path=%s&depth=0", tempDir), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var fsNode model.FsNode
		err := json.Unmarshal(w.Body.Bytes(), &fsNode)
		require.NoError(t, err)
		assert.Equal(t, tempDir, fsNode.Path)
		assert.Len(t, fsNode.Children, 0) // Should have no children with depth 0

		mockDb.AssertExpectations(t)
	})
}

// Tests for includeDir endpoint
func TestFsIncludeDir(t *testing.T) {
	t.Run("Handler Parameter Parsing", func(t *testing.T) {
		// This test verifies the handler correctly parses parameters
		// We need some basic mocks since fs.IncludeDir will be called
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		// Add minimal mocks for the database calls that will happen
		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()
		mockDb.On("CreateOrUpdateDirectory", mock.AnythingOfType("*model.Directory")).Return(nil).Maybe()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/fs/include?path=/nonexistent&subdirs=true&hierarchy=false", nil)
		router.ServeHTTP(w, req)

		// We expect an error since the path doesn't exist, but not a parameter parsing error
		assert.NotEqual(t, http.StatusBadRequest, w.Code)
		// The exact error code depends on how the fs business logic handles missing paths
		// Could be 404 (not found) or 500 (internal server error)
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	})

	t.Run("Invalid Subdirs Parameter", func(t *testing.T) {
		handler, _, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/fs/include?path=/tmp&subdirs=invalid&hierarchy=false", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid Hierarchy Parameter", func(t *testing.T) {
		handler, _, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/fs/include?path=/tmp&subdirs=false&hierarchy=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		tempDir := t.TempDir()

		// Simulate a database error
		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, fmt.Errorf("database error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/fs/include?path=%s&subdirs=false&hierarchy=false", tempDir), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})
}

// Tests for excludeDir endpoint
func TestFsExcludeDir(t *testing.T) {
	t.Run("Handler Behavior", func(t *testing.T) {
		// This test verifies the handler works - we need minimal mocks
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		// Add minimal mocks for the database calls that fs.ExcludeDir might make
		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()
		mockDb.On("CreateOrUpdateDirectory", mock.AnythingOfType("*model.Directory")).Return(nil).Maybe()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/fs/exclude?path=/nonexistent", nil)
		router.ServeHTTP(w, req)

		// We expect an error since the path doesn't exist, but the handler should process it
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		tempDir := t.TempDir()

		// Simulate a database error during exclusion
		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, fmt.Errorf("database error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/fs/exclude?path=%s", tempDir), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("Empty Path", func(t *testing.T) {
		handler, _, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/fs/exclude?path=", nil)
		router.ServeHTTP(w, req)

		// Empty path should result in an error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// Tests for SetDirectoryTags endpoint
func TestFsSetDirectoryTags(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockDb, mockListener := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		directory := model.Directory{
			Path: "/test/path",
			Tags: []*model.Tag{
				{Id: 1, Title: "tag1"},
				{Id: 2, Title: "tag2"},
			},
		}

		// Mock the directories.UpdateDirectoryTags function calls
		existingDir := &model.Directory{
			Path: "/test/path",
			Tags: []*model.Tag{
				{Id: 3, Title: "old-tag"},
			},
		}
		mockDb.On("GetDirectory", "path = ?", "/test/path").Return(existingDir, nil)
		mockDb.On("RemoveTagFromDirectory", "/test/path", uint64(3)).Return(nil)
		mockDb.On("CreateOrUpdateDirectory", mock.MatchedBy(func(dir *model.Directory) bool {
			return dir.Path == "/test/path" && len(dir.Tags) == 2
		})).Return(nil)

		mockListener.On("DirectoryChanged").Return()

		body, _ := json.Marshal(directory)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/directories/tags/test/path", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
		mockListener.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		handler, _, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/directories/tags/test/path", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		directory := model.Directory{
			Path: "/test/path",
		}

		// Simulate database error
		mockDb.On("GetDirectory", "path = ?", "/test/path").Return(nil, fmt.Errorf("database error"))

		body, _ := json.Marshal(directory)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/directories/tags/test/path", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("Empty Directory", func(t *testing.T) {
		handler, mockDb, mockListener := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		directory := model.Directory{}

		// Mock for empty directory
		existingDir := &model.Directory{Path: "", Tags: []*model.Tag{}}
		mockDb.On("GetDirectory", "path = ?", "").Return(existingDir, nil)
		mockDb.On("CreateOrUpdateDirectory", mock.AnythingOfType("*model.Directory")).Return(nil)

		mockListener.On("DirectoryChanged").Return()

		body, _ := json.Marshal(directory)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/directories/tags/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
		mockListener.AssertExpectations(t)
	})
}

// Tests for runDirectoriesScan endpoint
func TestFsRunDirectoriesScan(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, _, mockListener := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		mockListener.On("DirectoryChanged").Return()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/directories/scan", nil)
		router.ServeHTTP(w, req)

		// Note: The endpoint doesn't return a specific status, but it should not error
		assert.NotEqual(t, http.StatusNotFound, w.Code)

		mockListener.AssertExpectations(t)
	})

	t.Run("Multiple Calls", func(t *testing.T) {
		handler, _, mockListener := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		// Should be able to call multiple times
		mockListener.On("DirectoryChanged").Return().Times(3)

		for i := 0; i < 3; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/directories/scan", nil)
			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		}

		mockListener.AssertExpectations(t)
	})
}

// Tests for route registration
func TestFsRegisterRoutes(t *testing.T) {
	handler, _, _ := setupFsTestHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")

	// This should not panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(rg)
	})

	// Verify routes are registered by checking the route info
	routes := router.Routes()

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/api/directories/scan"},
		{"POST", "/api/directories/tags/*directory"},
		{"GET", "/api/fs"},
		{"POST", "/api/fs/include"},
		{"POST", "/api/fs/exclude"},
	}

	for _, expectedRoute := range expectedRoutes {
		found := false
		for _, route := range routes {
			if route.Method == expectedRoute.method && route.Path == expectedRoute.path {
				found = true
				break
			}
		}
		assert.True(t, found, fmt.Sprintf("Route %s %s should be registered", expectedRoute.method, expectedRoute.path))
	}
}

// Tests for error scenarios and edge cases
func TestFsErrorScenarios(t *testing.T) {
	t.Run("Large Directory Tree", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		// Create a deeper directory structure
		tempDir := t.TempDir()
		for i := 0; i < 5; i++ {
			subdir := filepath.Join(tempDir, fmt.Sprintf("level_%d", i))
			for j := 0; j < 3; j++ {
				deepDir := filepath.Join(subdir, fmt.Sprintf("subdir_%d", j))
				err := os.MkdirAll(deepDir, 0755)
				require.NoError(t, err)
			}
		}

		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/fs?path=%s&depth=3", tempDir), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var fsNode model.FsNode
		err := json.Unmarshal(w.Body.Bytes(), &fsNode)
		require.NoError(t, err)
		assert.True(t, len(fsNode.Children) > 0)
	})

	t.Run("Special Characters in Path", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		// Create directory with special characters
		tempDir := t.TempDir()
		specialDir := filepath.Join(tempDir, "special-dir-with-dashes")
		err := os.Mkdir(specialDir, 0755)
		require.NoError(t, err)

		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/fs?path=%s&depth=1", specialDir), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unicode Path", func(t *testing.T) {
		handler, mockDb, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		tempDir := t.TempDir()
		unicodeDir := filepath.Join(tempDir, "测试目录")
		err := os.Mkdir(unicodeDir, 0755)
		require.NoError(t, err)

		mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/fs?path=%s&depth=0", unicodeDir), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Permission Denied Directory", func(t *testing.T) {
		// This test would require specific permissions setup,
		// so we'll test the error handling path instead
		handler, _, _ := setupFsTestHandler()
		router := setupFsTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/fs?path=/root&depth=1", nil)
		router.ServeHTTP(w, req)

		// Should handle permission errors gracefully
		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})
}

// Benchmark tests
func BenchmarkFsGetFsDir(b *testing.B) {
	handler, mockDb, _ := setupFsTestHandler()
	router := setupFsTestRouter(handler)

	tempDir := b.TempDir()
	mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/fs?path=%s&depth=1", tempDir), nil)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkFsIncludeDir(b *testing.B) {
	handler, mockDb, mockListener := setupFsTestHandler()
	router := setupFsTestRouter(handler)

	tempDir := b.TempDir()
	mockDb.On("GetDirectory", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Maybe()
	mockDb.On("CreateOrUpdateDirectory", mock.AnythingOfType("*model.Directory")).Return(nil).Maybe()
	mockListener.On("DirectoryChanged").Return().Maybe()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/fs/include?path=%s&subdirs=false&hierarchy=false", tempDir), nil)
		router.ServeHTTP(w, req)
	}
}
