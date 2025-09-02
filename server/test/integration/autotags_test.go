package integration

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
	"my-collection/server/test/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAutoTagsBasicBehavior tests basic AutoTags functionality
func TestAutoTagsBasicBehavior(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create files in different directories
	framework.CreateFile("movies/action/terminator.mp4", "content1")
	framework.CreateFile("movies/comedy/ghostbusters.mp4", "content2")
	framework.CreateFile("tv/drama/breaking_bad_s01e01.mp4", "content3")
	framework.Sync()

	// Verify AutoTags are correctly applied
	framework.AssertAutoTagsExist("movies/action", []string{"movies/action"})
	framework.AssertAutoTagsExist("movies/comedy", []string{"movies/comedy"})
	framework.AssertAutoTagsExist("tv/drama", []string{"tv/drama"})

	// Verify items have the correct AutoTags
	actionItem, err := framework.GetDirectoryItemsGetter().GetBelongingItem("movies/action", "terminator.mp4")
	require.NoError(t, err)
	require.NotNil(t, actionItem)

	hasActionTag := false
	for _, tag := range actionItem.Tags {
		if tag.Title == "movies/action" && tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			hasActionTag = true
			break
		}
	}
	assert.True(t, hasActionTag, "Item should have AutoTag 'movies/action'")
}

// TestAutoTagsFileMovement tests AutoTags when files are moved between directories
func TestAutoTagsFileMovement(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create file in source directory
	framework.CreateFile("source/movie.mp4", "movie content")
	framework.CreateDir("destination")
	framework.Sync()

	// Verify initial AutoTags
	framework.AssertAutoTagsExist("source", []string{"source"})

	// Move file to destination
	framework.MoveFile("source/movie.mp4", "destination/movie.mp4")
	framework.Sync()

	// Verify AutoTags are updated
	movedItem, err := framework.GetDirectoryItemsGetter().GetBelongingItem("destination", "movie.mp4")
	require.NoError(t, err)
	require.NotNil(t, movedItem)

	hasDestinationTag := false
	hasSourceTag := false
	for _, tag := range movedItem.Tags {
		if tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			if tag.Title == "destination" {
				hasDestinationTag = true
			}
			if tag.Title == "source" {
				hasSourceTag = true
			}
		}
	}

	assert.True(t, hasDestinationTag, "Moved item should have destination AutoTag")
	assert.False(t, hasSourceTag, "Moved item should not have source AutoTag")
}

// TestAutoTagsDirectoryRename tests AutoTags when directories are renamed
func TestAutoTagsDirectoryRename(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create files in directory
	framework.CreateFile("oldname/movie1.mp4", "content1")
	framework.CreateFile("oldname/movie2.mp4", "content2")
	framework.Sync()

	// Verify initial AutoTags
	framework.AssertAutoTagsExist("oldname", []string{"oldname"})

	// Rename directory
	framework.MoveDir("oldname", "newname")
	framework.Sync()

	// Verify AutoTags are updated for renamed directory
	framework.AssertAutoTagsExist("newname", []string{"newname"})

	// Verify individual items have updated AutoTags
	for _, filename := range []string{"movie1.mp4", "movie2.mp4"} {
		item, err := framework.GetDirectoryItemsGetter().GetBelongingItem("newname", filename)
		require.NoError(t, err)
		require.NotNil(t, item)

		hasNewTag := false
		hasOldTag := false
		for _, tag := range item.Tags {
			if tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
				if tag.Title == "newname" {
					hasNewTag = true
				}
				if tag.Title == "oldname" {
					hasOldTag = true
				}
			}
		}

		assert.True(t, hasNewTag, "Item %s should have new AutoTag 'newname'", filename)
		assert.False(t, hasOldTag, "Item %s should not have old AutoTag 'oldname'", filename)
	}
}

// TestAutoTagsNestedDirectoryOperations tests AutoTags in nested directory structures
func TestAutoTagsNestedDirectoryOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create nested structure
	framework.CreateFile("media/movies/action/2023/terminator.mp4", "content1")
	framework.CreateFile("media/movies/comedy/2022/ghostbusters.mp4", "content2")
	framework.Sync()

	// Verify nested AutoTags
	framework.AssertAutoTagsExist("media/movies/action/2023", []string{"media/movies/action/2023"})
	framework.AssertAutoTagsExist("media/movies/comedy/2022", []string{"media/movies/comedy/2022"})

	// Move entire media directory
	framework.MoveDir("media", "content")
	framework.Sync()

	// Verify AutoTags are updated throughout the hierarchy
	framework.AssertAutoTagsExist("content/movies/action/2023", []string{"content/movies/action/2023"})
	framework.AssertAutoTagsExist("content/movies/comedy/2022", []string{"content/movies/comedy/2022"})

	// Verify files have correct AutoTags
	terminatorItem, err := framework.GetDirectoryItemsGetter().GetBelongingItem("content/movies/action/2023", "terminator.mp4")
	require.NoError(t, err)
	require.NotNil(t, terminatorItem)

	hasCorrectTag := false
	for _, tag := range terminatorItem.Tags {
		if tag.Title == "content/movies/action/2023" && tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			hasCorrectTag = true
			break
		}
	}
	assert.True(t, hasCorrectTag, "Nested item should have correct AutoTag after directory move")
}

// TestAutoTagsDirectoryDeletion tests AutoTags cleanup when directories are deleted
func TestAutoTagsDirectoryDeletion(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create directory with files
	framework.CreateFile("temporary/file1.mp4", "content1")
	framework.CreateFile("temporary/file2.mp4", "content2")
	framework.Sync()

	// Verify AutoTags exist
	framework.AssertAutoTagsExist("temporary", []string{"temporary"})

	// Get initial tag count
	initialTags := framework.GetTags()
	temporaryTagExists := false
	for _, tag := range initialTags {
		if tag.Title == "temporary" && tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			temporaryTagExists = true
			break
		}
	}
	assert.True(t, temporaryTagExists, "Temporary AutoTag should exist initially")

	// Delete directory
	framework.DeleteDir("temporary")
	framework.Sync()

	// Verify AutoTag is removed
	finalTags := framework.GetTags()
	temporaryTagExists = false
	for _, tag := range finalTags {
		if tag.Title == "temporary" && tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			temporaryTagExists = true
			break
		}
	}
	assert.False(t, temporaryTagExists, "Temporary AutoTag should be removed after directory deletion")
}

// TestAutoTagsMultipleFileOperations tests AutoTags during complex multi-file operations
func TestAutoTagsMultipleFileOperations(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create multiple files in various directories
	files := map[string]string{
		"action/movie1.mp4":   "content1",
		"action/movie2.mp4":   "content2",
		"comedy/movie3.mp4":   "content3",
		"drama/movie4.mp4":    "content4",
		"thriller/movie5.mp4": "content5",
	}

	for path, content := range files {
		framework.CreateFile(path, content)
	}
	framework.Sync()

	// Move all action and comedy movies to a "favorites" directory
	framework.MoveFile("action/movie1.mp4", "favorites/movie1.mp4")
	framework.MoveFile("action/movie2.mp4", "favorites/movie2.mp4")
	framework.MoveFile("comedy/movie3.mp4", "favorites/movie3.mp4")
	framework.Sync()

	// Verify AutoTags for moved files
	for _, filename := range []string{"movie1.mp4", "movie2.mp4", "movie3.mp4"} {
		item, err := framework.GetDirectoryItemsGetter().GetBelongingItem("favorites", filename)
		require.NoError(t, err)
		require.NotNil(t, item)

		hasFavoritesTag := false
		hasOldTag := false
		for _, tag := range item.Tags {
			if tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
				if tag.Title == "favorites" {
					hasFavoritesTag = true
				}
				if tag.Title == "action" || tag.Title == "comedy" {
					hasOldTag = true
				}
			}
		}

		assert.True(t, hasFavoritesTag, "File %s should have 'favorites' AutoTag", filename)
		assert.False(t, hasOldTag, "File %s should not have old AutoTags", filename)
	}

	// Verify remaining files still have correct AutoTags
	dramaItem, err := framework.GetDirectoryItemsGetter().GetBelongingItem("drama", "movie4.mp4")
	require.NoError(t, err)
	require.NotNil(t, dramaItem)

	hasDramaTag := false
	for _, tag := range dramaItem.Tags {
		if tag.Title == "drama" && tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			hasDramaTag = true
			break
		}
	}
	assert.True(t, hasDramaTag, "Remaining drama file should still have drama AutoTag")
}

// TestAutoTagsWithManualTags tests interaction between AutoTags and manually assigned tags
func TestAutoTagsWithManualTags(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create a file with AutoTags
	framework.CreateFile("movies/action/terminator.mp4", "content")
	framework.Sync()

	// Get the item and verify it has AutoTags
	item, err := framework.GetDirectoryItemsGetter().GetBelongingItem("movies/action", "terminator.mp4")
	require.NoError(t, err)
	require.NotNil(t, item)

	autoTagCount := 0
	for _, tag := range item.Tags {
		if tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			autoTagCount++
		}
	}
	assert.Greater(t, autoTagCount, 0, "Item should have AutoTags")

	// Create a manual tag (not under directories parent)
	manualTag := &model.Tag{
		Title:    "sci-fi",
		ParentID: nil, // Manual tag has no parent
	}
	err = framework.GetDatabase().CreateOrUpdateTag(manualTag)
	require.NoError(t, err)

	// Add manual tag to item (simulating user action)
	item.Tags = append(item.Tags, manualTag)
	err = framework.GetDatabase().UpdateItem(item)
	require.NoError(t, err)

	// Move file to different directory
	framework.MoveFile("movies/action/terminator.mp4", "movies/sci-fi/terminator.mp4")
	framework.Sync()

	// Verify both manual tag and new AutoTag coexist
	movedItem, err := framework.GetDirectoryItemsGetter().GetBelongingItem("movies/sci-fi", "terminator.mp4")
	require.NoError(t, err)
	require.NotNil(t, movedItem)

	hasSciFiAutoTag := false
	hasSciFiManualTag := false
	hasOldAutoTag := false

	for _, tag := range movedItem.Tags {
		if tag.Title == "sci-fi" {
			if tag.ParentID == nil {
				hasSciFiManualTag = true
			}
		}
		if tag.Title == "movies/sci-fi" && tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			hasSciFiAutoTag = true
		}
		if tag.Title == "movies/action" && tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			hasOldAutoTag = true
		}
	}

	assert.True(t, hasSciFiAutoTag, "Item should have new AutoTag")
	assert.True(t, hasSciFiManualTag, "Item should retain manual tag")
	assert.False(t, hasOldAutoTag, "Item should not have old AutoTag")
}

// TestAutoTagsStaleCleanup tests that stale AutoTags are properly cleaned up
func TestAutoTagsStaleCleanup(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create files that will generate AutoTags
	framework.CreateFile("temp1/file1.mp4", "content1")
	framework.CreateFile("temp2/file2.mp4", "content2")
	framework.CreateFile("permanent/file3.mp4", "content3")
	framework.Sync()

	// Count initial AutoTags
	initialTags := framework.GetTags()
	autoTagCount := 0
	for _, tag := range initialTags {
		if tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			autoTagCount++
		}
	}

	// Remove temp directories externally (simulating external filesystem changes)
	framework.DeleteDir("temp1")
	framework.DeleteDir("temp2")
	framework.Sync()

	// Count final AutoTags
	finalTags := framework.GetTags()
	finalAutoTagCount := 0
	remainingTagTitles := make([]string, 0)
	for _, tag := range finalTags {
		if tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
			finalAutoTagCount++
			remainingTagTitles = append(remainingTagTitles, tag.Title)
		}
	}

	// Should have fewer AutoTags now
	assert.Less(t, finalAutoTagCount, autoTagCount, "Should have fewer AutoTags after cleanup")

	// Permanent directory AutoTag should still exist
	assert.Contains(t, remainingTagTitles, "permanent", "Permanent directory AutoTag should remain")

	// Temp AutoTags should be gone
	assert.NotContains(t, remainingTagTitles, "temp1", "temp1 AutoTag should be removed")
	assert.NotContains(t, remainingTagTitles, "temp2", "temp2 AutoTag should be removed")
}

// TestAutoTagsRealLibraryScenario tests AutoTags in a realistic media library scenario
func TestAutoTagsRealLibraryScenario(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	// Create realistic library with AutoTags behavior
	framework.CreateTestLibrary()
	framework.Sync()

	// Verify AutoTags are created for the library structure
	expectedAutoTags := map[string][]string{
		"movies/action/2023":        {"movies/action/2023"},
		"movies/action/2022":        {"movies/action/2022"},
		"movies/comedy/2023":        {"movies/comedy/2023"},
		"movies/drama/2022":         {"movies/drama/2022"},
		"tv/drama/breaking_bad/s01": {"tv/drama/breaking_bad/s01"},
		"tv/comedy/office/s01":      {"tv/comedy/office/s01"},
		"documentaries/nature":      {"documentaries/nature"},
		"music/rock/album1":         {"music/rock/album1"},
		"music/jazz/album2":         {"music/jazz/album2"},
	}

	for dir, expectedTags := range expectedAutoTags {
		framework.AssertAutoTagsExist(dir, expectedTags)
	}

	// Simulate library reorganization
	// Move all action movies to an "action-collection" directory
	framework.MoveFile("movies/action/2023/terminator.mp4", "action-collection/2023/terminator.mp4")
	framework.MoveFile("movies/action/2022/john_wick.mp4", "action-collection/2022/john_wick.mp4")
	framework.Sync()

	// Verify AutoTags are updated correctly
	framework.AssertAutoTagsExist("action-collection/2023", []string{"action-collection/2023"})
	framework.AssertAutoTagsExist("action-collection/2022", []string{"action-collection/2022"})

	// Verify original AutoTags are cleaned up if directories are empty
	// (This depends on implementation - empty directories may be removed)
}

// === CUSTOM AUTOTAGS TESTS ===
// These test the custom tags that users can manually assign to directories
// which then automatically get applied to all files in those directories

// TestCustomAutoTagsBasicBehavior tests basic custom AutoTags functionality
func TestCustomAutoTagsBasicBehavior(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create some files in directories
	framework.CreateFile("movies/action/terminator.mp4", "action movie")
	framework.CreateFile("movies/comedy/ghostbusters.mp4", "comedy movie")
	framework.Sync()

	// Create custom tags that users would assign to directories
	actionTag := &model.Tag{Title: "Action Genre"}
	err := framework.GetDatabase().CreateOrUpdateTag(actionTag)
	require.NoError(t, err)

	year2023Tag := &model.Tag{Title: "2023 Release"}
	err = framework.GetDatabase().CreateOrUpdateTag(year2023Tag)
	require.NoError(t, err)

	// Add custom tags to the action directory
	framework.AddCustomTagsToDirectory("movies/action", []*model.Tag{actionTag, year2023Tag})
	framework.Sync()

	// Verify that files in the action directory now have custom AutoTags
	terminatorItem, err := framework.GetDirectoryItemsGetter().GetBelongingItem("movies/action", "terminator.mp4")
	require.NoError(t, err)
	require.NotNil(t, terminatorItem)

	// Check for custom AutoTags: should have "Action" (child of "Action Genre") and "Action" (child of "2023 Release")
	hasActionAutoTag := false
	hasYearAutoTag := false
	for _, tag := range terminatorItem.Tags {
		if tag.ParentID != nil {
			if *tag.ParentID == actionTag.Id && tag.Title == "Action" {
				hasActionAutoTag = true
			}
			if *tag.ParentID == year2023Tag.Id && tag.Title == "Action" {
				hasYearAutoTag = true
			}
		}
	}

	assert.True(t, hasActionAutoTag, "File should have custom AutoTag derived from Action Genre")
	assert.True(t, hasYearAutoTag, "File should have custom AutoTag derived from 2023 Release")

	// Files in other directories should not have these custom AutoTags
	ghostbustersItem, err := framework.GetDirectoryItemsGetter().GetBelongingItem("movies/comedy", "ghostbusters.mp4")
	require.NoError(t, err)
	require.NotNil(t, ghostbustersItem)

	hasCustomTag := false
	for _, tag := range ghostbustersItem.Tags {
		if tag.ParentID != nil && (*tag.ParentID == actionTag.Id || *tag.ParentID == year2023Tag.Id) {
			hasCustomTag = true
		}
	}
	assert.False(t, hasCustomTag, "Files in other directories should not have custom AutoTags")
}

// TestCustomAutoTagsWithMultipleFiles tests custom AutoTags with multiple files in same directory
func TestCustomAutoTagsWithMultipleFiles(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create multiple files in same directory
	framework.CreateFile("collection/movie1.mp4", "movie 1")
	framework.CreateFile("collection/movie2.mp4", "movie 2")
	framework.CreateFile("collection/movie3.mp4", "movie 3")
	framework.Sync()

	// Create custom tags
	favoritesTag := &model.Tag{Title: "Favorites"}
	hdTag := &model.Tag{Title: "HD Quality"}
	err := framework.GetDatabase().CreateOrUpdateTag(favoritesTag)
	require.NoError(t, err)
	err = framework.GetDatabase().CreateOrUpdateTag(hdTag)
	require.NoError(t, err)

	// Add multiple custom tags to directory
	framework.AddCustomTagsToDirectory("collection", []*model.Tag{favoritesTag, hdTag})
	framework.Sync()

	// Verify all files get all custom AutoTags
	files := []string{"movie1.mp4", "movie2.mp4", "movie3.mp4"}
	for _, filename := range files {
		framework.AssertCustomAutoTagExists("collection", filename, favoritesTag, "Collection")
		framework.AssertCustomAutoTagExists("collection", filename, hdTag, "Collection")
	}

	// Add new file to directory
	framework.CreateFile("collection/movie4.mp4", "movie 4")
	framework.Sync()

	// Verify new file automatically gets the custom AutoTags
	framework.AssertCustomAutoTagExists("collection", "movie4.mp4", favoritesTag, "Collection")
	framework.AssertCustomAutoTagExists("collection", "movie4.mp4", hdTag, "Collection")
}
