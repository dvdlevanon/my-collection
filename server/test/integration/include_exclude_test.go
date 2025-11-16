package integration

import (
	"context"
	"my-collection/server/pkg/model"
	"my-collection/server/test/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDirectoryIncludeExcludeBasic tests the basic include/exclude functionality
func TestDirectoryIncludeExcludeBasic(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	// Create some files in directories
	framework.CreateFile("movies/action/terminator.mp4", "action movie")
	framework.CreateFile("movies/comedy/ghostbusters.mp4", "comedy movie")
	framework.Sync()

	// By default, new directories should be excluded
	framework.AssertDirectoryExcluded("movies")
	framework.AssertDirectoryExcluded("movies/action")
	framework.AssertDirectoryExcluded("movies/comedy")

	// Files in excluded directories should not be processed
	framework.AssertItemNotExists("movies/action", "terminator.mp4")
	framework.AssertItemNotExists("movies/comedy", "ghostbusters.mp4")

	// User includes the movies directory
	framework.IncludeDirectory("movies")
	framework.Sync()

	// Now movies directory should be included
	framework.AssertDirectoryIncluded("movies")
	// But subdirectories should still be excluded by default
	framework.AssertDirectoryExcluded("movies/action")
	framework.AssertDirectoryExcluded("movies/comedy")

	// Files still not processed because subdirectories are excluded
	framework.AssertItemNotExists("movies/action", "terminator.mp4")
	framework.AssertItemNotExists("movies/comedy", "ghostbusters.mp4")

	// User includes the action subdirectory
	framework.IncludeDirectory("movies/action")
	framework.Sync()

	// Now action directory is included and its files are processed
	framework.AssertDirectoryIncluded("movies/action")
	framework.AssertItemExists("movies/action", "terminator.mp4")

	// But comedy is still excluded
	framework.AssertDirectoryExcluded("movies/comedy")
	framework.AssertItemNotExists("movies/comedy", "ghostbusters.mp4")

	// User excludes the action directory again
	framework.ExcludeDirectory("movies/action")
	framework.Sync()

	// Directory is now excluded and files are removed from index
	framework.AssertDirectoryExcluded("movies/action")
	framework.AssertItemNotExists("movies/action", "terminator.mp4")
}

// TestDirectoryIncludeExcludeHierarchy tests hierarchical include/exclude behavior
func TestDirectoryIncludeExcludeHierarchy(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	// Create nested structure
	framework.CreateFile("media/movies/action/2023/terminator.mp4", "action movie")
	framework.CreateFile("media/movies/comedy/2022/ghostbusters.mp4", "comedy movie")
	framework.CreateFile("media/tv/drama/series1/episode1.mp4", "drama episode")
	framework.Sync()

	// All directories should start as excluded
	framework.AssertDirectoryExcluded("media")
	framework.AssertDirectoryExcluded("media/movies")
	framework.AssertDirectoryExcluded("media/movies/action")
	framework.AssertDirectoryExcluded("media/movies/action/2023")

	// No files should be processed
	framework.AssertItemNotExists("media/movies/action/2023", "terminator.mp4")
	framework.AssertItemNotExists("media/movies/comedy/2022", "ghostbusters.mp4")
	framework.AssertItemNotExists("media/tv/drama/series1", "episode1.mp4")

	// User includes the entire media directory
	framework.IncludeDirectory("media")
	framework.Sync()

	// Only the media directory is included, subdirectories remain excluded
	framework.AssertDirectoryIncluded("media")
	framework.AssertDirectoryExcluded("media/movies")
	framework.AssertDirectoryExcluded("media/tv")

	// User includes movies subdirectory
	framework.IncludeDirectory("media/movies")
	framework.Sync()

	// Movies is included but its subdirectories are still excluded
	framework.AssertDirectoryIncluded("media/movies")
	framework.AssertDirectoryExcluded("media/movies/action")
	framework.AssertDirectoryExcluded("media/movies/comedy")

	// User includes entire action hierarchy
	framework.IncludeDirectory("media/movies/action")
	framework.Sync()
	framework.IncludeDirectory("media/movies/action/2023")
	framework.Sync()

	// Now action files are processed
	framework.AssertDirectoryIncluded("media/movies/action")
	framework.AssertDirectoryIncluded("media/movies/action/2023")
	framework.AssertItemExists("media/movies/action/2023", "terminator.mp4")

	// But comedy is still excluded
	framework.AssertDirectoryExcluded("media/movies/comedy")
	framework.AssertItemNotExists("media/movies/comedy/2022", "ghostbusters.mp4")
}

// TestDirectoryIncludeExcludeWithFileOperations tests include/exclude with file movements
func TestDirectoryIncludeExcludeWithFileOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	// Create files and include directories to process them
	framework.CreateFile("included/movie1.mp4", "movie 1")
	framework.CreateFile("excluded/movie2.mp4", "movie 2")
	framework.Sync()

	// Include only the "included" directory
	framework.IncludeDirectory("included")
	framework.Sync()

	// Verify processing
	framework.AssertDirectoryIncluded("included")
	framework.AssertDirectoryExcluded("excluded")
	framework.AssertItemExists("included", "movie1.mp4")
	framework.AssertItemNotExists("excluded", "movie2.mp4")

	// Move file from included to excluded directory
	framework.MoveFile("included/movie1.mp4", "excluded/movie1.mp4")
	framework.Sync()

	// Directory should be automatically included when file is moved there
	framework.AssertDirectoryIncluded("excluded")
	// File should be processed because destination directory is now included
	framework.AssertItemExists("excluded", "movie1.mp4")

	// Move file from excluded to included directory
	framework.MoveFile("excluded/movie2.mp4", "included/movie2.mp4")
	framework.Sync()

	// File should now be processed because destination is included
	framework.AssertItemExists("included", "movie2.mp4")

	// Include the excluded directory
	framework.IncludeDirectory("excluded")
	framework.Sync()

	// Now the file in excluded directory should be processed
	framework.AssertDirectoryIncluded("excluded")
	framework.AssertItemExists("excluded", "movie1.mp4")
}

// TestDirectoryIncludeExcludeWithAutoTags tests include/exclude with AutoTags
func TestDirectoryIncludeExcludeWithAutoTags(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	// Create files in directories
	framework.CreateFile("movies/action/terminator.mp4", "action movie")
	framework.CreateFile("movies/comedy/ghostbusters.mp4", "comedy movie")
	framework.Sync()

	// Include directories to process files
	framework.IncludeDirectory("movies")
	framework.Sync()
	framework.IncludeDirectory("movies/action")
	framework.Sync()
	framework.IncludeDirectory("movies/comedy")
	framework.Sync()

	// Verify files are processed with AutoTags
	framework.AssertItemExists("movies/action", "terminator.mp4")
	framework.AssertItemExists("movies/comedy", "ghostbusters.mp4")
	framework.AssertAutoTagsExist("movies/action", []string{"movies/action"})
	framework.AssertAutoTagsExist("movies/comedy", []string{"movies/comedy"})

	// Exclude action directory
	framework.ExcludeDirectory("movies/action")
	framework.Sync()

	// Action files should be removed from index and AutoTags cleaned up
	framework.AssertItemNotExists("movies/action", "terminator.mp4")

	// Comedy should still work
	framework.AssertItemExists("movies/comedy", "ghostbusters.mp4")
	framework.AssertAutoTagsExist("movies/comedy", []string{"movies/comedy"})

	// Re-include action directory
	framework.IncludeDirectory("movies/action")
	framework.Sync()

	// Files and AutoTags should be restored
	framework.AssertItemExists("movies/action", "terminator.mp4")
	framework.AssertAutoTagsExist("movies/action", []string{"movies/action"})
}

// TestDirectoryIncludeExcludeWithCustomAutoTags tests include/exclude with custom AutoTags
func TestDirectoryIncludeExcludeWithCustomAutoTags(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	// Create files in directories
	framework.CreateFile("movies/action/terminator.mp4", "action movie")
	framework.Sync()

	// Create custom tag and add to directory
	actionTag := &model.Tag{Title: "Action Genre"}
	err := framework.GetDatabase().CreateOrUpdateTag(context.Background(), actionTag)
	require.NoError(t, err)

	// Include directory first, then add custom tags
	framework.IncludeDirectory("movies")
	framework.Sync()
	framework.IncludeDirectory("movies/action")
	framework.AddCustomTagsToDirectory("movies/action", []*model.Tag{actionTag})
	framework.Sync()

	// Verify custom AutoTags are applied
	framework.AssertItemExists("movies/action", "terminator.mp4")
	framework.AssertCustomAutoTagExists("movies/action", "terminator.mp4", actionTag, "Action")

	// Exclude the directory
	framework.ExcludeDirectory("movies/action")
	framework.Sync()

	// File should be removed from index (custom AutoTags cleaned up too)
	framework.AssertItemNotExists("movies/action", "terminator.mp4")

	// Re-include the directory
	framework.IncludeDirectory("movies/action")
	framework.Sync()

	// Files and custom AutoTags should be restored
	framework.AssertItemExists("movies/action", "terminator.mp4")
	framework.AssertCustomAutoTagExists("movies/action", "terminator.mp4", actionTag, "Action")
}

// TestDirectoryIncludeOrCreate tests the IncludeOrCreateDirectory method
func TestDirectoryIncludeOrCreate(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	// Directory doesn't exist yet
	framework.AssertDirectoryNotExists("newdir")

	// IncludeOrCreateDirectory should create and include it
	framework.IncludeOrCreateDirectory("newdir")

	// Directory should now exist and be included
	framework.AssertDirectoryExists("newdir")
	framework.AssertDirectoryIncluded("newdir")

	// Create a file in the directory
	framework.CreateFile("newdir/movie.mp4", "test movie")
	framework.Sync()

	// File should be processed because directory is included
	framework.AssertItemExists("newdir", "movie.mp4")
}

// TestDirectoryExcludeReInclude tests the workflow of excluding and re-including directories
func TestDirectoryExcludeReInclude(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	// Create files and include directory
	framework.CreateFile("temp/file1.mp4", "file 1")
	framework.CreateFile("temp/file2.mp4", "file 2")
	framework.Sync()
	framework.IncludeDirectory("temp")
	framework.Sync()

	// Verify files are processed
	framework.AssertItemExists("temp", "file1.mp4")
	framework.AssertItemExists("temp", "file2.mp4")
	items := framework.GetItems("temp")
	assert.Len(t, items, 2, "Should have 2 items in temp directory")

	// Exclude the directory
	framework.ExcludeDirectory("temp")
	framework.Sync()

	// Files should be removed from index
	framework.AssertItemNotExists("temp", "file1.mp4")
	framework.AssertItemNotExists("temp", "file2.mp4")
	excludedItems := framework.GetItems("temp")
	assert.Len(t, excludedItems, 0, "Should have 0 items in excluded directory")

	// Add more files while excluded
	framework.CreateFile("temp/file3.mp4", "file 3")
	framework.Sync()

	// New file should not be processed
	framework.AssertItemNotExists("temp", "file3.mp4")

	// Re-include the directory
	framework.IncludeDirectory("temp")
	framework.Sync()

	// All files (including new ones) should be processed
	framework.AssertItemExists("temp", "file1.mp4")
	framework.AssertItemExists("temp", "file2.mp4")
	framework.AssertItemExists("temp", "file3.mp4")
	finalItems := framework.GetItems("temp")
	assert.Len(t, finalItems, 3, "Should have 3 items after re-including directory")
}
