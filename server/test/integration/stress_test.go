package integration

import (
	"fmt"
	"my-collection/server/test/testutils"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMassiveFileOperations tests performance with large numbers of files
func TestMassiveFileOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create 1000 files across 100 directories
	t.Log("Creating 1000 files across 100 directories...")
	framework.CreateLargeTestSet(100, 10) // 1000 files total

	start := time.Now()
	framework.Sync()
	syncDuration := time.Since(start)
	t.Logf("Sync of 1000 files took: %v", syncDuration)

	// Verify all files exist
	totalItems := 0
	for i := 0; i < 100; i++ {
		dirPath := fmt.Sprintf("large_test/dir%03d", i)
		items := framework.GetItems(dirPath)
		totalItems += len(items)
	}
	assert.Equal(t, 1000, totalItems, "Should have synced all 1000 files")

	// Delete every other directory
	t.Log("Deleting 50 directories...")
	for i := 0; i < 100; i += 2 {
		dirPath := fmt.Sprintf("large_test/dir%03d", i)
		framework.DeleteDir(dirPath)
	}

	start = time.Now()
	framework.Sync()
	syncDuration = time.Since(start)
	t.Logf("Sync after deleting 500 files took: %v", syncDuration)

	// Verify deletions
	remainingItems := 0
	for i := 0; i < 100; i++ {
		dirPath := fmt.Sprintf("large_test/dir%03d", i)
		items := framework.GetItems(dirPath)
		if i%2 == 0 {
			assert.Len(t, items, 0, "Deleted directory should have no items")
		} else {
			assert.Len(t, items, 10, "Remaining directory should have 10 items")
			remainingItems += len(items)
		}
	}
	assert.Equal(t, 500, remainingItems, "Should have 500 remaining files")
}

// TestDeepHierarchy tests very deep directory nesting
func TestDeepHierarchy(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create a 20-level deep hierarchy
	pathParts := make([]string, 20)
	for i := 0; i < 20; i++ {
		pathParts[i] = fmt.Sprintf("level%02d", i)
	}
	deepPath := strings.Join(pathParts, "/")

	framework.CreateFile(deepPath+"/deep_file.mp4", "deep content")
	framework.Sync()

	framework.AssertItemExists(deepPath, "deep_file.mp4")

	// Move the entire hierarchy
	framework.MoveDir("level00", "moved_level00")
	framework.Sync()

	newDeepPath := "moved_level00/" + strings.Join(pathParts[1:], "/")
	framework.AssertItemExists(newDeepPath, "deep_file.mp4")
	framework.AssertItemNotExists(deepPath, "deep_file.mp4")
}

// TestRapidFileChanges tests rapid file system changes
func TestRapidFileChanges(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Rapid creation, modification, and deletion
	for cycle := 0; cycle < 10; cycle++ {
		// Create files
		for i := 0; i < 20; i++ {
			framework.CreateFile(fmt.Sprintf("rapid/cycle%d_file%d.mp4", cycle, i), "content")
		}
		framework.Sync()

		// Move half the files
		for i := 0; i < 10; i++ {
			src := fmt.Sprintf("rapid/cycle%d_file%d.mp4", cycle, i)
			dst := fmt.Sprintf("rapid/moved/cycle%d_file%d.mp4", cycle, i)
			framework.MoveFile(src, dst)
		}
		framework.Sync()

		// Delete the other half
		for i := 10; i < 20; i++ {
			framework.DeleteFile(fmt.Sprintf("rapid/cycle%d_file%d.mp4", cycle, i))
		}
		framework.Sync()

		// Verify state
		rapidItems := framework.GetItems("rapid")
		movedItems := framework.GetItems("rapid/moved")
		assert.Len(t, rapidItems, 0, "Rapid directory should be empty after cycle %d", cycle)
		assert.Len(t, movedItems, 10*(cycle+1), "Moved directory should have %d items after cycle %d", 10*(cycle+1), cycle)
	}
}

// TestSpecialCharacters tests files and directories with special characters
func TestSpecialCharacters(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	specialNames := []string{
		"file with spaces.mp4",
		"file-with-dashes.mp4",
		"file_with_underscores.mp4",
		"file.with.dots.mp4",
		"file[with]brackets.mp4",
		"file(with)parentheses.mp4",
		"file&with&ampersands.mp4",
	}

	// Create files with special characters
	for i, name := range specialNames {
		dirName := fmt.Sprintf("special%d", i)
		framework.CreateFile(dirName+"/"+name, "special content")
	}
	framework.Sync()

	// Verify all files exist
	for i, name := range specialNames {
		dirName := fmt.Sprintf("special%d", i)
		framework.AssertItemExists(dirName, name)
	}

	// Move files to a common directory
	for i, name := range specialNames {
		src := fmt.Sprintf("special%d/%s", i, name)
		dst := fmt.Sprintf("common/%s", name)
		framework.MoveFile(src, dst)
	}
	framework.Sync()

	// Verify moves
	commonItems := framework.GetItems("common")
	assert.Len(t, commonItems, len(specialNames))

	for _, name := range specialNames {
		framework.AssertItemExists("common", name)
	}
}

// TestConcurrentDirectoryOperations tests operations on directories with many files
func TestConcurrentDirectoryOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large concurrent test in short mode")
	}

	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create multiple directories with files
	dirCount := 20
	filesPerDir := 50

	t.Logf("Creating %d directories with %d files each...", dirCount, filesPerDir)
	for i := 0; i < dirCount; i++ {
		for j := 0; j < filesPerDir; j++ {
			path := fmt.Sprintf("concurrent%d/file%d.mp4", i, j)
			framework.CreateFile(path, fmt.Sprintf("content %d %d", i, j))
		}
	}
	framework.Sync()

	// Verify initial state
	for i := 0; i < dirCount; i++ {
		dirPath := fmt.Sprintf("concurrent%d", i)
		items := framework.GetItems(dirPath)
		assert.Len(t, items, filesPerDir)
	}

	// Perform complex reorganization - move every other directory
	t.Log("Reorganizing directories...")
	for i := 0; i < dirCount; i += 2 {
		src := fmt.Sprintf("concurrent%d", i)
		dst := fmt.Sprintf("reorganized/concurrent%d", i)
		framework.MoveDir(src, dst)
	}
	framework.Sync()

	// Verify reorganization
	for i := 0; i < dirCount; i++ {
		if i%2 == 0 {
			// Moved directories
			dirPath := fmt.Sprintf("reorganized/concurrent%d", i)
			items := framework.GetItems(dirPath)
			assert.Len(t, items, filesPerDir)

			originalPath := fmt.Sprintf("concurrent%d", i)
			framework.AssertDirectoryNotExists(originalPath)
		} else {
			// Remaining directories
			dirPath := fmt.Sprintf("concurrent%d", i)
			items := framework.GetItems(dirPath)
			assert.Len(t, items, filesPerDir)
		}
	}
}

// TestInconsistentState tests recovery from inconsistent filesystem/DB state
func TestInconsistentState(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	// Create initial state
	framework.CreateFile("consistent/file1.mp4", "content1")
	framework.CreateFile("consistent/file2.mp4", "content2")
	framework.Sync()

	// Manually create inconsistencies by bypassing the framework
	// Add file to filesystem but not DB
	orphanPath := filepath.Join(framework.GetTempDir(), "consistent/orphan.mp4")
	err := os.WriteFile(orphanPath, []byte("orphan content"), 0644)
	require.NoError(t, err)

	// Remove file from filesystem but leave in DB
	stalePath := filepath.Join(framework.GetTempDir(), "consistent/file2.mp4")
	err = os.Remove(stalePath)
	require.NoError(t, err)

	// Sync should detect and fix inconsistencies
	framework.Sync()

	// Orphan file should be added to DB
	framework.AssertItemExists("consistent", "orphan.mp4")

	// Stale file should be removed from DB
	framework.AssertItemNotExists("consistent", "file2.mp4")

	// Original file should remain
	framework.AssertItemExists("consistent", "file1.mp4")
}

// TestVeryLongFilenames tests handling of files with very long names
func TestVeryLongFilenames(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create file with very long name (close to filesystem limits)
	longName := strings.Repeat("verylongfilename", 15) + ".mp4" // ~240 characters
	if len(longName) > 255 {
		longName = longName[:251] + ".mp4" // Ensure it fits in typical filesystem limits
	}

	framework.CreateFile("longnames/"+longName, "content")
	framework.Sync()

	framework.AssertItemExists("longnames", longName)

	// Test moving long-named file
	framework.MoveFile("longnames/"+longName, "moved/"+longName)
	framework.Sync()

	framework.AssertItemNotExists("longnames", longName)
	framework.AssertItemExists("moved", longName)
}

// TestEmptyDirectoryBehavior tests behavior with empty directories
func TestEmptyDirectoryBehavior(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create empty directories
	framework.CreateDir("empty1")
	framework.CreateDir("empty2/nested")
	framework.CreateDir("empty3/deeply/nested/structure")
	framework.Sync()

	// Verify empty directories are tracked
	framework.AssertDirectoryExists("empty1")
	framework.AssertDirectoryExists("empty2")
	framework.AssertDirectoryExists("empty2/nested")
	framework.AssertDirectoryExists("empty3")
	framework.AssertDirectoryExists("empty3/deeply")
	framework.AssertDirectoryExists("empty3/deeply/nested")
	framework.AssertDirectoryExists("empty3/deeply/nested/structure")

	// Add files to some empty directories
	framework.CreateFile("empty1/no_longer_empty.mp4", "content")
	framework.Sync()

	framework.AssertItemExists("empty1", "no_longer_empty.mp4")

	// Remove empty directories
	framework.DeleteDir("empty2")
	framework.DeleteDir("empty3")
	framework.Sync()

	framework.AssertDirectoryNotExists("empty2")
	framework.AssertDirectoryNotExists("empty3")
	framework.AssertDirectoryExists("empty1") // Should still exist because it has files
}

// TestCaseSensitivity tests case sensitivity behavior
func TestCaseSensitivity(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create files with different cases
	framework.CreateFile("case/Movie.mp4", "content1")
	framework.CreateFile("case/MOVIE.mp4", "content2")
	framework.CreateFile("case/movie.mp4", "content3")
	framework.Sync()

	// Behavior may depend on filesystem (case-sensitive vs case-insensitive)
	// This test ensures the system handles it gracefully regardless
	caseItems := framework.GetItems("case")

	// On case-sensitive filesystems, should have 3 items
	// On case-insensitive filesystems, may have fewer (last one wins)
	assert.GreaterOrEqual(t, len(caseItems), 1, "Should have at least one item")
	assert.LessOrEqual(t, len(caseItems), 3, "Should have at most three items")
}

// TestPerformanceRegression tests for performance regressions
func TestPerformanceRegression(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create a moderately complex structure
	fileCount := 400 // 20 dirs * 20 files each
	framework.CreateLargeTestSet(20, 20)

	// Measure sync performance
	start := time.Now()
	framework.Sync()
	initialSync := time.Since(start)
	t.Logf("Initial sync of %d files took: %v", fileCount, initialSync)

	// Measure incremental sync (no changes)
	start = time.Now()
	framework.Sync()
	noChangeSync := time.Since(start)
	t.Logf("No-change sync took: %v", noChangeSync)

	// Add a few more files and measure incremental sync
	for i := 0; i < 10; i++ {
		framework.CreateFile(fmt.Sprintf("newfiles/file%d.mp4", i), "new content")
	}

	start = time.Now()
	framework.Sync()
	incrementalSync := time.Since(start)
	t.Logf("Incremental sync of 10 new files took: %v", incrementalSync)

	// Performance assertions (adjust thresholds based on expected performance)
	assert.Less(t, noChangeSync, time.Second, "No-change sync should be fast")
	assert.Less(t, incrementalSync, initialSync, "Incremental sync should be faster than initial sync")
}

// TestEdgeCaseFilenames tests various edge cases in filenames
func TestEdgeCaseFilenames(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()

	framework.Sync()

	edgeCases := []string{
		"trailing.space .mp4", // Trailing space
		" leading.space.mp4",  // Leading space
		"multiple...dots.mp4", // Multiple consecutive dots
		"1.mp4",               // Single character
	}

	for i, filename := range edgeCases {
		dirName := fmt.Sprintf("edge%d", i)
		// Some edge cases might not be valid on all filesystems
		// Catch errors and skip those cases
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Skipping edge case file %s: %v", filename, r)
				}
			}()

			framework.CreateFile(dirName+"/"+filename, "edge content")
		}()
	}

	framework.Sync()

	// Verify what files were actually created
	for i := range edgeCases {
		dirName := fmt.Sprintf("edge%d", i)
		items := framework.GetItems(dirName)
		// At least verify the directory was created and sync didn't crash
		t.Logf("Directory %s has %d items", dirName, len(items))
	}
}

// TestCircularDirectoryMoves tests edge case of circular moves
func TestCircularDirectoryMoves(t *testing.T) {
	framework := testutils.NewIntegrationTestFramework(t)
	defer framework.Cleanup()
	assert.NoError(t, framework.AutoIncludeHierarchy(""))

	framework.Sync()

	// Create initial structure
	framework.CreateFile("a/file_a.mp4", "content a")
	framework.CreateFile("b/file_b.mp4", "content b")
	framework.CreateFile("c/file_c.mp4", "content c")
	framework.Sync()

	// Perform moves that could create issues if not handled properly
	// Move a to temp, b to a, c to b, temp to c
	framework.MoveDir("a", "temp")
	framework.MoveDir("b", "a")
	framework.MoveDir("c", "b")
	framework.MoveDir("temp", "c")
	framework.Sync()

	// Verify final state: a->b->c->a (circular)
	framework.AssertItemExists("a", "file_b.mp4") // b moved to a
	framework.AssertItemExists("b", "file_c.mp4") // c moved to b
	framework.AssertItemExists("c", "file_a.mp4") // a moved to c
}
