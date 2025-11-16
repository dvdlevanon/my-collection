package testutils

import (
	"fmt"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/fssync"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IntegrationTestFramework provides utilities for integration testing with real DB and filesystem
type IntegrationTestFramework struct {
	t           *testing.T
	tempDir     string
	dbFile      string
	database    db.Database
	fsManager   *fssync.FsManager
	dig         model.DirectoryItemsGetter
	initialSync bool
}

type testFileFilter struct {
}

func (f testFileFilter) Filter(path string) bool {
	// Accept all files except hidden ones for testing
	base := filepath.Base(path)
	return !strings.HasPrefix(base, ".")
}

// NewIntegrationTestFramework creates a new test framework with fresh DB and filesystem
func NewIntegrationTestFramework(t *testing.T) *IntegrationTestFramework {
	tempDir, err := os.MkdirTemp("", "integration-test-*")
	require.NoError(t, err)
	rootDir := filepath.Join(tempDir, "root")
	require.NoError(t, os.Mkdir(rootDir, 0755))

	dbFile := filepath.Join(tempDir, "test.db")
	database, err := db.New(dbFile, false)
	require.NoError(t, err)

	// Initialize relativasor with our test directory
	relativasor.Init(rootDir)

	// Initialize directories and tags subsystems
	err = directories.Init(database)
	require.NoError(t, err)

	// Create FsManager with filter that accepts all files
	fsManager, err := fssync.NewFsManager(database, testFileFilter{}, time.Hour) // Long interval since we sync manually
	require.NoError(t, err)

	return &IntegrationTestFramework{
		t:         t,
		tempDir:   rootDir,
		dbFile:    dbFile,
		database:  database,
		fsManager: fsManager,
		dig:       fsManager,
	}
}

// Cleanup removes temporary files and closes database
func (f *IntegrationTestFramework) Cleanup() {
	if f.database != nil {
		// Close database connection if possible
		// Note: gorm doesn't expose Close directly, but it will be closed when the process ends
	}
	if f.tempDir != "" {
		os.RemoveAll(f.tempDir)
	}
}

// GetTempDir returns the temporary directory path for this test
func (f *IntegrationTestFramework) GetTempDir() string {
	return f.tempDir
}

// GetDatabase returns the test database instance
func (f *IntegrationTestFramework) GetDatabase() db.Database {
	return f.database
}

// GetFsManager returns the filesystem manager instance
func (f *IntegrationTestFramework) GetDirectoryItemsGetter() model.DirectoryItemsGetter {
	return f.dig
}

// CreateFile creates a file with given content at the specified relative path
func (f *IntegrationTestFramework) CreateFile(relativePath string, content string) {
	fullPath := filepath.Join(f.tempDir, relativePath)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	require.NoError(f.t, err)

	err = os.WriteFile(fullPath, []byte(content), 0644)
	require.NoError(f.t, err)
}

// CreateDir creates a directory at the specified relative path
func (f *IntegrationTestFramework) CreateDir(relativePath string) {
	fullPath := filepath.Join(f.tempDir, relativePath)
	err := os.MkdirAll(fullPath, 0755)
	require.NoError(f.t, err)
}

// DeleteFile deletes a file at the specified relative path
func (f *IntegrationTestFramework) DeleteFile(relativePath string) {
	fullPath := filepath.Join(f.tempDir, relativePath)
	err := os.Remove(fullPath)
	require.NoError(f.t, err)
}

// DeleteDir deletes a directory at the specified relative path
func (f *IntegrationTestFramework) DeleteDir(relativePath string) {
	fullPath := filepath.Join(f.tempDir, relativePath)
	err := os.RemoveAll(fullPath)
	require.NoError(f.t, err)
}

// MoveFile moves a file from src to dst (both relative paths)
func (f *IntegrationTestFramework) MoveFile(src, dst string) {
	srcPath := filepath.Join(f.tempDir, src)
	dstPath := filepath.Join(f.tempDir, dst)

	// Ensure destination directory exists
	err := os.MkdirAll(filepath.Dir(dstPath), 0755)
	require.NoError(f.t, err)

	err = os.Rename(srcPath, dstPath)
	require.NoError(f.t, err)
}

// MoveDir moves a directory from src to dst (both relative paths)
func (f *IntegrationTestFramework) MoveDir(src, dst string) {
	srcPath := filepath.Join(f.tempDir, src)
	dstPath := filepath.Join(f.tempDir, dst)

	// Ensure parent of destination directory exists
	err := os.MkdirAll(filepath.Dir(dstPath), 0755)
	require.NoError(f.t, err)

	err = os.Rename(srcPath, dstPath)
	require.NoError(f.t, err)
}

// Sync performs a filesystem synchronization
func (f *IntegrationTestFramework) Sync() {
	err := f.fsManager.Sync()
	require.NoError(f.t, err)
	f.initialSync = true
}

// includeAllDirectories makes all directories non-excluded for testing
// Returns (hadExcluded, error) where hadExcluded indicates if any directories were included

// GetItems returns all items for a given directory path
func (f *IntegrationTestFramework) GetItems(dirPath string) []model.Item {
	items, err := f.fsManager.GetBelongingItems(directories.NormalizeDirectoryPath(dirPath))
	require.NoError(f.t, err)
	if items == nil {
		return []model.Item{}
	}
	return *items
}

// GetDirectories returns all directories from the database
func (f *IntegrationTestFramework) GetDirectories() []model.Directory {
	dirs, err := f.database.GetAllDirectories()
	require.NoError(f.t, err)
	return *dirs
}

// GetTags returns all tags from the database
func (f *IntegrationTestFramework) GetTags() []model.Tag {
	tags, err := f.database.GetAllTags()
	require.NoError(f.t, err)
	return *tags
}

// AssertItemExists checks that an item with the given title exists in the given directory
func (f *IntegrationTestFramework) AssertItemExists(dirPath, filename string) {
	item, err := f.fsManager.GetBelongingItem(directories.NormalizeDirectoryPath(dirPath), filename)
	require.NoError(f.t, err)
	assert.NotNil(f.t, item, "Item %s should exist in directory %s", filename, dirPath)
}

// AssertItemNotExists checks that an item with the given title does not exist in the given directory
func (f *IntegrationTestFramework) AssertItemNotExists(dirPath, filename string) {
	item, err := f.fsManager.GetBelongingItem(directories.NormalizeDirectoryPath(dirPath), filename)
	require.NoError(f.t, err)
	assert.Nil(f.t, item, "Item %s should not exist in directory %s", filename, dirPath)
}

// AssertDirectoryExists checks that a directory exists in the database
func (f *IntegrationTestFramework) AssertDirectoryExists(path string) {
	dir, err := f.database.GetDirectory("path = ?", directories.NormalizeDirectoryPath(path))
	require.NoError(f.t, err)
	assert.NotNil(f.t, dir, "Directory %s should exist in database", path)
}

// AssertDirectoryNotExists checks that a directory does not exist in the database
func (f *IntegrationTestFramework) AssertDirectoryNotExists(path string) {
	dir, err := f.database.GetDirectory("path = ?", directories.NormalizeDirectoryPath(path))
	if err != nil {
		// Record not found is expected
		return
	}
	assert.Nil(f.t, dir, "Directory %s should not exist in database", path)
}

// AssertAutoTagsExist checks that AutoTags are properly applied to items in a directory
func (f *IntegrationTestFramework) AssertAutoTagsExist(dirPath string, expectedAutoTagTitles []string) {
	items := f.GetItems(dirPath)
	for _, item := range items {
		itemAutoTags := make([]string, 0)
		for _, tag := range item.Tags {
			// Check if this is an AutoTag (directory-based tag)
			if tag.ParentID != nil && *tag.ParentID == directories.GetDirectoriesTagId() {
				itemAutoTags = append(itemAutoTags, tag.Title)
			}
		}

		for _, expectedTag := range expectedAutoTagTitles {
			assert.Contains(f.t, itemAutoTags, expectedTag,
				"Item %s in directory %s should have AutoTag %s", item.Title, dirPath, expectedTag)
		}
	}
}

// CreateTestLibrary creates a realistic test library structure
func (f *IntegrationTestFramework) CreateTestLibrary() {
	// Create a realistic media library structure
	structure := map[string]string{
		"movies/action/2023/terminator.mp4":            "action movie content",
		"movies/action/2022/john_wick.mp4":             "action movie content",
		"movies/comedy/2023/ghostbusters.mp4":          "comedy movie content",
		"movies/drama/2022/the_batman.mp4":             "drama movie content",
		"tv/drama/breaking_bad/s01/e01.mp4":            "tv episode content",
		"tv/drama/breaking_bad/s01/e02.mp4":            "tv episode content",
		"tv/comedy/office/s01/e01.mp4":                 "tv episode content",
		"documentaries/nature/planet_earth_s01e01.mp4": "documentary content",
		"music/rock/album1/song1.mp3":                  "music content",
		"music/jazz/album2/song2.mp3":                  "music content",
	}

	for path, content := range structure {
		f.CreateFile(path, content)
	}
}

// AssertTestLibraryExists verifies the test library was created correctly
func (f *IntegrationTestFramework) AssertTestLibraryExists() {
	expectedItems := map[string][]string{
		"movies/action/2023":        {"terminator.mp4"},
		"movies/action/2022":        {"john_wick.mp4"},
		"movies/comedy/2023":        {"ghostbusters.mp4"},
		"movies/drama/2022":         {"the_batman.mp4"},
		"tv/drama/breaking_bad/s01": {"e01.mp4", "e02.mp4"},
		"tv/comedy/office/s01":      {"e01.mp4"},
		"documentaries/nature":      {"planet_earth_s01e01.mp4"},
		"music/rock/album1":         {"song1.mp3"},
		"music/jazz/album2":         {"song2.mp3"},
	}

	for dir, files := range expectedItems {
		for _, file := range files {
			f.AssertItemExists(dir, file)
		}
	}
}

// WaitForOperation waits for an operation to complete (useful for async operations)
func (f *IntegrationTestFramework) WaitForOperation(timeout time.Duration, operation func() bool) bool {
	start := time.Now()
	for time.Since(start) < timeout {
		if operation() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// CreateLargeTestSet creates a large number of files for stress testing
func (f *IntegrationTestFramework) CreateLargeTestSet(dirCount, filesPerDir int) {
	for i := 0; i < dirCount; i++ {
		for j := 0; j < filesPerDir; j++ {
			path := fmt.Sprintf("large_test/dir%03d/file%03d.mp4", i, j)
			content := fmt.Sprintf("test content %d_%d", i, j)
			f.CreateFile(path, content)
		}
	}
}

// GetAllItems returns all items from the database
func (f *IntegrationTestFramework) GetAllItems() []model.Item {
	items, err := f.database.GetAllItems()
	require.NoError(f.t, err)
	if items == nil {
		return []model.Item{}
	}
	return *items
}

// CountItemsInDirectory counts all items in a directory and its subdirectories
func (f *IntegrationTestFramework) CountItemsInDirectory(dirPath string) int {
	allItems := f.GetAllItems()
	count := 0
	normalizedPath := directories.NormalizeDirectoryPath(dirPath)

	for _, item := range allItems {
		// Check if item's origin starts with the directory path
		if normalizedPath == model.ROOT_DIRECTORY_PATH {
			// For root directory, count all items
			count++
		} else if strings.HasPrefix(item.Origin, normalizedPath+"/") || item.Origin == normalizedPath {
			count++
		}
	}

	return count
}

func (f *IntegrationTestFramework) AutoIncludeHierarchy(dirPath string) error {
	return directories.AutoIncludeHierarchy(f.database, dirPath)
}

// AddCustomTagsToDirectory adds custom tags to a directory that will be applied as AutoTags to files
func (f *IntegrationTestFramework) AddCustomTagsToDirectory(dirPath string, customTags []*model.Tag) {
	normalizedPath := directories.NormalizeDirectoryPath(dirPath)

	// Get the directory from database
	dir, err := f.database.GetDirectory("path = ?", normalizedPath)
	require.NoError(f.t, err)
	require.NotNil(f.t, dir, "Directory %s should exist before adding custom tags", dirPath)

	// Update directory tags
	dir.Tags = customTags
	err = f.database.CreateOrUpdateDirectory(dir)
	require.NoError(f.t, err)
}

// AssertCustomAutoTagExists checks that a file has a custom AutoTag (child of a specific parent tag)
func (f *IntegrationTestFramework) AssertCustomAutoTagExists(dirPath, filename string, parentTag *model.Tag, expectedTitle string) {
	item, err := f.fsManager.GetBelongingItem(directories.NormalizeDirectoryPath(dirPath), filename)
	require.NoError(f.t, err)
	require.NotNil(f.t, item, "Item %s should exist in directory %s", filename, dirPath)

	hasCustomAutoTag := false
	for _, tag := range item.Tags {
		if tag.ParentID != nil && *tag.ParentID == parentTag.Id && tag.Title == expectedTitle {
			hasCustomAutoTag = true
			break
		}
	}

	assert.True(f.t, hasCustomAutoTag,
		"Item %s should have custom AutoTag '%s' (child of '%s') in directory %s",
		filename, expectedTitle, parentTag.Title, dirPath)
}

// AssertCustomAutoTagNotExists checks that a file does not have a custom AutoTag from a specific parent
func (f *IntegrationTestFramework) AssertCustomAutoTagNotExists(dirPath, filename string, parentTag *model.Tag) {
	item, err := f.fsManager.GetBelongingItem(directories.NormalizeDirectoryPath(dirPath), filename)
	require.NoError(f.t, err)
	require.NotNil(f.t, item, "Item %s should exist in directory %s", filename, dirPath)

	hasCustomAutoTag := false
	for _, tag := range item.Tags {
		if tag.ParentID != nil && *tag.ParentID == parentTag.Id {
			hasCustomAutoTag = true
			break
		}
	}

	assert.False(f.t, hasCustomAutoTag,
		"Item %s should not have any custom AutoTag from parent '%s' in directory %s",
		filename, parentTag.Title, dirPath)
}

// IncludeDirectory includes a directory (emulates user clicking include button)
func (f *IntegrationTestFramework) IncludeDirectory(dirPath string) {
	err := directories.IncludeDirectory(f.database, dirPath)
	require.NoError(f.t, err)
}

// ExcludeDirectory excludes a directory (emulates user clicking exclude button)
func (f *IntegrationTestFramework) ExcludeDirectory(dirPath string) {
	err := directories.ExcludeDirectory(f.database, dirPath)
	require.NoError(f.t, err)
}

// IncludeOrCreateDirectory includes a directory or creates it if missing
func (f *IntegrationTestFramework) IncludeOrCreateDirectory(dirPath string) {
	err := directories.IncludeOrCreateDirectory(f.database, dirPath)
	require.NoError(f.t, err)
}

// AssertDirectoryIncluded checks that a directory is included (not excluded)
func (f *IntegrationTestFramework) AssertDirectoryIncluded(dirPath string) {
	dir, err := f.database.GetDirectory("path = ?", directories.NormalizeDirectoryPath(dirPath))
	require.NoError(f.t, err)
	require.NotNil(f.t, dir, "Directory %s should exist", dirPath)
	assert.False(f.t, directories.IsExcluded(dir), "Directory %s should be included (not excluded)", dirPath)
}

// AssertDirectoryExcluded checks that a directory is excluded
func (f *IntegrationTestFramework) AssertDirectoryExcluded(dirPath string) {
	dir, err := f.database.GetDirectory("path = ?", directories.NormalizeDirectoryPath(dirPath))
	if err != nil {
		// If directory doesn't exist in database yet, it's effectively excluded
		assert.Contains(f.t, err.Error(), "record not found",
			"Directory %s should either be excluded or not exist yet (which means excluded)", dirPath)
		return
	}
	require.NotNil(f.t, dir, "Directory %s should exist", dirPath)
	assert.True(f.t, directories.IsExcluded(dir), "Directory %s should be excluded", dirPath)
}
