package integration

import (
	"my-collection/server/test/testutils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestBasicFileOperations tests basic file add/remove operations
func TestBasicFileOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	// Initial sync with empty filesystem
	framework.Sync()

	// Add some files
	framework.CreateFile("video1.mp4", "test video content")
	framework.CreateFile("subdir/video2.mp4", "test video content 2")
	framework.CreateFile("subdir/nested/video3.mp4", "test video content 3")

	// Sync and verify files are added
	framework.Sync()

	framework.AssertItemExists("", "video1.mp4")
	framework.AssertDirectoryExists("subdir")
	framework.IncludeDirectory("subdir")
	framework.Sync()
	framework.AssertDirectoryExists("subdir")
	framework.AssertItemExists("subdir", "video2.mp4")
	framework.IncludeDirectory("subdir/nested")
	framework.Sync()
	framework.AssertDirectoryExists("subdir/nested")
	framework.AssertItemExists("subdir/nested", "video3.mp4")

	// Remove a file
	framework.DeleteFile("subdir/video2.mp4")
	framework.Sync()

	framework.AssertItemNotExists("subdir", "video2.mp4")
	framework.AssertItemExists("", "video1.mp4") // Other files should remain
	framework.AssertItemExists("subdir/nested", "video3.mp4")
}

// TestBasicDirectoryOperations tests basic directory add/remove operations
func TestBasicDirectoryOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Initial sync
	framework.Sync()

	// Create directories
	framework.CreateDir("movies")
	framework.CreateDir("movies/action")
	framework.CreateDir("movies/comedy")

	framework.Sync()

	framework.AssertDirectoryExists("movies")
	framework.AssertDirectoryExists("movies/action")
	framework.AssertDirectoryExists("movies/comedy")

	// Remove a directory
	framework.DeleteDir("movies/comedy")
	framework.Sync()

	framework.AssertDirectoryNotExists("movies/comedy")
	framework.AssertDirectoryExists("movies")
	framework.AssertDirectoryExists("movies/action")
}

// TestFileMovementOperations tests file movement between directories
func TestFileMovementOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Setup initial structure
	framework.CreateFile("source/movie1.mp4", "content1")
	framework.CreateFile("source/movie2.mp4", "content2")
	framework.CreateDir("destination")
	framework.Sync()

	// Verify initial state
	framework.AssertItemExists("source", "movie1.mp4")
	framework.AssertItemExists("source", "movie2.mp4")

	// Move a file
	framework.MoveFile("source/movie1.mp4", "destination/movie1.mp4")
	framework.Sync()

	// Verify move
	framework.AssertItemNotExists("source", "movie1.mp4")
	framework.AssertItemExists("destination", "movie1.mp4")
	framework.AssertItemExists("source", "movie2.mp4") // Other file should remain
}

// TestDirectoryMovementOperations tests directory movement/renaming
func TestDirectoryMovementOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Setup initial structure with files in directories
	framework.CreateFile("oldname/file1.mp4", "content1")
	framework.CreateFile("oldname/subdir/file2.mp4", "content2")
	framework.Sync()

	// Verify initial state
	framework.AssertDirectoryExists("oldname")
	framework.AssertDirectoryExists("oldname/subdir")
	framework.AssertItemExists("oldname", "file1.mp4")
	framework.AssertItemExists("oldname/subdir", "file2.mp4")

	// Rename directory
	framework.MoveDir("oldname", "newname")
	framework.Sync()

	// Verify rename
	framework.AssertDirectoryNotExists("oldname")
	framework.AssertDirectoryExists("newname")
	framework.AssertDirectoryExists("newname/subdir")
	framework.AssertItemExists("newname", "file1.mp4")
	framework.AssertItemExists("newname/subdir", "file2.mp4")
}

// TestComplexHierarchyOperations tests complex nested directory operations
func TestComplexHierarchyOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Create a complex hierarchy
	paths := []string{
		"media/movies/action/2023/movie1.mp4",
		"media/movies/action/2023/movie2.mp4",
		"media/movies/comedy/2022/movie3.mp4",
		"media/tv/series1/season1/episode1.mp4",
		"media/tv/series1/season1/episode2.mp4",
		"media/tv/series1/season2/episode1.mp4",
	}

	for _, path := range paths {
		framework.CreateFile(path, "content")
	}
	framework.Sync()

	// Verify all items exist
	framework.AssertItemExists("media/movies/action/2023", "movie1.mp4")
	framework.AssertItemExists("media/movies/action/2023", "movie2.mp4")
	framework.AssertItemExists("media/movies/comedy/2022", "movie3.mp4")
	framework.AssertItemExists("media/tv/series1/season1", "episode1.mp4")
	framework.AssertItemExists("media/tv/series1/season1", "episode2.mp4")
	framework.AssertItemExists("media/tv/series1/season2", "episode1.mp4")

	// Reorganize: move TV series to separate location
	framework.MoveDir("media/tv/series1", "tv-shows/series1")
	framework.Sync()

	// Verify reorganization
	framework.AssertDirectoryNotExists("media/tv/series1")
	framework.AssertDirectoryExists("tv-shows/series1")
	framework.AssertItemExists("tv-shows/series1/season1", "episode1.mp4")
	framework.AssertItemExists("tv-shows/series1/season1", "episode2.mp4")
	framework.AssertItemExists("tv-shows/series1/season2", "episode1.mp4")

	// Movies should remain unchanged
	framework.AssertItemExists("media/movies/action/2023", "movie1.mp4")
	framework.AssertItemExists("media/movies/comedy/2022", "movie3.mp4")
}

// TestStaleFileHandling tests handling of stale files in database
func TestStaleFileHandling(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Create files and sync
	framework.CreateFile("temp/file1.mp4", "content1")
	framework.CreateFile("temp/file2.mp4", "content2")
	framework.Sync()

	framework.AssertItemExists("temp", "file1.mp4")
	framework.AssertItemExists("temp", "file2.mp4")

	// Delete entire directory from filesystem (simulating external changes)
	framework.DeleteDir("temp")
	framework.Sync()

	// Stale items should be removed from database
	framework.AssertItemNotExists("temp", "file1.mp4")
	framework.AssertItemNotExists("temp", "file2.mp4")
	framework.AssertDirectoryNotExists("temp")
}

// TestMixedOperations tests multiple operations in sequence
func TestMixedOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Step 1: Create initial structure
	framework.CreateFile("movies/action/movie1.mp4", "content1")
	framework.CreateFile("movies/comedy/movie2.mp4", "content2")
	framework.Sync()

	// Step 2: Add more files
	framework.CreateFile("movies/action/movie3.mp4", "content3")
	framework.CreateFile("tv/show1/episode1.mp4", "content4")
	framework.Sync()

	// Step 3: Move and rename
	framework.MoveFile("movies/comedy/movie2.mp4", "movies/action/movie2.mp4")
	framework.MoveDir("tv", "television")
	framework.Sync()

	// Step 4: Delete some items
	framework.DeleteFile("movies/action/movie1.mp4")
	framework.DeleteDir("movies/comedy")
	framework.Sync()

	// Verify final state
	framework.AssertItemNotExists("movies/action", "movie1.mp4")
	framework.AssertItemExists("movies/action", "movie2.mp4")
	framework.AssertItemExists("movies/action", "movie3.mp4")
	framework.AssertDirectoryNotExists("movies/comedy")
	framework.AssertDirectoryNotExists("tv")
	framework.AssertDirectoryExists("television")
	framework.AssertItemExists("television/show1", "episode1.mp4")
}

// TestRealisticLibraryOperations tests operations on a realistic media library
func TestRealisticLibraryOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Create realistic library structure
	framework.CreateTestLibrary()
	framework.Sync()

	// Verify library was created correctly
	framework.AssertTestLibraryExists()

	// Reorganize: move all 2023 content to archive
	framework.MoveFile("movies/action/2023/terminator.mp4", "archive/2023/action/terminator.mp4")
	framework.MoveFile("movies/comedy/2023/ghostbusters.mp4", "archive/2023/comedy/ghostbusters.mp4")
	framework.Sync()

	// Verify reorganization
	framework.AssertItemExists("archive/2023/action", "terminator.mp4")
	framework.AssertItemExists("archive/2023/comedy", "ghostbusters.mp4")
	framework.AssertItemNotExists("movies/action/2023", "terminator.mp4")
	framework.AssertItemNotExists("movies/comedy/2023", "ghostbusters.mp4")

	// Other content should remain
	framework.AssertItemExists("movies/action/2022", "john_wick.mp4")
	framework.AssertItemExists("tv/drama/breaking_bad/s01", "e01.mp4")
}

// TestSyncPerformance tests synchronization performance
func TestSyncPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Initial sync
	framework.Sync()

	// Create moderate number of files
	framework.CreateLargeTestSet(20, 25) // 500 files across 20 directories

	// Measure sync time
	start := time.Now()
	framework.Sync()
	syncDuration := time.Since(start)

	t.Logf("Sync of 500 files took: %v", syncDuration)

	// Verify all files were synced
	itemCount := framework.CountItemsInDirectory("large_test")
	assert.GreaterOrEqual(t, itemCount, 500, "Should have synced all 500 files")

	// Test incremental sync performance
	start = time.Now()
	framework.Sync() // No changes
	noChangeDuration := time.Since(start)

	t.Logf("No-change sync took: %v", noChangeDuration)

	// No-change sync should be very fast
	assert.Less(t, noChangeDuration, time.Second, "No-change sync should be fast")
}
