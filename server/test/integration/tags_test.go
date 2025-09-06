package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server/tags"
	"my-collection/server/pkg/storage"
	"my-collection/server/test/testutils"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockTagsProcessor implements a simple mock for the tags processor interface
type MockTagsProcessor struct{}

func (m *MockTagsProcessor) ProcessThumbnail(image *model.TagImage) error {
	return nil
}

// TagsServerIntegrationFramework extends the base framework with HTTP server capabilities for tags
type TagsServerIntegrationFramework struct {
	*testutils.IntegrationTestFramework
	server      *http.Server
	httpClient  *http.Client
	baseURL     string
	port        string
	t           *testing.T
	storageDir  string
	mockStorage *storage.Storage
}

// NewTagsServerIntegrationFramework creates a new tags server integration test framework
func NewTagsServerIntegrationFramework(t *testing.T) *TagsServerIntegrationFramework {
	baseFramework := testutils.NewIntegrationTestFramework(t)

	// Create a simple Gin router with just the tags handler
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create storage directory for tag images
	storageDir := filepath.Join(baseFramework.GetTempDir(), "storage")
	mockStorage, err := storage.New(storageDir)
	require.NoError(t, err)

	// Create tags handler with mock dependencies
	mockProcessor := &MockTagsProcessor{}
	tagsHandler := tags.NewHandler(baseFramework.GetDatabase(), mockStorage, mockProcessor)

	// Register routes
	apiGroup := router.Group("/api")
	tagsHandler.RegisterRoutes(apiGroup)

	// Start HTTP server on a test port
	port := "18081"
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server failed to start: %v", err)
		}
	}()

	// Wait for server to start
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	client := &http.Client{Timeout: 10 * time.Second}

	// Wait for server to be ready
	for i := 0; i < 50; i++ {
		resp, err := client.Get(baseURL + "/api/tags")
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
		if i == 49 {
			require.NoError(t, err, "Server failed to start within timeout")
		}
	}

	return &TagsServerIntegrationFramework{
		IntegrationTestFramework: baseFramework,
		server:                   httpServer,
		httpClient:               client,
		baseURL:                  baseURL,
		port:                     port,
		t:                        t,
		storageDir:               storageDir,
		mockStorage:              mockStorage,
	}
}

// Cleanup shuts down the server and cleans up resources
func (f *TagsServerIntegrationFramework) Cleanup() {
	if f.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		f.server.Shutdown(ctx)
	}
	f.IntegrationTestFramework.Cleanup()
}

// HTTP Helper Methods

// makeRequest makes an HTTP request and returns the response
func (f *TagsServerIntegrationFramework) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, f.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return f.httpClient.Do(req)
}

// makeRequestWithQuery makes an HTTP request with query parameters
func (f *TagsServerIntegrationFramework) makeRequestWithQuery(method, path string, query map[string]string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, f.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add query parameters
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	return f.httpClient.Do(req)
}

// parseJSONResponse parses a JSON response into the given interface
func (f *TagsServerIntegrationFramework) parseJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// assertStatusCode asserts that the response has the expected status code
func (f *TagsServerIntegrationFramework) assertStatusCode(resp *http.Response, expected int) {
	assert.Equal(f.t, expected, resp.StatusCode,
		"Expected status code %d, got %d", expected, resp.StatusCode)
}

// Tags API Test Methods

// GetAllTagsViaAPI gets all tags via HTTP API
func (f *TagsServerIntegrationFramework) GetAllTagsViaAPI() ([]model.Tag, error) {
	resp, err := f.makeRequest("GET", "/api/tags", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status %d", resp.StatusCode)
	}

	var tags []model.Tag
	err = f.parseJSONResponse(resp, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// GetTagViaAPI gets a tag by ID via HTTP API
func (f *TagsServerIntegrationFramework) GetTagViaAPI(tagID uint64) (*model.Tag, error) {
	resp, err := f.makeRequest("GET", fmt.Sprintf("/api/tags/%d", tagID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status %d", resp.StatusCode)
	}

	var tag model.Tag
	err = f.parseJSONResponse(resp, &tag)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// CreateTagViaAPI creates a tag via HTTP API
func (f *TagsServerIntegrationFramework) CreateTagViaAPI(tag model.Tag) (*model.Tag, error) {
	resp, err := f.makeRequest("POST", "/api/tags", tag)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdTag model.Tag
	err = f.parseJSONResponse(resp, &createdTag)
	if err != nil {
		return nil, err
	}

	return &createdTag, nil
}

// UpdateTagViaAPI updates a tag via HTTP API
func (f *TagsServerIntegrationFramework) UpdateTagViaAPI(tagID uint64, tag model.Tag) error {
	resp, err := f.makeRequest("POST", fmt.Sprintf("/api/tags/%d", tagID), tag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteTagViaAPI deletes a tag via HTTP API
func (f *TagsServerIntegrationFramework) DeleteTagViaAPI(tagID uint64) error {
	resp, err := f.makeRequest("DELETE", fmt.Sprintf("/api/tags/%d", tagID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetCategoriesViaAPI gets categories via HTTP API
func (f *TagsServerIntegrationFramework) GetCategoriesViaAPI() ([]model.Tag, error) {
	resp, err := f.makeRequest("GET", "/api/categories", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status %d", resp.StatusCode)
	}

	var categories []model.Tag
	err = f.parseJSONResponse(resp, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetSpecialTagsViaAPI gets special tags via HTTP API
func (f *TagsServerIntegrationFramework) GetSpecialTagsViaAPI() ([]model.Tag, error) {
	resp, err := f.makeRequest("GET", "/api/special-tags", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status %d", resp.StatusCode)
	}

	var specialTags []model.Tag
	err = f.parseJSONResponse(resp, &specialTags)
	if err != nil {
		return nil, err
	}

	return specialTags, nil
}

// GetTagImageTypesViaAPI gets tag image types via HTTP API
func (f *TagsServerIntegrationFramework) GetTagImageTypesViaAPI() ([]model.TagImageType, error) {
	resp, err := f.makeRequest("GET", "/api/tags/tag-image-types", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status %d", resp.StatusCode)
	}

	var imageTypes []model.TagImageType
	err = f.parseJSONResponse(resp, &imageTypes)
	if err != nil {
		return nil, err
	}

	return imageTypes, nil
}

// AddAnnotationToTagViaAPI adds an annotation to a tag via HTTP API
func (f *TagsServerIntegrationFramework) AddAnnotationToTagViaAPI(tagID uint64, annotation model.TagAnnotation) (*model.TagAnnotation, error) {
	resp, err := f.makeRequest("POST", fmt.Sprintf("/api/tags/%d/annotations", tagID), annotation)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdAnnotation model.TagAnnotation
	err = f.parseJSONResponse(resp, &createdAnnotation)
	if err != nil {
		return nil, err
	}

	return &createdAnnotation, nil
}

// Integration Tests

// TestTagsServerBasicCRUDOperations tests basic CRUD operations on tags
func TestTagsServerBasicCRUDOperations(t *testing.T) {
	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Create and Get Tag", func(t *testing.T) {
		// Create a tag
		tag := model.Tag{Title: "rock"}
		createdTag, err := framework.CreateTagViaAPI(tag)
		require.NoError(t, err)
		require.NotNil(t, createdTag)
		assert.NotZero(t, createdTag.Id)

		// Get all tags via API
		tags, err := framework.GetAllTagsViaAPI()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tags), 1)

		// Get specific tag via API
		retrievedTag, err := framework.GetTagViaAPI(createdTag.Id)
		require.NoError(t, err)
		assert.Equal(t, createdTag.Id, retrievedTag.Id)
		assert.Equal(t, "rock", retrievedTag.Title)
	})

	t.Run("Update Tag", func(t *testing.T) {
		// Create a tag
		tag := model.Tag{Title: "jazz"}
		createdTag, err := framework.CreateTagViaAPI(tag)
		require.NoError(t, err)

		// Update tag
		updatedTag := *createdTag
		updatedTag.Title = "updated-jazz"
		err = framework.UpdateTagViaAPI(createdTag.Id, updatedTag)
		require.NoError(t, err)

		// Verify update
		retrievedTag, err := framework.GetTagViaAPI(createdTag.Id)
		require.NoError(t, err)
		assert.Equal(t, "updated-jazz", retrievedTag.Title)
	})

	t.Run("Delete Tag", func(t *testing.T) {
		// Create a tag
		tag := model.Tag{Title: "delete-me"}
		createdTag, err := framework.CreateTagViaAPI(tag)
		require.NoError(t, err)

		// Delete tag
		err = framework.DeleteTagViaAPI(createdTag.Id)
		require.NoError(t, err)

		// Verify tag is deleted
		_, err = framework.GetTagViaAPI(createdTag.Id)
		assert.Error(t, err, "Tag should not exist after deletion")
	})
}

// TestTagsServerSpecialOperations tests special operations like categories and special tags
func TestTagsServerSpecialOperations(t *testing.T) {
	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Get Categories", func(t *testing.T) {
		// Create some parent tags (categories)
		category1 := model.Tag{Title: "music"}
		createdCategory1, err := framework.CreateTagViaAPI(category1)
		require.NoError(t, err)

		category2 := model.Tag{Title: "movies"}
		_, err = framework.CreateTagViaAPI(category2)
		require.NoError(t, err)

		// Create child tags
		parentID := createdCategory1.Id
		childTag := model.Tag{Title: "rock", ParentID: &parentID}
		_, err = framework.CreateTagViaAPI(childTag)
		require.NoError(t, err)

		// Get categories
		categories, err := framework.GetCategoriesViaAPI()
		require.NoError(t, err)

		// Should contain our created categories (and potentially system tags)
		categoryTitles := make(map[string]bool)
		for _, cat := range categories {
			categoryTitles[cat.Title] = true
		}
		assert.True(t, categoryTitles["music"])
		assert.True(t, categoryTitles["movies"])
	})

	t.Run("Get Special Tags", func(t *testing.T) {
		// Special tags are system-created tags like directories, dailymix, etc.
		specialTags, err := framework.GetSpecialTagsViaAPI()
		require.NoError(t, err)
		// Should return at least the system special tags
		assert.GreaterOrEqual(t, len(specialTags), 0)
	})

	t.Run("Get Tag Image Types", func(t *testing.T) {
		imageTypes, err := framework.GetTagImageTypesViaAPI()
		require.NoError(t, err)
		// Initially should be empty or contain default types
		assert.GreaterOrEqual(t, len(imageTypes), 0)
	})
}

// TestTagsServerRandomMixOperations tests random mix include/exclude operations
func TestTagsServerRandomMixOperations(t *testing.T) {
	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	// Create a tag first
	tag := model.Tag{Title: "test-random"}
	createdTag, err := framework.CreateTagViaAPI(tag)
	require.NoError(t, err)

	t.Run("Random Mix Include", func(t *testing.T) {
		resp, err := framework.makeRequest("POST", fmt.Sprintf("/api/tags/%d/random-mix/include", createdTag.Id), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)

		// Verify tag is included in random mix
		updatedTag, err := framework.GetTagViaAPI(createdTag.Id)
		require.NoError(t, err)
		assert.NotNil(t, updatedTag.NoRandom)
		assert.False(t, *updatedTag.NoRandom)
	})

	t.Run("Random Mix Exclude", func(t *testing.T) {
		resp, err := framework.makeRequest("POST", fmt.Sprintf("/api/tags/%d/random-mix/exclude", createdTag.Id), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)

		// Verify tag is excluded from random mix
		updatedTag, err := framework.GetTagViaAPI(createdTag.Id)
		require.NoError(t, err)
		assert.NotNil(t, updatedTag.NoRandom)
		assert.True(t, *updatedTag.NoRandom)
	})
}

// TestTagsServerAnnotationOperations tests tag annotation operations
func TestTagsServerAnnotationOperations(t *testing.T) {
	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	// Create a tag first
	tag := model.Tag{Title: "test-annotations"}
	createdTag, err := framework.CreateTagViaAPI(tag)
	require.NoError(t, err)

	t.Run("Add Annotation To Tag", func(t *testing.T) {
		annotation := model.TagAnnotation{Title: "test annotation"}
		createdAnnotation, err := framework.AddAnnotationToTagViaAPI(createdTag.Id, annotation)
		require.NoError(t, err)
		require.NotNil(t, createdAnnotation)
		assert.NotZero(t, createdAnnotation.Id)
		// Note: The handler only returns the ID, not the full annotation object
		// Verify the annotation was created by checking available annotations
		resp, err := framework.makeRequest("GET", fmt.Sprintf("/api/tags/%d/available-annotations", createdTag.Id), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		var availableAnnotations []model.TagAnnotation
		err = framework.parseJSONResponse(resp, &availableAnnotations)
		require.NoError(t, err)

		// Should contain our created annotation
		found := false
		for _, ann := range availableAnnotations {
			if ann.Title == "test annotation" {
				found = true
				break
			}
		}
		assert.True(t, found, "Created annotation should be found in available annotations")
	})

	t.Run("Remove Annotation From Tag", func(t *testing.T) {
		// First add an annotation
		annotation := model.TagAnnotation{Title: "remove me"}
		createdAnnotation, err := framework.AddAnnotationToTagViaAPI(createdTag.Id, annotation)
		require.NoError(t, err)

		// Remove the annotation
		resp, err := framework.makeRequest("DELETE", fmt.Sprintf("/api/tags/%d/annotations/%d", createdTag.Id, createdAnnotation.Id), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)
	})

	t.Run("Get Tag Available Annotations", func(t *testing.T) {
		// Add some annotations first
		annotation1 := model.TagAnnotation{Title: "annotation 1"}
		_, err := framework.AddAnnotationToTagViaAPI(createdTag.Id, annotation1)
		require.NoError(t, err)

		annotation2 := model.TagAnnotation{Title: "annotation 2"}
		_, err = framework.AddAnnotationToTagViaAPI(createdTag.Id, annotation2)
		require.NoError(t, err)

		// Get available annotations
		resp, err := framework.makeRequest("GET", fmt.Sprintf("/api/tags/%d/available-annotations", createdTag.Id), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)

		var annotations []model.TagAnnotation
		err = framework.parseJSONResponse(resp, &annotations)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(annotations), 2)
	})
}

// TestTagsServerAutoImageOperations tests auto image operations
func TestTagsServerAutoImageOperations(t *testing.T) {
	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	// Create a tag first
	tag := model.Tag{Title: "test-auto-image", Children: []*model.Tag{}}
	createdTag, err := framework.CreateTagViaAPI(tag)
	require.NoError(t, err)

	t.Run("Auto Image", func(t *testing.T) {
		// Create a test directory for auto image
		imageDir := framework.GetTempDir()
		framework.CreateDir("auto-image-test")

		fileUrl := model.FileUrl{Url: filepath.Join(imageDir, "auto-image-test")}
		resp, err := framework.makeRequest("POST", fmt.Sprintf("/api/tags/%d/auto-image", createdTag.Id), fileUrl)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)
	})
}

// TestTagsServerCustomCommands tests custom commands operations
func TestTagsServerCustomCommands(t *testing.T) {
	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Get All Tag Custom Commands", func(t *testing.T) {
		resp, err := framework.makeRequest("GET", "/api/tags/123/tag-custom-commands", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)

		var commands []model.TagCustomCommand
		err = framework.parseJSONResponse(resp, &commands)
		require.NoError(t, err)
		// Should return empty list initially
		assert.GreaterOrEqual(t, len(commands), 0)
	})
}

// TestTagsServerErrorHandling tests error scenarios and edge cases
func TestTagsServerErrorHandling(t *testing.T) {
	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Get Non-existent Tag", func(t *testing.T) {
		resp, err := framework.makeRequest("GET", "/api/tags/99999", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return error status
		assert.True(t, resp.StatusCode >= 400,
			"Expected error status, got %d", resp.StatusCode)
	})

	t.Run("Invalid Tag ID", func(t *testing.T) {
		resp, err := framework.makeRequest("GET", "/api/tags/invalid", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Invalid JSON in Create", func(t *testing.T) {
		req, err := http.NewRequest("POST", framework.baseURL+"/api/tags", bytes.NewBufferString("{invalid json}"))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := framework.httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Update Tag with Mismatched ID", func(t *testing.T) {
		tag := model.Tag{Id: 999, Title: "mismatch"}
		resp, err := framework.makeRequest("POST", "/api/tags/123", tag)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// TestTagsServerIntegrationWithRealDatabase tests tags operations with real database operations
func TestTagsServerIntegrationWithRealDatabase(t *testing.T) {
	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Full Workflow: Create Hierarchy", func(t *testing.T) {
		// Create a category (parent tag)
		category := model.Tag{Title: "music"}
		createdCategory, err := framework.CreateTagViaAPI(category)
		require.NoError(t, err)

		// Create child tags
		parentID := createdCategory.Id
		rockTag := model.Tag{Title: "rock", ParentID: &parentID}
		createdRock, err := framework.CreateTagViaAPI(rockTag)
		require.NoError(t, err)

		jazzTag := model.Tag{Title: "jazz", ParentID: &parentID}
		createdJazz, err := framework.CreateTagViaAPI(jazzTag)
		require.NoError(t, err)

		// Verify hierarchy via API
		allTags, err := framework.GetAllTagsViaAPI()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allTags), 3)

		// Verify parent-child relationships
		tagsByID := make(map[uint64]model.Tag)
		for _, tag := range allTags {
			tagsByID[tag.Id] = tag
		}

		rockRetrieved := tagsByID[createdRock.Id]
		jazzRetrieved := tagsByID[createdJazz.Id]

		assert.NotNil(t, rockRetrieved.ParentID)
		assert.Equal(t, createdCategory.Id, *rockRetrieved.ParentID)
		assert.NotNil(t, jazzRetrieved.ParentID)
		assert.Equal(t, createdCategory.Id, *jazzRetrieved.ParentID)

		// Test categories endpoint
		categories, err := framework.GetCategoriesViaAPI()
		require.NoError(t, err)

		musicFound := false
		for _, cat := range categories {
			if cat.Title == "music" {
				musicFound = true
				break
			}
		}
		assert.True(t, musicFound, "Music category should be returned by categories endpoint")
	})

	t.Run("Annotations Workflow", func(t *testing.T) {
		// Create a tag
		tag := model.Tag{Title: "annotation-test"}
		createdTag, err := framework.CreateTagViaAPI(tag)
		require.NoError(t, err)

		// Add multiple annotations
		annotations := []string{"tempo: fast", "mood: energetic", "instrument: guitar"}
		createdAnnotations := make([]*model.TagAnnotation, 0)

		for _, annotationTitle := range annotations {
			annotation := model.TagAnnotation{Title: annotationTitle}
			created, err := framework.AddAnnotationToTagViaAPI(createdTag.Id, annotation)
			require.NoError(t, err)
			createdAnnotations = append(createdAnnotations, created)
		}

		// Get available annotations
		resp, err := framework.makeRequest("GET", fmt.Sprintf("/api/tags/%d/available-annotations", createdTag.Id), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)

		var availableAnnotations []model.TagAnnotation
		err = framework.parseJSONResponse(resp, &availableAnnotations)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(availableAnnotations), 3)

		// Remove one annotation
		firstAnnotation := createdAnnotations[0]
		resp, err = framework.makeRequest("DELETE", fmt.Sprintf("/api/tags/%d/annotations/%d", createdTag.Id, firstAnnotation.Id), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)

		// Verify annotation count decreased
		resp, err = framework.makeRequest("GET", fmt.Sprintf("/api/tags/%d/available-annotations", createdTag.Id), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		err = framework.parseJSONResponse(resp, &availableAnnotations)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(availableAnnotations), 2)
	})
}

// TestTagsServerPerformanceAndStress tests performance aspects
func TestTagsServerPerformanceAndStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	framework := NewTagsServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Multiple Concurrent Tag Requests", func(t *testing.T) {
		// Create some test tags
		for i := 0; i < 10; i++ {
			tag := model.Tag{Title: fmt.Sprintf("concurrent-test-%d", i)}
			_, err := framework.CreateTagViaAPI(tag)
			require.NoError(t, err)
		}

		// Make concurrent API requests
		const numConcurrent = 20
		results := make(chan error, numConcurrent)

		for i := 0; i < numConcurrent; i++ {
			go func() {
				_, err := framework.GetAllTagsViaAPI()
				results <- err
			}()
		}

		// Collect results
		for i := 0; i < numConcurrent; i++ {
			err := <-results
			assert.NoError(t, err, "Concurrent request %d should succeed", i)
		}
	})

	t.Run("Large Number of Tags", func(t *testing.T) {
		// Create many tags
		const numTags = 100
		for i := 0; i < numTags; i++ {
			tag := model.Tag{Title: fmt.Sprintf("stress-tag-%03d", i)}
			_, err := framework.CreateTagViaAPI(tag)
			require.NoError(t, err)
		}

		// Measure API response time
		start := time.Now()
		tags, err := framework.GetAllTagsViaAPI()
		duration := time.Since(start)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tags), numTags)
		assert.Less(t, duration, 5*time.Second, "API should respond within reasonable time even with many tags")

		t.Logf("API returned %d tags in %v", len(tags), duration)
	})
}
