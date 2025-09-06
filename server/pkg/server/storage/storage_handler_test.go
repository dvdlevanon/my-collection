package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"my-collection/server/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStorageHandlerStorage implements a mock for the storageHandlerStorage interface
type MockStorageHandlerStorage struct {
	mock.Mock
}

func (m *MockStorageHandlerStorage) IsStorageUrl(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *MockStorageHandlerStorage) GetFile(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *MockStorageHandlerStorage) GetStorageUrl(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *MockStorageHandlerStorage) GetFileForWriting(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func (m *MockStorageHandlerStorage) GetTempFile() string {
	args := m.Called()
	return args.String(0)
}

// Test setup functions
func setupStorageTestHandler() (*storageHandler, *MockStorageHandlerStorage) {
	mockStorage := &MockStorageHandlerStorage{}
	handler := NewHandler(mockStorage)
	return handler, mockStorage
}

func setupStorageTestRouter(handler *storageHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")
	handler.RegisterRoutes(rg)
	return router
}

// Helper function to create multipart form data for file upload
func createMultipartForm(fieldName, fileName, filePath, content string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the path field
	_ = writer.WriteField("path", filePath)

	// Add the file field
	part, _ := writer.CreateFormFile("file", fileName)
	_, _ = part.Write([]byte(content))

	_ = writer.Close()
	return body, writer.FormDataContentType()
}

// Tests for uploadFile endpoint
func TestStorageUploadFile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "test-file-uuid.jpg")
		mockStorage.On("GetFileForWriting", mock.MatchedBy(func(relativeFile string) bool {
			return strings.Contains(relativeFile, "test-file.jpg") && strings.Contains(relativeFile, "uploads/")
		})).Return(mockFile, nil)
		mockStorage.On("GetStorageUrl", mock.AnythingOfType("string")).Return("http://storage.example.com/uploads/test-file-uuid.jpg")

		// Create test file content
		body, contentType := createMultipartForm("file", "test-file.jpg", "uploads/", "test file content")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", contentType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.FileUrl
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "http://storage.example.com/uploads/test-file-uuid.jpg", response.Url)

		mockStorage.AssertExpectations(t)
	})

	t.Run("Storage GetFileForWriting Error", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		mockStorage.On("GetFileForWriting", mock.AnythingOfType("string")).Return("", fmt.Errorf("storage error"))

		body, contentType := createMultipartForm("file", "test-file.jpg", "uploads/", "test content")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", contentType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Invalid Multipart Form", func(t *testing.T) {
		handler, _ := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", bytes.NewBufferString("invalid multipart data"))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=invalid")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Missing File Field", func(t *testing.T) {
		// Note: This test demonstrates that the handler will panic when file field is missing
		// This is a limitation of the current handler implementation
		handler, _ := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("path", "uploads/")
		// Intentionally not adding file field
		_ = writer.Close()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// This will panic due to handler implementation, so we test that it doesn't return 200
		assert.Panics(t, func() {
			router.ServeHTTP(w, req)
		})
	})

	t.Run("Missing Path Field", func(t *testing.T) {
		// Note: This test demonstrates that the handler will panic when path field is missing
		// This is a limitation of the current handler implementation
		handler, _ := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.jpg")
		_, _ = part.Write([]byte("content"))
		// Intentionally not adding path field
		_ = writer.Close()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// This will panic due to handler implementation
		assert.Panics(t, func() {
			router.ServeHTTP(w, req)
		})
	})

	t.Run("File Save Error", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		// Use an invalid directory that doesn't exist and can't be created
		invalidPath := "/invalid/readonly/path/test-file.jpg"
		mockStorage.On("GetFileForWriting", mock.AnythingOfType("string")).Return(invalidPath, nil)

		body, contentType := createMultipartForm("file", "test-file.jpg", "uploads/", "test content")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", contentType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockStorage.AssertExpectations(t)
	})
}

// Tests for uploadFileFromUrl endpoint
func TestStorageUploadFileFromUrl(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		// Create a test HTTP server to serve the file
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			_, _ = w.Write([]byte("fake image content"))
		}))
		defer testServer.Close()

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "downloaded-file.jpg")
		mockStorage.On("GetFileForWriting", "uploads/image.jpg").Return(mockFile, nil)
		mockStorage.On("GetStorageUrl", "uploads/image.jpg").Return("http://storage.example.com/uploads/image.jpg")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/upload-file-from-url?url=%s&path=%s", testServer.URL, "uploads/image.jpg"), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.FileUrl
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "http://storage.example.com/uploads/image.jpg", response.Url)

		// Verify file was actually created
		_, err = os.Stat(mockFile)
		assert.NoError(t, err)

		mockStorage.AssertExpectations(t)
	})

	t.Run("Storage GetFileForWriting Error", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		mockStorage.On("GetFileForWriting", "uploads/image.jpg").Return("", fmt.Errorf("storage error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file-from-url?url=http://example.com/image.jpg&path=uploads/image.jpg", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "file.jpg")
		mockStorage.On("GetFileForWriting", "uploads/image.jpg").Return(mockFile, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file-from-url?url=invalid-url&path=uploads/image.jpg", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Download Error - Server Not Found", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "file.jpg")
		mockStorage.On("GetFileForWriting", "uploads/image.jpg").Return(mockFile, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file-from-url?url=http://nonexistent.example.com/image.jpg&path=uploads/image.jpg", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Missing URL Parameter", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "file.jpg")
		mockStorage.On("GetFileForWriting", "uploads/image.jpg").Return(mockFile, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file-from-url?path=uploads/image.jpg", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Missing Path Parameter", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "file.jpg")
		mockStorage.On("GetFileForWriting", "").Return(mockFile, nil)
		mockStorage.On("GetStorageUrl", "").Return("http://storage.example.com/")

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("content"))
		}))
		defer testServer.Close()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/upload-file-from-url?url=%s", testServer.URL), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code) // Should still work with empty path
		mockStorage.AssertExpectations(t)
	})
}

// Tests for getFile endpoint
func TestStorageGetFile(t *testing.T) {
	t.Run("Success - Storage URL", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		// Create a temporary file to serve
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		testContent := "Hello, World!"
		err := os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(t, err)

		mockStorage.On("IsStorageUrl", "uploads/test.txt").Return(true)
		mockStorage.On("GetFile", "uploads/test.txt").Return(testFile)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/file/uploads/test.txt", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, testContent, w.Body.String())

		mockStorage.AssertExpectations(t)
	})

	t.Run("Success - Regular Path with Relativasor", func(t *testing.T) {
		// This test verifies that regular paths (non-storage URLs) are processed through relativasor
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		// Create a test file in the current working directory (which is the server root in tests)
		testFile := "test-regular-file.txt"
		testContent := "Regular file content via relativasor"

		// Create the file in the current directory
		err := os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(t, err)
		defer os.Remove(testFile) // Clean up after test

		// The handler will call relativasor.GetAbsoluteFile() for this path
		mockStorage.On("IsStorageUrl", testFile).Return(false)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/file/"+testFile, nil)
		router.ServeHTTP(w, req)

		// The file should be found and served
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, testContent, w.Body.String())

		mockStorage.AssertExpectations(t)
	})

	t.Run("File Not Found - Regular Path with Relativasor", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		// Test with a non-existent file that will be processed by relativasor
		nonExistentFile := "non-existent-regular-file.txt"
		mockStorage.On("IsStorageUrl", nonExistentFile).Return(false)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/file/"+nonExistentFile, nil)
		router.ServeHTTP(w, req)

		// Should return 404 when file doesn't exist in the root directory
		assert.Equal(t, http.StatusNotFound, w.Code)

		mockStorage.AssertExpectations(t)
	})

	t.Run("File Not Found - Storage URL", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		nonExistentFile := "/path/to/nonexistent/file.txt"
		mockStorage.On("IsStorageUrl", "uploads/nonexistent.txt").Return(true)
		mockStorage.On("GetFile", "uploads/nonexistent.txt").Return(nonExistentFile)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/file/uploads/nonexistent.txt", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockStorage.AssertExpectations(t)
	})

	t.Run("Empty Path", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		mockStorage.On("IsStorageUrl", "").Return(false)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/file/", nil)
		router.ServeHTTP(w, req)

		// Should handle empty path gracefully
		assert.NotEqual(t, http.StatusMethodNotAllowed, w.Code)

		mockStorage.AssertExpectations(t)
	})

	t.Run("Deep Nested Path - Storage URL", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		deepPath := "uploads/2024/01/subfolder/file.jpg"
		mockStorage.On("IsStorageUrl", deepPath).Return(true)
		mockStorage.On("GetFile", deepPath).Return("/tmp/file.jpg")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/file/uploads/2024/01/subfolder/file.jpg", nil)
		router.ServeHTTP(w, req)

		// The exact response depends on whether the file exists, but the route should work
		assert.NotEqual(t, http.StatusMethodNotAllowed, w.Code)

		mockStorage.AssertExpectations(t)
	})

	t.Run("Deep Nested Path - Regular Path with Relativasor", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		// Create a nested directory structure for testing relativasor behavior
		nestedPath := "testdir/subdir/nested-file.txt"
		err := os.MkdirAll("testdir/subdir", 0755)
		require.NoError(t, err)
		defer os.RemoveAll("testdir") // Clean up after test

		testContent := "Nested file content"
		err = os.WriteFile(nestedPath, []byte(testContent), 0644)
		require.NoError(t, err)

		mockStorage.On("IsStorageUrl", nestedPath).Return(false)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/file/"+nestedPath, nil)
		router.ServeHTTP(w, req)

		// Should find and serve the nested file via relativasor
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, testContent, w.Body.String())

		mockStorage.AssertExpectations(t)
	})
}

// Tests for route registration
func TestStorageRegisterRoutes(t *testing.T) {
	handler, _ := setupStorageTestHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")

	// This should not panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(rg)
	})

	// We can verify the routes exist by checking the route info
	// This is safer than actually making HTTP requests which require complex mocking
	routes := router.Routes()

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/file/*path"},
		{"POST", "/api/upload-file"},
		{"POST", "/api/upload-file-from-url"},
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
func TestStorageErrorScenarios(t *testing.T) {
	t.Run("Large File Upload", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "large-file.bin")
		mockStorage.On("GetFileForWriting", mock.AnythingOfType("string")).Return(mockFile, nil)
		mockStorage.On("GetStorageUrl", mock.AnythingOfType("string")).Return("http://storage.example.com/large-file.bin")

		// Create a large content string (1MB)
		largeContent := strings.Repeat("A", 1024*1024)
		body, contentType := createMultipartForm("file", "large-file.bin", "uploads/", largeContent)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", contentType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Special Characters in Filename", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "special-file.jpg")
		mockStorage.On("GetFileForWriting", mock.MatchedBy(func(relativeFile string) bool {
			return strings.Contains(relativeFile, "special file & symbols.jpg")
		})).Return(mockFile, nil)
		mockStorage.On("GetStorageUrl", mock.AnythingOfType("string")).Return("http://storage.example.com/special-file.jpg")

		body, contentType := createMultipartForm("file", "special file & symbols.jpg", "uploads/", "content")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", contentType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Empty File Upload", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "empty-file.txt")
		mockStorage.On("GetFileForWriting", mock.AnythingOfType("string")).Return(mockFile, nil)
		mockStorage.On("GetStorageUrl", mock.AnythingOfType("string")).Return("http://storage.example.com/empty-file.txt")

		body, contentType := createMultipartForm("file", "empty-file.txt", "uploads/", "")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", contentType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Unicode Filename", func(t *testing.T) {
		handler, mockStorage := setupStorageTestHandler()
		router := setupStorageTestRouter(handler)

		tempDir := t.TempDir()
		mockFile := filepath.Join(tempDir, "unicode-file.txt")
		mockStorage.On("GetFileForWriting", mock.MatchedBy(func(relativeFile string) bool {
			return strings.Contains(relativeFile, "测试文件.txt")
		})).Return(mockFile, nil)
		mockStorage.On("GetStorageUrl", mock.AnythingOfType("string")).Return("http://storage.example.com/unicode-file.txt")

		body, contentType := createMultipartForm("file", "测试文件.txt", "uploads/", "Unicode content")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", body)
		req.Header.Set("Content-Type", contentType)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockStorage.AssertExpectations(t)
	})
}

// Benchmark tests
func BenchmarkStorageUploadFile(b *testing.B) {
	handler, mockStorage := setupStorageTestHandler()
	router := setupStorageTestRouter(handler)

	tempDir := b.TempDir()
	mockFile := filepath.Join(tempDir, "bench-file.txt")
	mockStorage.On("GetFileForWriting", mock.AnythingOfType("string")).Return(mockFile, nil)
	mockStorage.On("GetStorageUrl", mock.AnythingOfType("string")).Return("http://storage.example.com/bench-file.txt")

	body, contentType := createMultipartForm("file", "bench.txt", "uploads/", "benchmark content")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/upload-file", bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", contentType)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkStorageGetFile(b *testing.B) {
	handler, mockStorage := setupStorageTestHandler()
	router := setupStorageTestRouter(handler)

	// Create a test file
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "bench-get.txt")
	_ = os.WriteFile(testFile, []byte("benchmark content"), 0644)

	mockStorage.On("IsStorageUrl", "bench-get.txt").Return(true)
	mockStorage.On("GetFile", "bench-get.txt").Return(testFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/file/bench-get.txt", nil)
		router.ServeHTTP(w, req)
	}
}
