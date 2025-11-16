package management

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"my-collection/server/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockManagementDb implements a mock for the managementDb interface
type MockManagementDb struct {
	mock.Mock
}

func (m *MockManagementDb) GetAllItems(ctx context.Context) (*[]model.Item, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Item), args.Error(1)
}

func (m *MockManagementDb) GetItem(ctx context.Context, conds ...interface{}) (*model.Item, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Item), args.Error(1)
}

func (m *MockManagementDb) GetItems(ctx context.Context, conds ...interface{}) (*[]model.Item, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Item), args.Error(1)
}

func (m *MockManagementDb) GetItemsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockManagementDb) GetTotalDurationSeconds(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockManagementDb) GetAllTags(ctx context.Context) (*[]model.Tag, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

func (m *MockManagementDb) GetTag(ctx context.Context, conds ...interface{}) (*model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockManagementDb) GetTags(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

func (m *MockManagementDb) GetTagsWithoutChildren(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

func (m *MockManagementDb) GetTagsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockManagementProcessor implements a mock for the managementProcessor interface
type MockManagementProcessor struct {
	mock.Mock
}

func (m *MockManagementProcessor) EnqueueAllItemsCovers(ctx context.Context, force bool) error {
	args := m.Called(ctx, force)
	return args.Error(0)
}

func (m *MockManagementProcessor) EnqueueAllItemsFileMetadata(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockManagementProcessor) EnqueueAllItemsPreview(ctx context.Context, force bool) error {
	args := m.Called(ctx, force)
	return args.Error(0)
}

func (m *MockManagementProcessor) EnqueueAllItemsVideoMetadata(ctx context.Context, force bool) error {
	args := m.Called(ctx, force)
	return args.Error(0)
}

func (m *MockManagementProcessor) GenerateMixOnDemand(ctx context.Context, ctg model.CurrentTimeGetter, desc string, tags []model.Tag) (*model.Tag, error) {
	args := m.Called(ctx, ctg, desc, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockManagementProcessor) EnqueueItemOptimizer() {
	m.Called()
}

func (m *MockManagementProcessor) EnqueueSpecTagger() {
	m.Called()
}

// Test setup functions
func setupManagementTestHandler() (*managementHandler, *MockManagementDb, *MockManagementProcessor) {
	mockDb := &MockManagementDb{}
	mockProcessor := &MockManagementProcessor{}

	handler := NewHandler(mockDb, mockProcessor)
	return handler, mockDb, mockProcessor
}

func setupManagementTestRouter(handler *managementHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")
	handler.RegisterRoutes(rg)
	return router
}

// Tests for stats endpoint
func TestManagementGetStats(t *testing.T) {
	handler, mockDb, _ := setupManagementTestHandler()
	router := setupManagementTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		mockDb.On("GetItemsCount", mock.Anything).Return(int64(100), nil)
		mockDb.On("GetTagsCount", mock.Anything).Return(int64(50), nil)
		mockDb.On("GetTotalDurationSeconds", mock.Anything).Return(3600.5, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/stats", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var stats model.Stats
		err := json.Unmarshal(w.Body.Bytes(), &stats)
		require.NoError(t, err)

		assert.Equal(t, int64(100), stats.ItemsCount)
		assert.Equal(t, int64(50), stats.TagsCount)
		assert.Equal(t, 3600.5, stats.TotalDurationSeconds)

		mockDb.AssertExpectations(t)
	})

	t.Run("Database Error - Items Count", func(t *testing.T) {
		handler, mockDb, _ := setupManagementTestHandler()
		router := setupManagementTestRouter(handler)

		mockDb.On("GetItemsCount", mock.Anything).Return(int64(0), assert.AnError)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/stats", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockDb.AssertExpectations(t)
	})
}

// Tests for refresh operations
func TestManagementRefreshOperations(t *testing.T) {
	t.Run("Refresh Items Covers", func(t *testing.T) {
		handler, _, mockProcessor := setupManagementTestHandler()
		router := setupManagementTestRouter(handler)

		t.Run("With Force=true", func(t *testing.T) {
			mockProcessor.On("EnqueueAllItemsCovers", mock.Anything, true).Return(nil)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/items/refresh-covers?force=true", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockProcessor.AssertExpectations(t)
		})

		t.Run("With Force=false", func(t *testing.T) {
			handler, _, mockProcessor := setupManagementTestHandler()
			router := setupManagementTestRouter(handler)

			mockProcessor.On("EnqueueAllItemsCovers", mock.Anything, false).Return(nil)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/items/refresh-covers?force=false", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockProcessor.AssertExpectations(t)
		})

		t.Run("Default Force=false", func(t *testing.T) {
			handler, _, mockProcessor := setupManagementTestHandler()
			router := setupManagementTestRouter(handler)

			mockProcessor.On("EnqueueAllItemsCovers", mock.Anything, false).Return(nil)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/items/refresh-covers", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockProcessor.AssertExpectations(t)
		})

		t.Run("Processor Error", func(t *testing.T) {
			handler, _, mockProcessor := setupManagementTestHandler()
			router := setupManagementTestRouter(handler)

			mockProcessor.On("EnqueueAllItemsCovers", mock.Anything, false).Return(assert.AnError)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/items/refresh-covers", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			mockProcessor.AssertExpectations(t)
		})
	})

	t.Run("Refresh Items Preview", func(t *testing.T) {
		handler, _, mockProcessor := setupManagementTestHandler()
		router := setupManagementTestRouter(handler)

		mockProcessor.On("EnqueueAllItemsPreview", mock.Anything, true).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/refresh-preview?force=true", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("Refresh Items Video Metadata", func(t *testing.T) {
		handler, _, mockProcessor := setupManagementTestHandler()
		router := setupManagementTestRouter(handler)

		mockProcessor.On("EnqueueAllItemsVideoMetadata", mock.Anything, false).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/refresh-video-metadata", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("Refresh Items File Metadata", func(t *testing.T) {
		handler, _, mockProcessor := setupManagementTestHandler()
		router := setupManagementTestRouter(handler)

		mockProcessor.On("EnqueueAllItemsFileMetadata", mock.Anything).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/refresh-file-metadata", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockProcessor.AssertExpectations(t)
	})
}

// Tests for processor operations
func TestManagementProcessorOperations(t *testing.T) {
	t.Run("Run Spec Tagger", func(t *testing.T) {
		handler, _, mockProcessor := setupManagementTestHandler()
		router := setupManagementTestRouter(handler)

		mockProcessor.On("EnqueueSpecTagger").Return()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/spectagger/run", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("Run Items Optimizer", func(t *testing.T) {
		handler, _, mockProcessor := setupManagementTestHandler()
		router := setupManagementTestRouter(handler)

		mockProcessor.On("EnqueueItemOptimizer").Return()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/itemsoptimizer/run", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockProcessor.AssertExpectations(t)
	})
}

// Tests for mix on demand
func TestManagementGenerateMixOnDemand(t *testing.T) {
	handler, _, mockProcessor := setupManagementTestHandler()
	router := setupManagementTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		tags := []model.Tag{
			{Id: 1, Title: "rock"},
			{Id: 2, Title: "classical"},
		}

		expectedTag := &model.Tag{Id: 123, Title: "generated-mix"}
		mockProcessor.On("GenerateMixOnDemand", mock.Anything, mock.AnythingOfType("utils.NowTimeGetter"), "My Mix", tags).Return(expectedTag, nil)

		body, _ := json.Marshal(tags)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/mix-on-demand?desc=My Mix", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resultTag model.Tag
		err := json.Unmarshal(w.Body.Bytes(), &resultTag)
		require.NoError(t, err)
		assert.Equal(t, expectedTag.Id, resultTag.Id)

		mockProcessor.AssertExpectations(t)
	})

	t.Run("Missing Description", func(t *testing.T) {
		tags := []model.Tag{{Id: 1, Title: "rock"}}
		body, _ := json.Marshal(tags)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/mix-on-demand", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/mix-on-demand?desc=Test", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Processor Error", func(t *testing.T) {
		tags := []model.Tag{{Id: 1, Title: "rock"}}
		mockProcessor.On("GenerateMixOnDemand", mock.Anything, mock.AnythingOfType("utils.NowTimeGetter"), "My Mix", tags).Return(nil, assert.AnError)

		body, _ := json.Marshal(tags)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/mix-on-demand?desc=My Mix", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockProcessor.AssertExpectations(t)
	})
}

// Tests for export metadata
func TestManagementExportMetadata(t *testing.T) {
	handler, mockDb, _ := setupManagementTestHandler()
	router := setupManagementTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		// Mock the backup.Export function by setting up the required data
		mockDb.On("GetAllItems", mock.Anything).Return(&[]model.Item{{Id: 1, Title: "test.mp4"}}, nil)
		mockDb.On("GetAllTags", mock.Anything).Return(&[]model.Tag{{Id: 1, Title: "test-tag"}}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/export-metadata.json", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, "gallery-metadata.json", w.Header().Get("Content-Disposition"))
		assert.NotEmpty(t, w.Body.String())

		mockDb.AssertExpectations(t)
	})
}

// Test route registration
func TestManagementRegisterRoutes(t *testing.T) {
	handler, _, _ := setupManagementTestHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")

	// This should not panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(rg)
	})
}
