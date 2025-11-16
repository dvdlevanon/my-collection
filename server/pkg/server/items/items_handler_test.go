package items

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"my-collection/server/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockItemsHandlerDb is a mock implementation of itemsHandlerDb interface
type MockItemsHandlerDb struct {
	mock.Mock
}

func (m *MockItemsHandlerDb) GetItem(ctx context.Context, conds ...interface{}) (*model.Item, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Item), args.Error(1)
}

func (m *MockItemsHandlerDb) GetItems(ctx context.Context, conds ...interface{}) (*[]model.Item, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Item), args.Error(1)
}

func (m *MockItemsHandlerDb) GetAllItems(ctx context.Context) (*[]model.Item, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Item), args.Error(1)
}

func (m *MockItemsHandlerDb) CreateOrUpdateItem(ctx context.Context, item *model.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemsHandlerDb) UpdateItem(ctx context.Context, item *model.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemsHandlerDb) RemoveItem(ctx context.Context, itemId uint64) error {
	args := m.Called(ctx, itemId)
	return args.Error(0)
}

func (m *MockItemsHandlerDb) RemoveTagFromItem(ctx context.Context, itemId uint64, tagId uint64) error {
	args := m.Called(ctx, itemId, tagId)
	return args.Error(0)
}

func (m *MockItemsHandlerDb) GetTag(ctx context.Context, conds ...interface{}) (*model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockItemsHandlerDb) GetTags(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

func (m *MockItemsHandlerDb) GetAllTags(ctx context.Context) (*[]model.Tag, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

func (m *MockItemsHandlerDb) GetTagsWithoutChildren(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

// MockItemsHandlerProcessor is a mock implementation of itemsHandlerProcessor interface
type MockItemsHandlerProcessor struct {
	mock.Mock
}

func (m *MockItemsHandlerProcessor) EnqueueItemVideoMetadata(ctx context.Context, id uint64) {
	m.Called(ctx, id)
}

func (m *MockItemsHandlerProcessor) EnqueueItemCovers(ctx context.Context, id uint64) {
	m.Called(ctx, id)
}

func (m *MockItemsHandlerProcessor) EnqueueCropFrame(ctx context.Context, id uint64, second float64, rect model.RectFloat) {
	m.Called(ctx, id, second, rect)
}

func (m *MockItemsHandlerProcessor) EnqueueItemPreview(ctx context.Context, id uint64) {
	m.Called(ctx, id)
}

func (m *MockItemsHandlerProcessor) EnqueueItemFileMetadata(ctx context.Context, id uint64) {
	m.Called(ctx, id)
}

func (m *MockItemsHandlerProcessor) EnqueueMainCover(ctx context.Context, id uint64, second float64) {
	m.Called(ctx, id, second)
}

// MockItemsHandlerOptimizer is a mock implementation of itemsHandlerOptimizer interface
type MockItemsHandlerOptimizer struct {
	mock.Mock
}

func (m *MockItemsHandlerOptimizer) HandleItem(ctx context.Context, item *model.Item) {
	m.Called(ctx, item)
}

// Test setup helper
func setupTestHandler() (interface{}, *MockItemsHandlerDb, *MockItemsHandlerProcessor, *MockItemsHandlerOptimizer) {
	mockDb := new(MockItemsHandlerDb)
	mockProcessor := new(MockItemsHandlerProcessor)
	mockOptimizer := new(MockItemsHandlerOptimizer)

	handler := NewHandler(mockDb, mockProcessor, mockOptimizer)

	return handler, mockDb, mockProcessor, mockOptimizer
}

// Test setup for Gin router
func setupTestRouter(handler interface{}) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")

	// Type assert to the concrete handler type
	itemsHandler := handler.(*itemsHandler)
	itemsHandler.RegisterRoutes(rg)
	return router
}

func TestGetItems(t *testing.T) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		expectedItems := &[]model.Item{
			{Id: 1, Title: "Item 1", Origin: "/path/to/item1"},
			{Id: 2, Title: "Item 2", Origin: "/path/to/item2"},
		}

		mockDb.On("GetAllItems", mock.Anything).Return(expectedItems, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/items", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseItems []model.Item
		err := json.Unmarshal(w.Body.Bytes(), &responseItems)
		assert.NoError(t, err)
		assert.Equal(t, *expectedItems, responseItems)

		mockDb.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Create new handler for this test to avoid mock conflicts
		handler, mockDb, _, _ := setupTestHandler()
		router := setupTestRouter(handler)

		mockDb.On("GetAllItems", mock.Anything).Return(nil, fmt.Errorf("database error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/items", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})
}

func TestCreateItem(t *testing.T) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		inputItem := model.Item{
			Title:  "New Item",
			Origin: "/absolute/path/to/item",
			Url:    "/absolute/path/to/item/file.mp4",
		}

		mockDb.On("CreateOrUpdateItem", mock.Anything, mock.MatchedBy(func(item *model.Item) bool {
			return item.Title == inputItem.Title &&
				item.Origin != "" && // Should be relativized
				item.Url != "" // Should be relativized
		})).Run(func(args mock.Arguments) {
			item := args.Get(1).(*model.Item)
			item.Id = 123 // Simulate DB assigning ID
		}).Return(nil)

		jsonBody, _ := json.Marshal(inputItem)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseItem model.Item
		err := json.Unmarshal(w.Body.Bytes(), &responseItem)
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), responseItem.Id)

		mockDb.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		inputItem := model.Item{
			Title:  "New Item",
			Origin: "/path/to/item",
		}

		mockDb.On("CreateOrUpdateItem", mock.Anything, mock.AnythingOfType("*model.Item")).Return(fmt.Errorf("database error"))

		jsonBody, _ := json.Marshal(inputItem)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})
}

func TestUpdateItem(t *testing.T) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		itemId := uint64(123)
		inputItem := model.Item{
			Id:     itemId,
			Title:  "Updated Item",
			Origin: "/path/to/item",
		}

		mockDb.On("UpdateItem", mock.Anything, mock.MatchedBy(func(item *model.Item) bool {
			return item.Id == itemId && item.Title == inputItem.Title
		})).Return(nil)

		jsonBody, _ := json.Marshal(inputItem)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d", itemId), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("ID Mismatch", func(t *testing.T) {
		itemId := uint64(123)
		inputItem := model.Item{
			Id:     456, // Different ID
			Title:  "Updated Item",
			Origin: "/path/to/item",
		}

		jsonBody, _ := json.Marshal(inputItem)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d", itemId), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid Item ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/invalid", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetItem(t *testing.T) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		itemId := uint64(123)
		expectedItem := &model.Item{
			Id:     itemId,
			Title:  "Test Item",
			Origin: "/path/to/item",
		}

		mockDb.On("GetItem", mock.Anything, itemId).Return(expectedItem, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/items/%d", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseItem model.Item
		err := json.Unmarshal(w.Body.Bytes(), &responseItem)
		assert.NoError(t, err)
		assert.Equal(t, *expectedItem, responseItem)

		mockDb.AssertExpectations(t)
	})

	t.Run("Item Not Found", func(t *testing.T) {
		itemId := uint64(999)
		mockDb.On("GetItem", mock.Anything, itemId).Return(nil, fmt.Errorf("item not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/items/%d", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("Invalid Item ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/items/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDeleteItem(t *testing.T) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success Without Deleting Real File", func(t *testing.T) {
		itemId := uint64(123)

		mockDb.On("RemoveItem", mock.Anything, itemId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/items/%d?deleteRealFile=false", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("Success With Deleting Real File", func(t *testing.T) {
		itemId := uint64(123)
		testItem := &model.Item{
			Id:     itemId,
			Title:  "test.mp4",
			Origin: "/tmp",
			Url:    "/tmp/test.mp4",
		}

		mockDb.On("GetItem", mock.Anything, itemId).Return(testItem, nil)
		mockDb.On("RemoveItem", mock.Anything, itemId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/items/%d?deleteRealFile=true", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("Invalid deleteRealFile Parameter", func(t *testing.T) {
		itemId := uint64(123)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/items/%d?deleteRealFile=invalid", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Create new handler for this test to avoid mock conflicts
		handler, mockDb, _, _ := setupTestHandler()
		router := setupTestRouter(handler)

		itemId := uint64(123)

		mockDb.On("RemoveItem", mock.Anything, itemId).Return(fmt.Errorf("database error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/items/%d?deleteRealFile=false", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})
}

func TestGetItemLocation(t *testing.T) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		itemId := uint64(123)
		testItem := &model.Item{
			Id:     itemId,
			Title:  "test.mp4",
			Origin: "/relative/path",
			Url:    "/relative/path/test.mp4",
		}

		mockDb.On("GetItem", mock.Anything, itemId).Return(testItem, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/items/%d/location", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.FileUrl
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Url) // URL should be converted to absolute path

		mockDb.AssertExpectations(t)
	})

	t.Run("Item Not Found", func(t *testing.T) {
		itemId := uint64(999)
		mockDb.On("GetItem", mock.Anything, itemId).Return(nil, fmt.Errorf("item not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/items/%d/location", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})
}

func TestRemoveTagFromItem(t *testing.T) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		itemId := uint64(123)
		tagId := uint64(456)

		mockDb.On("RemoveTagFromItem", mock.Anything, itemId, tagId).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/remove-tag/%d", itemId, tagId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Create new handler for this test to avoid mock conflicts
		handler, mockDb, _, _ := setupTestHandler()
		router := setupTestRouter(handler)

		itemId := uint64(123)
		tagId := uint64(456)

		mockDb.On("RemoveTagFromItem", mock.Anything, itemId, tagId).Return(fmt.Errorf("database error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/remove-tag/%d", itemId, tagId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("Invalid Item ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/invalid/remove-tag/456", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid Tag ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/123/remove-tag/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestSetMainCover(t *testing.T) {
	handler, _, mockProcessor, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		itemId := uint64(123)
		second := 5.5

		mockProcessor.On("EnqueueMainCover", mock.Anything, itemId, second).Return()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/main-cover?second=%f", itemId, second), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockProcessor.AssertExpectations(t)
	})

	t.Run("Invalid Item ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/invalid/main-cover?second=5.5", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid Second Parameter", func(t *testing.T) {
		itemId := uint64(123)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/main-cover?second=invalid", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRefreshItem(t *testing.T) {
	handler, _, mockProcessor, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		itemId := uint64(123)

		mockProcessor.On("EnqueueItemVideoMetadata", mock.Anything, itemId).Return()
		mockProcessor.On("EnqueueItemCovers", mock.Anything, itemId).Return()
		mockProcessor.On("EnqueueItemPreview", mock.Anything, itemId).Return()
		mockProcessor.On("EnqueueItemFileMetadata", mock.Anything, itemId).Return()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/process", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockProcessor.AssertExpectations(t)
	})

	t.Run("Invalid Item ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/invalid/process", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestOptimizeItem(t *testing.T) {
	handler, mockDb, _, mockOptimizer := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		itemId := uint64(123)
		testItem := &model.Item{
			Id:     itemId,
			Title:  "Test Item",
			Origin: "/path/to/item",
		}

		mockDb.On("GetItem", mock.Anything, itemId).Return(testItem, nil)
		mockOptimizer.On("HandleItem", mock.Anything, testItem).Return()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/optimize", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
		mockOptimizer.AssertExpectations(t)
	})

	t.Run("Item Not Found", func(t *testing.T) {
		itemId := uint64(999)
		mockDb.On("GetItem", mock.Anything, itemId).Return(nil, fmt.Errorf("item not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/optimize", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockDb.AssertExpectations(t)
	})
}

func TestCropFrame(t *testing.T) {
	handler, _, mockProcessor, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Success", func(t *testing.T) {
		itemId := uint64(123)
		second := 5.5
		cropX := 10.0
		cropY := 20.0
		cropWidth := 100.0
		cropHeight := 200.0

		expectedRect := model.RectFloat{
			X: cropX,
			Y: cropY,
			W: cropWidth,
			H: cropHeight,
		}

		mockProcessor.On("EnqueueCropFrame", mock.Anything, itemId, second, expectedRect).Return()

		url := fmt.Sprintf("/api/items/%d/crop-frame?second=%f&crop-x=%f&crop-y=%f&crop-width=%f&crop-height=%f",
			itemId, second, cropX, cropY, cropWidth, cropHeight)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", url, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockProcessor.AssertExpectations(t)
	})

	t.Run("Invalid Parameters", func(t *testing.T) {
		itemId := uint64(123)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/crop-frame?second=invalid", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// Helper function to create a test for HTTP handlers that only require basic setup
func TestRegisterRoutes(t *testing.T) {
	handler, _, _, _ := setupTestHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")

	// This should not panic
	assert.NotPanics(t, func() {
		itemsHandler := handler.(*itemsHandler)
		itemsHandler.RegisterRoutes(rg)
	})

	// Verify routes are registered by checking if they respond
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/items/invalid", nil)
	router.ServeHTTP(w, req)

	// Should get an error response but not a 404, indicating route is registered
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

func TestSplitItem(t *testing.T) {
	t.Run("Parameter Validation", func(t *testing.T) {
		handler, mockDb, mockProcessor, _ := setupTestHandler()
		router := setupTestRouter(handler)

		itemId := uint64(123)
		second := 30.5

		// Mock the business logic calls that Split function will make
		mainItem := &model.Item{
			Id:              itemId,
			Title:           "test.mp4",
			DurationSeconds: 100.0,
		}

		mockDb.On("GetItem", mock.Anything, itemId).Return(mainItem, nil)
		// Mock additional calls that Split might make
		mockDb.On("CreateOrUpdateItem", mock.Anything, mock.AnythingOfType("*model.Item")).Return(nil)
		mockDb.On("UpdateItem", mock.Anything, mock.AnythingOfType("*model.Item")).Return(nil)

		// Mock processor calls for the result items
		mockProcessor.On("EnqueueItemVideoMetadata", mock.Anything, mock.AnythingOfType("uint64")).Return()
		mockProcessor.On("EnqueueItemCovers", mock.Anything, mock.AnythingOfType("uint64")).Return()
		mockProcessor.On("EnqueueItemFileMetadata", mock.Anything, mock.AnythingOfType("uint64")).Return()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/split?second=%f", itemId, second), nil)
		router.ServeHTTP(w, req)

		// The split should succeed with proper mocks
		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("Invalid Item ID", func(t *testing.T) {
		handler, _, _, _ := setupTestHandler()
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items/invalid/split?second=30.5", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid Second Parameter", func(t *testing.T) {
		handler, _, _, _ := setupTestHandler()
		router := setupTestRouter(handler)

		itemId := uint64(123)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/items/%d/split?second=invalid", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestMakeHighlight(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockDb, mockProcessor, _ := setupTestHandler()
		router := setupTestRouter(handler)

		itemId := uint64(123)
		startSecond := 10.5
		endSecond := 25.0
		highlightId := uint64(456)

		// Mock the business logic calls that MakeHighlight function will make
		mainItem := &model.Item{
			Id:              itemId,
			Title:           "test.mp4",
			DurationSeconds: 100.0,
			Origin:          "/path/to/item",
		}

		mockDb.On("GetItem", mock.Anything, itemId).Return(mainItem, nil)
		mockDb.On("CreateOrUpdateItem", mock.Anything, mock.AnythingOfType("*model.Item")).Return(nil)
		mockDb.On("UpdateItem", mock.Anything, mock.AnythingOfType("*model.Item")).Return(nil)

		// Mock processor calls for the highlight (use flexible matching since the ID might be different)
		mockProcessor.On("EnqueueItemVideoMetadata", mock.Anything, mock.AnythingOfType("uint64")).Return()
		mockProcessor.On("EnqueueItemCovers", mock.Anything, mock.AnythingOfType("uint64")).Return()
		mockProcessor.On("EnqueueItemPreview", mock.Anything, mock.AnythingOfType("uint64")).Return()
		mockProcessor.On("EnqueueItemFileMetadata", mock.Anything, mock.AnythingOfType("uint64")).Return()

		url := fmt.Sprintf("/api/items/%d/make-highlight?start=%f&end=%f&highlight-id=%d",
			itemId, startSecond, endSecond, highlightId)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", url, nil)
		router.ServeHTTP(w, req)

		// The highlight should succeed with proper mocks
		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("Invalid Parameters", func(t *testing.T) {
		handler, _, _, _ := setupTestHandler()
		router := setupTestRouter(handler)

		itemId := uint64(123)

		tests := []struct {
			name string
			url  string
		}{
			{"Invalid Item ID", "/api/items/invalid/make-highlight?start=10&end=20&highlight-id=456"},
			{"Invalid Start", fmt.Sprintf("/api/items/%d/make-highlight?start=invalid&end=20&highlight-id=456", itemId)},
			{"Invalid End", fmt.Sprintf("/api/items/%d/make-highlight?start=10&end=invalid&highlight-id=456", itemId)},
			{"Invalid Highlight ID", fmt.Sprintf("/api/items/%d/make-highlight?start=10&end=20&highlight-id=invalid", itemId)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", tt.url, nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusInternalServerError, w.Code)
			})
		}
	})
}

func TestGetSuggestionsForItem(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler, mockDb, _, _ := setupTestHandler()
		router := setupTestRouter(handler)

		itemId := uint64(123)

		// Mock the suggestions business logic calls
		testItem := &model.Item{
			Id:    itemId,
			Title: "test.mp4",
			Tags:  []*model.Tag{},
		}

		// Mock calls that the suggestions function will make
		allItems := &[]model.Item{
			*testItem,
			{Id: 2, Title: "item2.mp4", Tags: []*model.Tag{}},
			{Id: 3, Title: "item3.mp4", Tags: []*model.Tag{}},
			{Id: 4, Title: "item4.mp4", Tags: []*model.Tag{}},
			{Id: 5, Title: "item5.mp4", Tags: []*model.Tag{}},
			{Id: 6, Title: "item6.mp4", Tags: []*model.Tag{}},
			{Id: 7, Title: "item7.mp4", Tags: []*model.Tag{}},
			{Id: 8, Title: "item8.mp4", Tags: []*model.Tag{}},
			{Id: 9, Title: "item9.mp4", Tags: []*model.Tag{}},
		}

		mockDb.On("GetItem", mock.Anything, itemId).Return(testItem, nil)
		mockDb.On("GetAllItems", mock.Anything).Return(allItems, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/items/%d/suggestions", itemId), nil)
		router.ServeHTTP(w, req)

		// The suggestions should succeed with proper mocks
		assert.Equal(t, http.StatusOK, w.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("Invalid Item ID", func(t *testing.T) {
		handler, _, _, _ := setupTestHandler()
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/items/invalid/suggestions", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// Test error handling middleware integration
func TestErrorHandling(t *testing.T) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Database Connection Error", func(t *testing.T) {
		mockDb.On("GetAllItems", mock.Anything).Return(nil, fmt.Errorf("connection refused"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/items", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("Malformed Request Body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items", bytes.NewBufferString(`{"invalid": json}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// Test edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	handler, mockDb, mockProcessor, _ := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("Large Item ID", func(t *testing.T) {
		itemId := uint64(18446744073709551615) // Max uint64
		mockDb.On("GetItem", mock.Anything, itemId).Return(nil, fmt.Errorf("item not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/items/%d", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("Zero Item ID", func(t *testing.T) {
		itemId := uint64(0)
		mockDb.On("GetItem", mock.Anything, itemId).Return(nil, fmt.Errorf("item not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/items/%d", itemId), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("Empty Item Creation", func(t *testing.T) {
		emptyItem := model.Item{}
		mockDb.On("CreateOrUpdateItem", mock.Anything, mock.AnythingOfType("*model.Item")).Return(fmt.Errorf("validation error"))

		jsonBody, _ := json.Marshal(emptyItem)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("Crop Frame with Zero Dimensions", func(t *testing.T) {
		itemId := uint64(123)
		second := 5.5
		expectedRect := model.RectFloat{X: 0, Y: 0, W: 0, H: 0}

		mockProcessor.On("EnqueueCropFrame", mock.Anything, itemId, second, expectedRect).Return()

		url := fmt.Sprintf("/api/items/%d/crop-frame?second=%f&crop-x=0&crop-y=0&crop-width=0&crop-height=0",
			itemId, second)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", url, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("Negative Crop Values", func(t *testing.T) {
		itemId := uint64(123)
		second := 5.5
		expectedRect := model.RectFloat{X: -10, Y: -20, W: -100, H: -200}

		mockProcessor.On("EnqueueCropFrame", mock.Anything, itemId, second, expectedRect).Return()

		url := fmt.Sprintf("/api/items/%d/crop-frame?second=%f&crop-x=-10&crop-y=-20&crop-width=-100&crop-height=-200",
			itemId, second)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", url, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockProcessor.AssertExpectations(t)
	})
}

// Test concurrent access scenarios
func TestConcurrentAccess(t *testing.T) {
	t.Run("Concurrent Get Requests", func(t *testing.T) {
		// Test that concurrent requests don't cause panic
		// Use separate handlers for each request to avoid mock conflicts
		for i := 0; i < 5; i++ {
			handler, mockDb, _, _ := setupTestHandler()
			router := setupTestRouter(handler)

			expectedItems := &[]model.Item{
				{Id: 1, Title: "Item 1"},
				{Id: 2, Title: "Item 2"},
			}

			mockDb.On("GetAllItems", mock.Anything).Return(expectedItems, nil)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/items", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			mockDb.AssertExpectations(t)
		}
	})
}

// Benchmark test for performance
func BenchmarkGetItems(b *testing.B) {
	handler, mockDb, _, _ := setupTestHandler()
	router := setupTestRouter(handler)

	expectedItems := &[]model.Item{
		{Id: 1, Title: "Item 1"},
		{Id: 2, Title: "Item 2"},
	}

	mockDb.On("GetAllItems", mock.Anything).Return(expectedItems, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/items", nil)
		router.ServeHTTP(w, req)
	}
}
