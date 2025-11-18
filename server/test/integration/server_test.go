package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/server/items"
	"my-collection/server/test/testutils"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockProcessor implements a simple mock for the processor interface
type MockProcessor struct{}

func (m *MockProcessor) EnqueueItemVideoMetadata(ctx context.Context, id uint64) {}
func (m *MockProcessor) EnqueueItemCovers(ctx context.Context, id uint64)        {}
func (m *MockProcessor) EnqueueCropFrame(ctx context.Context, id uint64, second float64, rect model.RectFloat) {
}
func (m *MockProcessor) EnqueueItemPreview(ctx context.Context, id uint64)               {}
func (m *MockProcessor) EnqueueItemFileMetadata(ctx context.Context, id uint64)          {}
func (m *MockProcessor) EnqueueMainCover(ctx context.Context, id uint64, second float64) {}

// MockOptimizer implements a simple mock for the optimizer interface
type MockOptimizer struct{}

func (m *MockOptimizer) HandleItem(ctx context.Context, item *model.Item) {}

// ServerIntegrationFramework extends the base framework with HTTP server capabilities
type ServerIntegrationFramework struct {
	*testutils.IntegrationTestFramework
	server     *http.Server
	httpClient *http.Client
	baseURL    string
	port       string
	t          *testing.T
}

// NewServerIntegrationFramework creates a new server integration test framework
func NewServerIntegrationFramework(t *testing.T) *ServerIntegrationFramework {
	baseFramework := testutils.NewIntegrationTestFramework(t)

	// Create a simple Gin router with just the items handler
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create items handler with mock dependencies
	mockProcessor := &MockProcessor{}
	mockOptimizer := &MockOptimizer{}
	itemsHandler := items.NewHandler(baseFramework.GetDatabase(), mockProcessor, mockOptimizer)

	// Register routes
	apiGroup := router.Group("/api")
	itemsHandler.RegisterRoutes(apiGroup)

	// Start HTTP server on a test port
	port := "18080"
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
		resp, err := client.Get(baseURL + "/api/items")
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
		if i == 49 {
			require.NoError(t, err, "Server failed to start within timeout")
		}
	}

	return &ServerIntegrationFramework{
		IntegrationTestFramework: baseFramework,
		server:                   httpServer,
		httpClient:               client,
		baseURL:                  baseURL,
		port:                     port,
		t:                        t,
	}
}

// Cleanup shuts down the server and cleans up resources
func (f *ServerIntegrationFramework) Cleanup() {
	if f.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		f.server.Shutdown(ctx)
	}
	f.IntegrationTestFramework.Cleanup()
}

// HTTP Helper Methods

// makeRequest makes an HTTP request and returns the response
func (f *ServerIntegrationFramework) makeRequest(method, path string, body interface{}) (*http.Response, error) {
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
func (f *ServerIntegrationFramework) makeRequestWithQuery(method, path string, query map[string]string, body interface{}) (*http.Response, error) {
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
func (f *ServerIntegrationFramework) parseJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// assertStatusCode asserts that the response has the expected status code
func (f *ServerIntegrationFramework) assertStatusCode(resp *http.Response, expected int) {
	assert.Equal(f.t, expected, resp.StatusCode,
		"Expected status code %d, got %d", expected, resp.StatusCode)
}

// Server API Test Methods

// GetAllItemsViaAPI gets all items via HTTP API
func (f *ServerIntegrationFramework) GetAllItemsViaAPI() ([]model.Item, error) {
	resp, err := f.makeRequest("GET", "/api/items", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status %d", resp.StatusCode)
	}

	var items []model.Item
	err = f.parseJSONResponse(resp, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetItemViaAPI gets an item by ID via HTTP API
func (f *ServerIntegrationFramework) GetItemViaAPI(itemID uint64) (*model.Item, error) {
	resp, err := f.makeRequest("GET", fmt.Sprintf("/api/items/%d", itemID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status %d", resp.StatusCode)
	}

	var item model.Item
	err = f.parseJSONResponse(resp, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// UpdateItemViaAPI updates an item via HTTP API
func (f *ServerIntegrationFramework) UpdateItemViaAPI(itemID uint64, item model.Item) error {
	resp, err := f.makeRequest("POST", fmt.Sprintf("/api/items/%d", itemID), item)
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

// DeleteItemViaAPI deletes an item via HTTP API
func (f *ServerIntegrationFramework) DeleteItemViaAPI(itemID uint64, deleteRealFile bool) error {
	query := map[string]string{
		"deleteRealFile": strconv.FormatBool(deleteRealFile),
	}
	resp, err := f.makeRequestWithQuery("DELETE", fmt.Sprintf("/api/items/%d", itemID), query, nil)
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

// Integration Tests

// TestServerBasicItemOperations tests basic CRUD operations on items
func TestServerBasicItemOperations(t *testing.T) {
	framework := NewServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Create and Get Item", func(t *testing.T) {
		// Create a test file
		framework.CreateFile("test_video.mp4", "test video content")
		framework.Sync()

		// Get all items via API
		items, err := framework.GetAllItemsViaAPI()
		require.NoError(t, err)
		require.Len(t, items, 1)

		item := items[0]
		assert.Equal(t, "test_video.mp4", item.Title)
		assert.NotEmpty(t, item.Origin)

		// Get specific item via API
		retrievedItem, err := framework.GetItemViaAPI(item.Id)
		require.NoError(t, err)
		assert.Equal(t, item.Id, retrievedItem.Id)
		assert.Equal(t, item.Title, retrievedItem.Title)
	})

	t.Run("Update Item", func(t *testing.T) {
		// Create a test file
		framework.CreateFile("update_test.mp4", "test content")
		framework.Sync()

		items, err := framework.GetAllItemsViaAPI()
		require.NoError(t, err)
		require.NotEmpty(t, items)

		item := items[0]
		originalTitle := item.Title

		// Update item metadata (not changing file-related fields)
		item.DurationSeconds = 120.5
		err = framework.UpdateItemViaAPI(item.Id, item)
		require.NoError(t, err)

		// Verify update
		updatedItem, err := framework.GetItemViaAPI(item.Id)
		require.NoError(t, err)
		assert.Equal(t, 120.5, updatedItem.DurationSeconds)
		assert.Equal(t, originalTitle, updatedItem.Title) // Title should remain the same
	})

	t.Run("Delete Item", func(t *testing.T) {
		// Create a test file
		framework.CreateFile("delete_test.mp4", "test content")
		framework.Sync()

		items, err := framework.GetAllItemsViaAPI()
		require.NoError(t, err)
		require.NotEmpty(t, items)

		item := items[0]
		itemID := item.Id

		// Delete item (not the real file)
		err = framework.DeleteItemViaAPI(itemID, false)
		require.NoError(t, err)

		// Verify item is deleted
		_, err = framework.GetItemViaAPI(itemID)
		assert.Error(t, err, "Item should not exist after deletion")

		// Verify item count decreased
		itemsAfterDelete, err := framework.GetAllItemsViaAPI()
		require.NoError(t, err)
		assert.Len(t, itemsAfterDelete, len(items)-1)
	})
}

// TestServerItemAdvancedOperations tests advanced operations like main cover, split, highlight
func TestServerItemAdvancedOperations(t *testing.T) {
	framework := NewServerIntegrationFramework(t)
	defer framework.Cleanup()

	// Create test files
	framework.CreateFile("movie.mp4", "movie content")
	framework.Sync()

	items, err := framework.GetAllItemsViaAPI()
	require.NoError(t, err)
	require.NotEmpty(t, items)

	itemID := items[0].Id

	t.Run("Set Main Cover", func(t *testing.T) {
		query := map[string]string{
			"second": "15.5",
		}
		resp, err := framework.makeRequestWithQuery("POST", fmt.Sprintf("/api/items/%d/main-cover", itemID), query, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)
	})

	t.Run("Refresh Item Processing", func(t *testing.T) {
		resp, err := framework.makeRequest("POST", fmt.Sprintf("/api/items/%d/process", itemID), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)
	})

	t.Run("Get Item Location", func(t *testing.T) {
		resp, err := framework.makeRequest("GET", fmt.Sprintf("/api/items/%d/location", itemID), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)

		var location model.FileUrl
		err = framework.parseJSONResponse(resp, &location)
		require.NoError(t, err)
		assert.NotEmpty(t, location.Url)
	})

	t.Run("Get Suggestions", func(t *testing.T) {
		// Now that we fixed the infinite loop bug, suggestions should work normally
		resp, err := framework.makeRequest("GET", fmt.Sprintf("/api/items/%d/suggestions", itemID), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)

		// Parse suggestions response (should be an array of items)
		var suggestions []model.Item
		err = framework.parseJSONResponse(resp, &suggestions)
		require.NoError(t, err)

		// With only 1 item in the database, suggestions should return empty array or single item
		assert.True(t, len(suggestions) <= 1, "Should return at most 1 suggestion when only 1 item exists")
	})

	t.Run("Optimize Item", func(t *testing.T) {
		resp, err := framework.makeRequest("POST", fmt.Sprintf("/api/items/%d/optimize", itemID), nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)
	})

	t.Run("Crop Frame", func(t *testing.T) {
		query := map[string]string{
			"second":      "10.0",
			"crop-x":      "10",
			"crop-y":      "20",
			"crop-width":  "100",
			"crop-height": "200",
		}
		resp, err := framework.makeRequestWithQuery("POST", fmt.Sprintf("/api/items/%d/crop-frame", itemID), query, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		framework.assertStatusCode(resp, http.StatusOK)
	})
}

// TestServerErrorHandling tests error scenarios and edge cases
func TestServerErrorHandling(t *testing.T) {
	framework := NewServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Get Non-existent Item", func(t *testing.T) {
		resp, err := framework.makeRequest("GET", "/api/items/99999", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Could be either 404 or 500 depending on how the handler processes the error
		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError,
			"Expected 404 or 500, got %d", resp.StatusCode)
	})

	t.Run("Invalid Item ID", func(t *testing.T) {
		resp, err := framework.makeRequest("GET", "/api/items/invalid", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Invalid JSON in Create", func(t *testing.T) {
		req, err := http.NewRequest("POST", framework.baseURL+"/api/items", bytes.NewBufferString("{invalid json}"))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := framework.httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Invalid Parameters in Advanced Operations", func(t *testing.T) {
		// Test invalid second parameter
		query := map[string]string{
			"second": "invalid",
		}
		resp, err := framework.makeRequestWithQuery("POST", "/api/items/1/main-cover", query, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Test invalid crop parameters
		invalidCropQuery := map[string]string{
			"second":      "10.0",
			"crop-x":      "invalid",
			"crop-y":      "20",
			"crop-width":  "100",
			"crop-height": "200",
		}
		resp2, err := framework.makeRequestWithQuery("POST", "/api/items/1/crop-frame", invalidCropQuery, nil)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp2.StatusCode)
	})
}

// TestServerWithRealFilesystem tests server operations with actual filesystem changes
func TestServerWithRealFilesystem(t *testing.T) {
	framework := NewServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Full Workflow: File to API", func(t *testing.T) {
		// Create files in root directory for simpler testing
		framework.CreateFile("terminator.mp4", "action movie")
		framework.CreateFile("ghostbusters.mp4", "comedy movie")
		framework.CreateFile("song.mp3", "rock music")

		// Sync filesystem to database
		framework.Sync()

		// Get items via API
		items, err := framework.GetAllItemsViaAPI()
		require.NoError(t, err)
		assert.Len(t, items, 3)

		// Verify each item exists and has correct data
		itemsByTitle := make(map[string]model.Item)
		for _, item := range items {
			itemsByTitle[item.Title] = item
		}

		assert.Contains(t, itemsByTitle, "terminator.mp4")
		assert.Contains(t, itemsByTitle, "ghostbusters.mp4")
		assert.Contains(t, itemsByTitle, "song.mp3")

		// Test updating one item via API
		terminatorItem := itemsByTitle["terminator.mp4"]
		terminatorItem.DurationSeconds = 108.0
		err = framework.UpdateItemViaAPI(terminatorItem.Id, terminatorItem)
		require.NoError(t, err)

		// Verify update
		updatedItem, err := framework.GetItemViaAPI(terminatorItem.Id)
		require.NoError(t, err)
		assert.Equal(t, 108.0, updatedItem.DurationSeconds)

		// Test deleting file from filesystem and syncing
		framework.DeleteFile("song.mp3")
		framework.Sync()

		// Verify item is removed from API
		finalItems, err := framework.GetAllItemsViaAPI()
		require.NoError(t, err)
		assert.Len(t, finalItems, 2)

		finalItemsByTitle := make(map[string]bool)
		for _, item := range finalItems {
			finalItemsByTitle[item.Title] = true
		}
		assert.False(t, finalItemsByTitle["song.mp3"])
		assert.True(t, finalItemsByTitle["terminator.mp4"])
		assert.True(t, finalItemsByTitle["ghostbusters.mp4"])
	})
}

// TestServerPerformanceAndStress tests performance aspects
func TestServerPerformanceAndStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	framework := NewServerIntegrationFramework(t)
	defer framework.Cleanup()

	t.Run("Multiple Concurrent Requests", func(t *testing.T) {
		// Create some test files
		for i := 0; i < 10; i++ {
			framework.CreateFile(fmt.Sprintf("test%d.mp4", i), fmt.Sprintf("content %d", i))
		}
		framework.Sync()

		// Make concurrent API requests
		const numConcurrent = 20
		results := make(chan error, numConcurrent)

		for i := 0; i < numConcurrent; i++ {
			go func() {
				_, err := framework.GetAllItemsViaAPI()
				results <- err
			}()
		}

		// Collect results
		for i := 0; i < numConcurrent; i++ {
			err := <-results
			assert.NoError(t, err, "Concurrent request %d should succeed", i)
		}
	})

	t.Run("Large Number of Items", func(t *testing.T) {
		// Create many files
		const numFiles = 100
		for i := 0; i < numFiles; i++ {
			framework.CreateFile(fmt.Sprintf("large/file%03d.mp4", i), fmt.Sprintf("content %d", i))
		}
		framework.Sync()

		// Measure API response time
		start := time.Now()
		items, err := framework.GetAllItemsViaAPI()
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Len(t, items, numFiles)
		assert.Less(t, duration, 5*time.Second, "API should respond within reasonable time even with many items")

		t.Logf("API returned %d items in %v", len(items), duration)
	})
}
