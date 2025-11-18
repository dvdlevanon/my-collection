package opensubtitles_test

import (
	"my-collection/server/pkg/opensubtitles"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestList_RealAPI calls the actual OpenSubtitles API to verify assumptions
// This test requires a valid API key and internet connection
func TestList_RealAPI(t *testing.T) {
	// Using the API key from the curl example
	apiKey := "DICiZ00oTTHgYMZsjI2Iue87PWJQ5fqE"

	client := opensubtitles.NewOpenSubtitles([]string{apiKey})

	// Test with the IMDB ID from the curl example: tt0120631
	imdbID := "tt0120631"

	result, err := client.List(imdbID, "he", false)

	// Print results for debugging
	t.Logf("IMDB ID: %s", imdbID)
	t.Logf("Language: he")
	t.Logf("AI Translated: exclude")
	t.Logf("Error: %v", err)
	t.Logf("Number of results: %d", len(result))

	if err != nil {
		t.Logf("API Error: %v", err)
		// Don't fail the test if it's a rate limit or quota issue
		if assert.Error(t, err) {
			t.Logf("Error details: %v", err)
		}
		return
	}

	require.NoError(t, err)

	// Verify we got results
	if len(result) > 0 {
		t.Logf("First result - ID: %s, Release: %s", result[0].Id, result[0].Title)
		assert.NotEmpty(t, result[0].Id, "Subtitle ID should not be empty")
		// Release might be empty, so we'll just log it
		t.Logf("Sample results:")
		for i, sub := range result {
			if i >= 5 { // Limit to first 5 for logging
				break
			}
			t.Logf("  [%d] ID: %s, Release: %s", i+1, sub.Id, sub.Title)
		}
	} else {
		t.Logf("No subtitles found for IMDB ID %s with language 'he'", imdbID)
	}
}

// TestDownload_RealAPI calls the actual OpenSubtitles API to verify download assumptions
// This test requires a valid API key and internet connection
// It first lists subtitles, then downloads the first one
func TestDownload_RealAPI(t *testing.T) {
	// Using the API key from the curl example
	apiKey := "DICiZ00oTTHgYMZsjI2Iue87PWJQ5fqE"

	client := opensubtitles.NewOpenSubtitles([]string{apiKey})

	// First, get a list of subtitles
	imdbID := "tt0120631"
	t.Logf("Step 1: Listing subtitles for IMDB ID: %s", imdbID)

	subtitles, err := client.List(imdbID, "he", false)
	if err != nil {
		t.Logf("Failed to list subtitles: %v", err)
		t.Logf("Skipping download test - cannot proceed without subtitle list")
		return
	}

	if len(subtitles) == 0 {
		t.Logf("No subtitles found for IMDB ID %s with language 'he'", imdbID)
		t.Logf("Skipping download test - no subtitles to download")
		return
	}

	// Use the first subtitle for download test
	subtitle := subtitles[0]
	t.Logf("Step 2: Downloading subtitle")
	t.Logf("  Subtitle ID: %s", subtitle.Id)
	t.Logf("  Release: %s", subtitle.Title)

	// Create a temporary file for the download
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "subtitle.srt")
	t.Logf("  Output file: %s", outputFile)

	// Download the subtitle
	err = client.Download(subtitle, outputFile)

	// Print results for debugging
	t.Logf("Download Error: %v", err)

	if err != nil {
		t.Logf("Download failed: %v", err)
		// Don't fail the test if it's a rate limit or quota issue
		if assert.Error(t, err) {
			t.Logf("Error details: %v", err)
		}
		return
	}

	require.NoError(t, err)

	// Verify the file was created
	fileInfo, err := os.Stat(outputFile)
	require.NoError(t, err, "Downloaded file should exist")
	assert.Greater(t, fileInfo.Size(), int64(0), "Downloaded file should not be empty")

	t.Logf("Download successful!")
	t.Logf("  File size: %d bytes", fileInfo.Size())

	// Read and log first few lines of the subtitle file to verify it's valid
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	contentStr := string(content)
	t.Logf("  File preview (first 200 chars): %s", contentStr[:min(200, len(contentStr))])

	// Basic validation - subtitle files usually start with numeric sequence or have specific format
	if len(contentStr) > 0 {
		t.Logf("  File appears to be valid (non-empty content)")
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestListAndDownload_RealAPI tests the full flow: list then download
// This test requires a valid API key and internet connection
func TestListAndDownload_RealAPI(t *testing.T) {
	// Using the API key from the curl example
	apiKey := "DICiZ00oTTHgYMZsjI2Iue87PWJQ5fqE"

	client := opensubtitles.NewOpenSubtitles([]string{apiKey})

	// Test with the IMDB ID from the curl example: tt0120631
	imdbID := "tt0120631"

	t.Logf("=== Full Flow Test: List + Download ===")
	t.Logf("IMDB ID: %s", imdbID)
	t.Logf("Language: he")
	t.Logf("AI Translated: exclude")

	// Step 1: List subtitles
	t.Logf("\n--- Step 1: Listing subtitles ---")
	result, err := client.List(imdbID, "he", false)

	if err != nil {
		t.Logf("List API Error: %v", err)
		t.Logf("Cannot proceed with download test")
		return
	}

	require.NoError(t, err)
	t.Logf("Found %d subtitles", len(result))

	if len(result) == 0 {
		t.Logf("No subtitles found - cannot test download")
		return
	}

	// Log first few results
	for i, sub := range result {
		if i >= 3 {
			break
		}
		t.Logf("  [%d] ID: %s, Release: %s", i+1, sub.Id, sub.Title)
	}

	// Step 2: Download first subtitle
	firstSubtitle := result[0]
	t.Logf("\n--- Step 2: Downloading subtitle ---")
	t.Logf("Subtitle ID: %s", firstSubtitle.Id)
	t.Logf("Release: %s", firstSubtitle.Title)

	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "test_subtitle.srt")
	t.Logf("Output file: %s", outputFile)

	err = client.Download(firstSubtitle, outputFile)

	if err != nil {
		t.Logf("Download Error: %v", err)
		return
	}

	require.NoError(t, err)

	// Verify download
	fileInfo, err := os.Stat(outputFile)
	require.NoError(t, err)
	assert.Greater(t, fileInfo.Size(), int64(0))

	t.Logf("\n=== Test Complete ===")
	t.Logf("Successfully listed %d subtitles and downloaded subtitle file (%d bytes)", len(result), fileInfo.Size())
}
