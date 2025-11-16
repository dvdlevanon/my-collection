package tags

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
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// MockTagsHandlerDb implements a mock for the tagsHandlerDb interface
type MockTagsHandlerDb struct {
	mock.Mock
}

func (m *MockTagsHandlerDb) GetAllTags(ctx context.Context) (*[]model.Tag, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

func (m *MockTagsHandlerDb) GetTag(ctx context.Context, conds ...interface{}) (*model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagsHandlerDb) GetTags(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

func (m *MockTagsHandlerDb) GetTagsWithoutChildren(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Tag), args.Error(1)
}

func (m *MockTagsHandlerDb) CreateOrUpdateTag(ctx context.Context, tag *model.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagsHandlerDb) UpdateTag(ctx context.Context, tag *model.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagsHandlerDb) RemoveTag(ctx context.Context, tagId uint64) error {
	args := m.Called(ctx, tagId)
	return args.Error(0)
}

func (m *MockTagsHandlerDb) RemoveTagImageFromTag(ctx context.Context, tagId uint64, imageId uint64) error {
	args := m.Called(ctx, tagId, imageId)
	return args.Error(0)
}

func (m *MockTagsHandlerDb) GetTagCustomCommand(ctx context.Context, conds ...interface{}) (*[]model.TagCustomCommand, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.TagCustomCommand), args.Error(1)
}

func (m *MockTagsHandlerDb) GetAllTagCustomCommands(ctx context.Context) (*[]model.TagCustomCommand, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.TagCustomCommand), args.Error(1)
}

func (m *MockTagsHandlerDb) UpdateTagImage(ctx context.Context, image *model.TagImage) error {
	args := m.Called(ctx, image)
	return args.Error(0)
}

func (m *MockTagsHandlerDb) GetAllTagImageTypes(ctx context.Context) (*[]model.TagImageType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.TagImageType), args.Error(1)
}

func (m *MockTagsHandlerDb) GetTagImageType(ctx context.Context, conds ...interface{}) (*model.TagImageType, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TagImageType), args.Error(1)
}

func (m *MockTagsHandlerDb) CreateOrUpdateTagImageType(ctx context.Context, tit *model.TagImageType) error {
	args := m.Called(ctx, tit)
	return args.Error(0)
}

func (m *MockTagsHandlerDb) GetTagAnnotation(ctx context.Context, conds ...interface{}) (*model.TagAnnotation, error) {
	args := m.Called(append([]interface{}{ctx}, conds...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TagAnnotation), args.Error(1)
}

func (m *MockTagsHandlerDb) GetTagAnnotations(ctx context.Context, tagId uint64) ([]model.TagAnnotation, error) {
	args := m.Called(ctx, tagId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TagAnnotation), args.Error(1)
}

func (m *MockTagsHandlerDb) CreateTagAnnotation(ctx context.Context, tagAnnotation *model.TagAnnotation) error {
	args := m.Called(ctx, tagAnnotation)
	return args.Error(0)
}

func (m *MockTagsHandlerDb) RemoveTagAnnotationFromTag(ctx context.Context, tagId uint64, annotationId uint64) error {
	args := m.Called(ctx, tagId, annotationId)
	return args.Error(0)
}

// MockStorageUploader implements a mock for the storage uploader interface
type MockStorageUploader struct {
	mock.Mock
}

func (m *MockStorageUploader) GetStorageUrl(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *MockStorageUploader) GetFileForWriting(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func (m *MockStorageUploader) GetTempFile() string {
	args := m.Called()
	return args.String(0)
}

// MockTagsHandlerProcessor implements a mock for the processor interface
type MockTagsHandlerProcessor struct {
	mock.Mock
}

func (m *MockTagsHandlerProcessor) ProcessThumbnail(ctx context.Context, image *model.TagImage) error {
	args := m.Called(ctx, image)
	return args.Error(0)
}

// Test setup functions
func setupTagsTestHandler() (*tagsHandler, *MockTagsHandlerDb, *MockStorageUploader, *MockTagsHandlerProcessor) {
	mockDb := &MockTagsHandlerDb{}
	mockStorage := &MockStorageUploader{}
	mockProcessor := &MockTagsHandlerProcessor{}

	handler := NewHandler(mockDb, mockStorage, mockProcessor)
	return handler, mockDb, mockStorage, mockProcessor
}

func setupTagsTestRouter(handler *tagsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")
	handler.RegisterRoutes(rg)
	return router
}

// Tests for CRUD operations on tags
func TestTagsCRUDOperations(t *testing.T) {
	t.Run("Create Tag", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		t.Run("Success", func(t *testing.T) {
			tag := model.Tag{Title: "rock"}
			mockDb.On("CreateOrUpdateTag", mock.Anything, &tag).Return(nil).Run(func(args mock.Arguments) {
				// Simulate database setting the ID
				tagArg := args.Get(1).(*model.Tag)
				tagArg.Id = 123
			})

			body, _ := json.Marshal(tag)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var resultTag model.Tag
			err := json.Unmarshal(w.Body.Bytes(), &resultTag)
			require.NoError(t, err)
			assert.Equal(t, uint64(123), resultTag.Id)

			mockDb.AssertExpectations(t)
		})

		t.Run("Invalid JSON", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags", bytes.NewBufferString("{invalid json}"))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("Database Error", func(t *testing.T) {
			handler, mockDb, _, _ := setupTagsTestHandler()
			router := setupTagsTestRouter(handler)

			tag := model.Tag{Title: "rock"}
			mockDb.On("CreateOrUpdateTag", mock.Anything, &tag).Return(assert.AnError)

			body, _ := json.Marshal(tag)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			mockDb.AssertExpectations(t)
		})
	})

	t.Run("Get Tag", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		t.Run("Success", func(t *testing.T) {
			expectedTag := &model.Tag{Id: 123, Title: "rock"}
			mockDb.On("GetTag", mock.Anything, mock.Anything).Return(expectedTag, nil)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/tags/123", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var resultTag model.Tag
			err := json.Unmarshal(w.Body.Bytes(), &resultTag)
			require.NoError(t, err)
			assert.Equal(t, expectedTag.Id, resultTag.Id)
			assert.Equal(t, expectedTag.Title, resultTag.Title)

			mockDb.AssertExpectations(t)
		})

		t.Run("Invalid ID", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/tags/invalid", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("Tag Not Found", func(t *testing.T) {
			handler, mockDb, _, _ := setupTagsTestHandler()
			router := setupTagsTestRouter(handler)

			mockDb.On("GetTag", mock.Anything, mock.Anything).Return(nil, assert.AnError)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/tags/999", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			mockDb.AssertExpectations(t)
		})
	})

	t.Run("Update Tag", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		t.Run("Success", func(t *testing.T) {
			tag := model.Tag{Id: 123, Title: "updated-rock"}
			mockDb.On("UpdateTag", mock.Anything, &tag).Return(nil)

			body, _ := json.Marshal(tag)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags/123", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockDb.AssertExpectations(t)
		})

		t.Run("Mismatched IDs", func(t *testing.T) {
			tag := model.Tag{Id: 999, Title: "rock"}
			body, _ := json.Marshal(tag)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags/123", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("Database Error", func(t *testing.T) {
			handler, mockDb, _, _ := setupTagsTestHandler()
			router := setupTagsTestRouter(handler)

			tag := model.Tag{Id: 123, Title: "rock"}
			mockDb.On("UpdateTag", mock.Anything, &tag).Return(assert.AnError)

			body, _ := json.Marshal(tag)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags/123", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			mockDb.AssertExpectations(t)
		})
	})

	t.Run("Delete Tag", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		t.Run("Success", func(t *testing.T) {
			mockDb.On("RemoveTag", mock.Anything, uint64(123)).Return(nil)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/tags/123", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockDb.AssertExpectations(t)
		})

		t.Run("Database Error", func(t *testing.T) {
			handler, mockDb, _, _ := setupTagsTestHandler()
			router := setupTagsTestRouter(handler)

			mockDb.On("RemoveTag", mock.Anything, uint64(123)).Return(assert.AnError)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/tags/123", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			mockDb.AssertExpectations(t)
		})
	})

	t.Run("Get All Tags", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		t.Run("Success", func(t *testing.T) {
			expectedTags := &[]model.Tag{
				{Id: 1, Title: "rock"},
				{Id: 2, Title: "jazz"},
			}
			mockDb.On("GetAllTags", mock.Anything).Return(expectedTags, nil)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/tags", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var resultTags []model.Tag
			err := json.Unmarshal(w.Body.Bytes(), &resultTags)
			require.NoError(t, err)
			assert.Len(t, resultTags, 2)

			mockDb.AssertExpectations(t)
		})
	})
}

// Tests for special tag operations
func TestTagsSpecialOperations(t *testing.T) {
	t.Run("Get Categories", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		// Since this calls tags.GetCategories, we need to mock the underlying database calls
		expectedTags := &[]model.Tag{
			{Id: 1, Title: "music", ParentID: nil},
			{Id: 2, Title: "movies", ParentID: nil},
		}
		// Mock the GetTags call that tags.GetCategories makes
		mockDb.On("GetTags", mock.Anything, mock.Anything).Return(expectedTags, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/categories", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("Get Special Tags", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		expectedTags := &[]model.Tag{
			{Id: 1, Title: "special-tag-1"},
			{Id: 2, Title: "special-tag-2"},
		}
		mockDb.On("GetTagsWithoutChildren", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedTags, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/special-tags", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("Get Tag Image Types", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		expectedTypes := &[]model.TagImageType{
			{Id: 1, Nickname: "banner"},
			{Id: 2, Nickname: "thumbnail"},
		}
		mockDb.On("GetAllTagImageTypes", mock.Anything).Return(expectedTypes, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tags/tag-image-types", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resultTypes []model.TagImageType
		err := json.Unmarshal(w.Body.Bytes(), &resultTypes)
		require.NoError(t, err)
		assert.Len(t, resultTypes, 2)

		mockDb.AssertExpectations(t)
	})
}

// Tests for tag image operations
func TestTagsImageOperations(t *testing.T) {
	t.Run("Auto Image", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		t.Run("Success", func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir := t.TempDir()

			tag := &model.Tag{Id: 123, Title: "rock", Children: []*model.Tag{}}
			fileUrl := model.FileUrl{Url: tempDir}

			mockDb.On("GetTag", mock.Anything, uint64(123)).Return(tag, nil)
			// Since AutoImageChildren reads directory contents and no children exist,
			// it should return successfully without doing anything

			body, _ := json.Marshal(fileUrl)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags/123/auto-image", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockDb.AssertExpectations(t)
		})

		t.Run("Invalid Tag ID", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags/invalid/auto-image", bytes.NewBufferString("{}"))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})

	t.Run("Update Tag Image", func(t *testing.T) {
		handler, mockDb, _, mockProcessor := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		t.Run("Success", func(t *testing.T) {
			image := model.TagImage{Id: 456, TagId: 123, Url: "/path/to/image.jpg"}
			mockDb.On("UpdateTagImage", mock.Anything, &image).Return(nil)
			mockProcessor.On("ProcessThumbnail", mock.Anything, &image).Return(nil)

			body, _ := json.Marshal(image)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags/123/images/456", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			mockDb.AssertExpectations(t)
			// Note: ProcessThumbnail runs in goroutine, so we can't easily assert it
		})

		t.Run("Mismatched Image ID", func(t *testing.T) {
			image := model.TagImage{Id: 999, TagId: 123}
			body, _ := json.Marshal(image)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags/123/images/456", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("Mismatched Tag ID", func(t *testing.T) {
			image := model.TagImage{Id: 456, TagId: 999}
			body, _ := json.Marshal(image)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/tags/123/images/456", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})

	t.Run("Remove Tag Image From Tag", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		// Mock the tags.RemoveTagImages call - since this is a business logic function,
		// we'll need to mock its dependencies
		mockDb.On("GetTag", mock.Anything, mock.Anything).Return(&model.Tag{Id: 123}, nil)
		mockDb.On("RemoveTagImageFromTag", mock.Anything, uint64(123), uint64(456)).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/tags/123/tit/456", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// Tests for random mix operations
func TestTagsRandomMixOperations(t *testing.T) {
	t.Run("Random Mix Include", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		tag := &model.Tag{Id: 123, Title: "rock"}
		mockDb.On("GetTag", mock.Anything, mock.Anything).Return(tag, nil)
		mockDb.On("UpdateTag", mock.Anything, mock.AnythingOfType("*model.Tag")).Return(nil).Run(func(args mock.Arguments) {
			tagArg := args.Get(1).(*model.Tag)
			assert.NotNil(t, tagArg.NoRandom)
			assert.False(t, *tagArg.NoRandom)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tags/123/random-mix/include", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("Random Mix Exclude", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		tag := &model.Tag{Id: 123, Title: "rock"}
		mockDb.On("GetTag", mock.Anything, mock.Anything).Return(tag, nil)
		mockDb.On("UpdateTag", mock.Anything, mock.AnythingOfType("*model.Tag")).Return(nil).Run(func(args mock.Arguments) {
			tagArg := args.Get(1).(*model.Tag)
			assert.NotNil(t, tagArg.NoRandom)
			assert.True(t, *tagArg.NoRandom)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tags/123/random-mix/exclude", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDb.AssertExpectations(t)
	})
}

// Tests for tag annotations
func TestTagsAnnotationOperations(t *testing.T) {
	t.Run("Add Annotation To Tag", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		annotation := model.TagAnnotation{Title: "test annotation"}

		// Mock the tag_annotations.AddAnnotationToTag business logic calls
		// First it checks if annotation exists, then gets the tag, then creates/updates
		mockDb.On("GetTagAnnotation", mock.Anything, mock.AnythingOfType("*model.TagAnnotation")).Return(nil, gorm.ErrRecordNotFound)
		mockDb.On("CreateTagAnnotation", mock.Anything, mock.AnythingOfType("*model.TagAnnotation")).Return(nil).Run(func(args mock.Arguments) {
			// Simulate setting an ID
			annotationArg := args.Get(1).(*model.TagAnnotation)
			annotationArg.Id = 456
		})
		mockDb.On("GetTag", mock.Anything, uint64(123)).Return(&model.Tag{Id: 123, Annotations: []*model.TagAnnotation{}}, nil)
		mockDb.On("CreateOrUpdateTag", mock.Anything, mock.AnythingOfType("*model.Tag")).Return(nil)

		body, _ := json.Marshal(annotation)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tags/123/annotations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resultAnnotation model.TagAnnotation
		err := json.Unmarshal(w.Body.Bytes(), &resultAnnotation)
		require.NoError(t, err)
		assert.Equal(t, uint64(456), resultAnnotation.Id)
	})

	t.Run("Remove Annotation From Tag", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		mockDb.On("RemoveTagAnnotationFromTag", mock.Anything, uint64(123), uint64(456)).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/tags/123/annotations/456", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("Get Tag Available Annotations", func(t *testing.T) {
		handler, mockDb, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		expectedAnnotations := []model.TagAnnotation{
			{Id: 1, Title: "annotation 1"},
			{Id: 2, Title: "annotation 2"},
		}

		// Mock the tag_annotations.GetTagAvailableAnnotations business logic calls
		mockDb.On("GetTag", mock.Anything, mock.Anything).Return(&model.Tag{Id: 123}, nil)
		mockDb.On("GetTagAnnotations", mock.Anything, uint64(123)).Return(expectedAnnotations, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tags/123/available-annotations", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resultAnnotations []model.TagAnnotation
		err := json.Unmarshal(w.Body.Bytes(), &resultAnnotations)
		require.NoError(t, err)
		assert.Len(t, resultAnnotations, 2)
	})
}

// Tests for tag custom commands
func TestTagsCustomCommands(t *testing.T) {
	handler, mockDb, _, _ := setupTagsTestHandler()
	router := setupTagsTestRouter(handler)

	t.Run("Get All Tag Custom Commands", func(t *testing.T) {
		expectedCommands := &[]model.TagCustomCommand{
			{Id: 1, Title: "ffmpeg command", Arg: "ffmpeg -i input output"},
			{Id: 2, Title: "convert command", Arg: "convert resize"},
		}
		mockDb.On("GetAllTagCustomCommands", mock.Anything).Return(expectedCommands, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tags/123/tag-custom-commands", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resultCommands []model.TagCustomCommand
		err := json.Unmarshal(w.Body.Bytes(), &resultCommands)
		require.NoError(t, err)
		assert.Len(t, resultCommands, 2)

		mockDb.AssertExpectations(t)
	})
}

// Test route registration
func TestTagsRegisterRoutes(t *testing.T) {
	handler, _, _, _ := setupTagsTestHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/api")

	// This should not panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(rg)
	})
}

// Test error scenarios
func TestTagsErrorScenarios(t *testing.T) {
	t.Run("Invalid Tag ID in Various Endpoints", func(t *testing.T) {
		handler, _, _, _ := setupTagsTestHandler()
		router := setupTagsTestRouter(handler)

		endpoints := []struct {
			method string
			path   string
		}{
			{"GET", "/api/tags/invalid"},
			{"POST", "/api/tags/invalid"},
			{"DELETE", "/api/tags/invalid"},
			{"POST", "/api/tags/invalid/auto-image"},
			{"POST", "/api/tags/invalid/random-mix/include"},
			{"POST", "/api/tags/invalid/random-mix/exclude"},
			{"POST", "/api/tags/invalid/annotations"},
			{"DELETE", "/api/tags/invalid/annotations/123"},
			{"GET", "/api/tags/invalid/available-annotations"},
		}

		for _, endpoint := range endpoints {
			w := httptest.NewRecorder()
			body := bytes.NewBufferString("{}")
			req, _ := http.NewRequest(endpoint.method, endpoint.path, body)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code,
				fmt.Sprintf("Expected error for %s %s", endpoint.method, endpoint.path))
		}
	})
}
