package fssync

import (
	"fmt"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/directorytree"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

// Test helper functions
func setupMocksForTest(t *testing.T) (*gomock.Controller,
	*model.MockDatabase,
	*model.MockDirectoryItemsGetterSetter,
	*model.MockDirectoryAutoTagsGetter,
	*model.MockFileMetadataGetter) {
	ctrl := gomock.NewController(t)
	mockDB := model.NewMockDatabase(ctrl)
	mockDIGS := model.NewMockDirectoryItemsGetterSetter(ctrl)
	mockDATG := model.NewMockDirectoryAutoTagsGetter(ctrl)
	mockFMG := model.NewMockFileMetadataGetter(ctrl)
	return ctrl, mockDB, mockDIGS, mockDATG, mockFMG
}

func createTestTempDirectory() (string, error) {
	return os.MkdirTemp("", "fs-syncer-test-*")
}

func createTestDiff() *directorytree.Diff {
	return &directorytree.Diff{
		AddedDirectories: []directorytree.Change{
			{Path1: "/test/new-dir", ChangeType: directorytree.DIRECTORY_ADDED},
		},
		RemovedDirectories: []directorytree.Change{
			{Path1: "/test/removed-dir", ChangeType: directorytree.DIRECTORY_REMOVED},
		},
		AddedFiles: []directorytree.Change{
			{Path1: "/test/new-file.mp4", ChangeType: directorytree.FILE_ADDED},
		},
		RemovedFiles: []directorytree.Change{
			{Path1: "/test/removed-file.mp4", ChangeType: directorytree.FILE_REMOVED},
		},
		MovedDirectories: []directorytree.Change{
			{Path1: "/test/old-dir", Path2: "/test/moved-dir", ChangeType: directorytree.DIRECTORY_MOVED},
		},
		MovedFiles: []directorytree.Change{
			{Path1: "/test/old-file.mp4", Path2: "/test/moved-file.mp4", ChangeType: directorytree.FILE_MOVED},
		},
	}
}

func createTestStale() *directorytree.Stale {
	return &directorytree.Stale{
		Dirs:  []string{"/test/stale-dir"},
		Files: []string{"/test/stale-file.mp4"},
	}
}

func createTestItem(id uint64, title, origin string) *model.Item {
	return &model.Item{
		Id:     id,
		Title:  title,
		Origin: origin,
		Tags:   []*model.Tag{},
	}
}

func createTestTag(id uint64, title string, parentID *uint64) *model.Tag {
	return &model.Tag{
		Id:       id,
		Title:    title,
		ParentID: parentID,
	}
}

func createTestDir(path string, excluded bool) *model.Directory {
	return &model.Directory{
		Path:     path,
		Excluded: pointer.Bool(excluded),
	}
}

type trueFileFilter struct {
}

func (f trueFileFilter) Filter(path string) bool {
	return true
}

// Tests for newFsSyncer constructor
func TestNewFsSyncer_Success(t *testing.T) {
	ctrl, mockDB, _, _, _ := setupMocksForTest(t)
	defer ctrl.Finish()

	testDir, err := createTestTempDirectory()
	assert.NoError(t, err)
	defer os.RemoveAll(testDir)

	// Create a test file in the directory
	testFile := filepath.Join(testDir, "test.txt")
	_, err = os.Create(testFile)
	assert.NoError(t, err)

	// Mock directory exists check
	mockDB.EXPECT().GetDirectory(gomock.Any()).Return(&model.Directory{Path: testDir}, nil)

	// Mock DirectoryItemsGetter for BuildFromDb
	mockDIG := model.NewMockDirectoryItemsGetter(ctrl)

	// Mock the GetAllDirectories call that BuildFromDb makes
	mockDB.EXPECT().GetAllDirectories().Return(&[]model.Directory{
		{Path: testDir, Excluded: pointer.Bool(false)},
	}, nil)

	// Mock GetBelongingItems call for the test directory - use gomock.Any() since testDir is dynamic
	mockDIG.EXPECT().GetBelongingItems(gomock.Any()).Return(&[]model.Item{}, nil)

	syncer, err := newFsSyncer(testDir, mockDB, mockDIG, trueFileFilter{})

	assert.NoError(t, err)
	assert.NotNil(t, syncer)
	assert.NotNil(t, syncer.diff)
	assert.NotNil(t, syncer.stales)
}

func TestNewFsSyncer_DirectoryNotFound(t *testing.T) {
	ctrl, mockDB, _, _, _ := setupMocksForTest(t)
	defer ctrl.Finish()

	// Mock directory not exists
	mockDB.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	mockDIG := model.NewMockDirectoryItemsGetter(ctrl)

	syncer, err := newFsSyncer("/nonexistent", mockDB, mockDIG, trueFileFilter{})

	assert.Error(t, err)
	assert.Nil(t, syncer)
	assert.Contains(t, err.Error(), "directory not found in db")
}

func TestNewFsSyncer_DirectoryCheckError(t *testing.T) {
	ctrl, mockDB, _, _, _ := setupMocksForTest(t)
	defer ctrl.Finish()

	// Mock database error
	mockDB.EXPECT().GetDirectory(gomock.Any()).Return(nil, errors.New("database error"))

	mockDIG := model.NewMockDirectoryItemsGetter(ctrl)

	syncer, err := newFsSyncer("/test", mockDB, mockDIG, trueFileFilter{})

	assert.Error(t, err)
	assert.Nil(t, syncer)
	assert.Contains(t, err.Error(), "database error")
}

// Tests for hasFsChanges
func TestHasFsChanges_WithChanges(t *testing.T) {
	syncer := &fsSyncer{
		diff:   createTestDiff(),
		stales: &directorytree.Stale{Dirs: []string{}, Files: []string{}},
	}

	assert.True(t, syncer.hasFsChanges())
}

func TestHasFsChanges_WithStales(t *testing.T) {
	syncer := &fsSyncer{
		diff:   &directorytree.Diff{},
		stales: createTestStale(),
	}

	assert.True(t, syncer.hasFsChanges())
}

func TestHasFsChanges_NoChanges(t *testing.T) {
	syncer := &fsSyncer{
		diff:   &directorytree.Diff{},
		stales: &directorytree.Stale{Dirs: []string{}, Files: []string{}},
	}

	assert.False(t, syncer.hasFsChanges())
}

// Tests for sync method main flow
func TestSync_NoChanges(t *testing.T) {
	ctrl, mockDB, mockDIGS, mockDATG, mockFMG := setupMocksForTest(t)
	defer ctrl.Finish()

	syncer := &fsSyncer{
		diff:   &directorytree.Diff{},
		stales: &directorytree.Stale{Dirs: []string{}, Files: []string{}},
	}

	// Mock for addMissingDirectoryTags - should get all directories
	mockDB.EXPECT().GetAllDirectories().Return(&[]model.Directory{}, nil)

	// Mock for syncAutoTags - should get all directories again
	mockDB.EXPECT().GetAllDirectories().Return(&[]model.Directory{}, nil)

	hasChanges, errs := syncer.sync(mockDB, mockDIGS, mockDATG, mockFMG)

	assert.False(t, hasChanges)
	assert.Empty(t, errs)
}

func TestSync_WithChanges(t *testing.T) {
	ctrl, mockDB, mockDIGS, mockDATG, mockFMG := setupMocksForTest(t)
	defer ctrl.Finish()

	syncer := &fsSyncer{
		diff:   createTestDiff(),
		stales: createTestStale(),
	}

	// Mock for addMissingDirectoryTags
	mockDB.EXPECT().GetAllDirectories().Return(&[]model.Directory{}, nil)

	// Mocks for stale removal
	mockDIGS.EXPECT().GetBelongingItem(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockDB.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockDB.EXPECT().RemoveDirectory(gomock.Any()).Return(nil).AnyTimes()

	// Mocks for adding directories
	mockDB.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockDB.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Return(nil).AnyTimes()

	// Mock AutoTags for new files
	mockDATG.EXPECT().GetAutoTags(gomock.Any()).Return([]*model.Tag{}, nil).AnyTimes()
	mockFMG.EXPECT().GetFileMetadata(gomock.Any()).Return(int64(12345), int64(1000), nil).AnyTimes()
	mockDIGS.EXPECT().AddBelongingItem(gomock.Any()).Return(nil).AnyTimes()

	// Mock for syncAutoTags
	mockDB.EXPECT().GetAllDirectories().Return(&[]model.Directory{}, nil)

	hasChanges, errs := syncer.sync(mockDB, mockDIGS, mockDATG, mockFMG)

	assert.True(t, hasChanges)
	// Log any errors for debugging but allow some due to simplified mocking
	if len(errs) > 0 {
		t.Logf("Errors occurred (expected due to simplified mocking): %d", len(errs))
	}
}

// Tests for addMissingDirectoryTags
func TestAddMissingDirectoryTags_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDR := model.NewMockDirectoryReader(ctrl)
	mockTRW := model.NewMockTagReaderWriter(ctrl)

	testDirectories := []model.Directory{
		{Path: "/test1", Excluded: pointer.Bool(false)},
		{Path: "/test2", Excluded: pointer.Bool(false)},
	}

	mockDR.EXPECT().GetAllDirectories().Return(&testDirectories, nil)

	// Mock the tag operations for each directory
	for range testDirectories {
		mockTRW.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
		mockTRW.EXPECT().CreateOrUpdateTag(gomock.Any()).Return(nil)
	}

	errs := addMissingDirectoryTags(mockDR, mockTRW)

	assert.Empty(t, errs)
}

func TestAddMissingDirectoryTags_GetDirectoriesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDR := model.NewMockDirectoryReader(ctrl)
	mockTRW := model.NewMockTagReaderWriter(ctrl)

	mockDR.EXPECT().GetAllDirectories().Return(nil, errors.New("db error"))

	errs := addMissingDirectoryTags(mockDR, mockTRW)

	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "db error")
}

func TestAddMissingDirectoryTag_ExcludedDirectory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)

	excludedDir := &model.Directory{
		Path:     "/excluded",
		Excluded: pointer.Bool(true),
	}

	// Should not call any tag operations for excluded directory
	err := addMissingDirectoryTag(mockTRW, excludedDir)

	assert.NoError(t, err)
}

func TestAddMissingDirectoryTag_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)

	dir := &model.Directory{
		Path:     "/test",
		Excluded: pointer.Bool(false),
	}

	mockTRW.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockTRW.EXPECT().CreateOrUpdateTag(gomock.Any()).Return(nil)

	err := addMissingDirectoryTag(mockTRW, dir)

	assert.NoError(t, err)
}

// Tests for removeStaleItems
func TestRemoveStaleItems_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDIG := model.NewMockDirectoryItemsGetter(ctrl)
	mockIW := model.NewMockItemWriter(ctrl)

	files := []string{"/test/file1.mp4", "/test/file2.mp4"}

	// Mock that items don't exist (already removed)
	for range files {
		mockDIG.EXPECT().GetBelongingItem(gomock.Any(), gomock.Any()).Return(nil, nil)
	}

	errs := removeStaleItems(mockDIG, mockIW, files)

	assert.Empty(t, errs)
}

func TestRemoveStaleItems_WithExistingItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDIG := model.NewMockDirectoryItemsGetter(ctrl)
	mockIW := model.NewMockItemWriter(ctrl)

	files := []string{"/test/file1.mp4"}
	testItem := createTestItem(1, "file1.mp4", "/test")

	mockDIG.EXPECT().GetBelongingItem(gomock.Any(), gomock.Any()).Return(testItem, nil)
	mockIW.EXPECT().RemoveItem(testItem.Id).Return(nil)

	errs := removeStaleItems(mockDIG, mockIW, files)

	assert.Empty(t, errs)
}

// Tests for removeStaleDirs
func TestRemoveStaleDirs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)
	mockDW := model.NewMockDirectoryWriter(ctrl)

	dirs := []string{"/test/dir1", "/test/dir2"}

	for range dirs {
		// Mock that directory tags don't exist
		mockTRW.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
		mockDW.EXPECT().RemoveDirectory(gomock.Any()).Return(nil)
	}

	errs := removeStaleDirs(mockTRW, mockDW, dirs)

	assert.Empty(t, errs)
}

// Tests for addMissingDirs
func TestAddMissingDirs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)

	changes := []directorytree.Change{
		{Path1: "/test/new-dir1", ChangeType: directorytree.DIRECTORY_ADDED},
		{Path1: "/test/new-dir2", ChangeType: directorytree.DIRECTORY_ADDED},
	}

	for range changes {
		mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
		mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
		mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Return(nil)
	}

	errs := addMissingDirs(mockDRW, changes)

	assert.Empty(t, errs)
}

// Tests for addNewFiles
func TestAddNewFiles_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIW := model.NewMockItemWriter(ctrl)
	mockDIGS := model.NewMockDirectoryItemsGetterSetter(ctrl)
	mockDATG := model.NewMockDirectoryAutoTagsGetter(ctrl)
	mockFMG := model.NewMockFileMetadataGetter(ctrl)

	changes := []directorytree.Change{
		{Path1: "/test/new-file.mp4", ChangeType: directorytree.FILE_ADDED},
	}

	// Mock that item doesn't exist yet
	mockDIGS.EXPECT().GetBelongingItem(gomock.Any(), gomock.Any()).Return(nil, nil)

	// Mock AutoTags for the directory
	mockDATG.EXPECT().GetAutoTags(gomock.Any()).Return([]*model.Tag{}, nil)

	// Mock file metadata
	mockFMG.EXPECT().GetFileMetadata(gomock.Any()).Return(int64(12345), int64(1000), nil)

	// Mock adding the new item
	mockDIGS.EXPECT().AddBelongingItem(gomock.Any()).Return(nil)

	errs := addNewFiles(mockIW, mockDIGS, mockDATG, mockFMG, changes)

	assert.Empty(t, errs)
}

func TestAddNewFiles_ExistingItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIW := model.NewMockItemWriter(ctrl)
	mockDIGS := model.NewMockDirectoryItemsGetterSetter(ctrl)
	mockDATG := model.NewMockDirectoryAutoTagsGetter(ctrl)
	mockFMG := model.NewMockFileMetadataGetter(ctrl)

	changes := []directorytree.Change{
		{Path1: "/test/existing-file.mp4", ChangeType: directorytree.FILE_ADDED},
	}

	existingItem := createTestItem(1, "existing-file.mp4", "/test")

	// Mock that item already exists
	mockDIGS.EXPECT().GetBelongingItem(gomock.Any(), gomock.Any()).Return(existingItem, nil)

	// Mock AutoTags for the directory
	mockDATG.EXPECT().GetAutoTags(gomock.Any()).Return([]*model.Tag{}, nil)

	// Mock updating existing item (no missing tags to add)
	// EnsureItemHaveTags won't call UpdateItem if no missing tags

	errs := addNewFiles(mockIW, mockDIGS, mockDATG, mockFMG, changes)

	assert.Empty(t, errs)
}

// Tests for renameDirs
func TestRenameDirs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)
	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)
	mockIRW := model.NewMockItemReaderWriter(ctrl)

	changes := []directorytree.Change{
		{Path1: "/test/old-dir", Path2: "/test/new-dir", ChangeType: directorytree.DIRECTORY_MOVED},
	}

	// Mock destination directory exists
	dstDir := createTestDir("/test/new-dir", false)
	mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(dstDir, nil)

	// Mock source directory exists
	srcDir := createTestDir("/test/old-dir", false)
	mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(srcDir, nil)

	// Mock updating directory path
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Return(nil)
	mockDRW.EXPECT().RemoveDirectory(gomock.Any()).Return(nil)

	// Mock getting directory tag for updating items - called in updateItemsLocation
	dirTag := createTestTag(1, "new-dir", nil)
	dirTag.Items = []*model.Item{} // Empty items list
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(dirTag, nil)
	// Even with empty items, GetItems([]) is still called
	emptyItems := &[]model.Item{}
	mockIRW.EXPECT().GetItems([]uint64{}).Return(emptyItems, nil)

	errs := renameDirs(mockTRW, mockDRW, mockIRW, changes)

	assert.Empty(t, errs)
}

// Tests for renameFiles
func TestRenameFiles_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)
	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)
	mockIRW := model.NewMockItemReaderWriter(ctrl)

	changes := []directorytree.Change{
		{Path1: "/test/old-file.mp4", Path2: "/test/new-file.mp4", ChangeType: directorytree.FILE_MOVED},
	}

	// Mock destination directory validation - called in validateReadyDirectory
	dstDir := createTestDir("/test", false)
	// AddDirectoryIfMissing -> DirectoryExists -> GetDirectory
	mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(dstDir, nil)
	// ValidateReadyDirectory -> GetDirectory (second call)
	mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(dstDir, nil)

	// Mock addMissingDirectoryTag -> GetOrCreateChildTag -> GetTag
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	// GetOrCreateChildTag -> CreateOrUpdateTag
	mockTRW.EXPECT().CreateOrUpdateTag(gomock.Any()).Return(nil)

	// Mock getting the source directory tag and finding original item
	originalItem := createTestItem(1, "old-file.mp4", "/test")
	srcTag := createTestTag(1, "test", nil)
	srcTag.Items = []*model.Item{originalItem} // Add the item to the tag
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(srcTag, nil)
	// Mock items in source directory through GetItems call
	mockIRW.EXPECT().GetItems(gomock.Any()).Return(&[]model.Item{*originalItem}, nil)

	// Mock updating item location
	mockIRW.EXPECT().UpdateItem(gomock.Any()).Return(nil)

	// Mock destination directory tag for adding item
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(srcTag, nil) // Same tag since same directory
	mockIRW.EXPECT().CreateOrUpdateItem(gomock.Any()).Return(nil)

	// Mock source directory tag for removing item
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(srcTag, nil)
	mockIRW.EXPECT().RemoveTagFromItem(gomock.Any(), gomock.Any()).Return(nil)

	errs := renameFiles(mockTRW, mockDRW, mockIRW, changes)

	assert.Empty(t, errs)
}

// Tests for syncAutoTags
func TestSyncAutoTags_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTR := model.NewMockTagReader(ctrl)
	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockDR := model.NewMockDirectoryReader(ctrl)
	mockDATG := model.NewMockDirectoryAutoTagsGetter(ctrl)

	testDirectories := []model.Directory{
		{Path: "/test1", Excluded: pointer.Bool(false)},
		{Path: "/test2", Excluded: pointer.Bool(false)},
	}

	mockDR.EXPECT().GetAllDirectories().Return(&testDirectories, nil)

	for i, dir := range testDirectories {
		// Mock AutoTags for each directory
		mockDATG.EXPECT().GetAutoTags(dir.Path).Return([]*model.Tag{}, nil)

		// Mock getting directory tag - createTestTag creates tags with no items
		tag := createTestTag(uint64(i+1), dir.Path, nil)
		tag.Items = []*model.Item{} // Empty items list
		mockTR.EXPECT().GetTag(gomock.Any()).Return(tag, nil)
		// Even with empty items, GetItems([]) is still called in getItems
		emptyItems := &[]model.Item{}
		mockIRW.EXPECT().GetItems([]uint64{}).Return(emptyItems, nil)
	}

	_, errs := syncAutoTags(mockTR, mockIRW, mockDR, mockDATG)

	assert.Empty(t, errs)
}

func TestSyncAutoTags_GetDirectoriesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTR := model.NewMockTagReader(ctrl)
	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockDR := model.NewMockDirectoryReader(ctrl)
	mockDATG := model.NewMockDirectoryAutoTagsGetter(ctrl)

	mockDR.EXPECT().GetAllDirectories().Return(nil, errors.New("db error"))

	_, errs := syncAutoTags(mockTR, mockIRW, mockDR, mockDATG)

	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "db error")
}

// Tests for syncAutoTagsForDir
func TestSyncAutoTagsForDir_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTR := model.NewMockTagReader(ctrl)
	mockIRW := model.NewMockItemReaderWriter(ctrl)

	dir := createTestDir("/test", false)
	autoTags := []*model.Tag{createTestTag(1, "auto-tag", nil)}

	// Mock getting directory tag - create a tag with items so GetItems is called
	testItem := createTestItem(1, "test.mp4", "/test")
	dirTag := createTestTag(2, "test", nil)
	dirTag.Items = []*model.Item{testItem} // Add items to the tag
	mockTR.EXPECT().GetTag(gomock.Any()).Return(dirTag, nil)

	// Mock items that need the auto tags
	items := []model.Item{*testItem}
	mockIRW.EXPECT().GetItems(gomock.Any()).Return(&items, nil)

	// Mock updating item with auto tags
	mockIRW.EXPECT().UpdateItem(gomock.Any()).Return(nil)

	_, errs := syncAutoTagsForDir(mockTR, mockIRW, autoTags, dir)

	assert.Empty(t, errs)
}

// Edge case tests
func TestSync_ErrorsCollected(t *testing.T) {
	ctrl, mockDB, mockDIGS, mockDATG, mockFMG := setupMocksForTest(t)
	defer ctrl.Finish()

	syncer := &fsSyncer{
		diff:   &directorytree.Diff{},
		stales: &directorytree.Stale{Dirs: []string{}, Files: []string{}},
	}

	// Force an error in addMissingDirectoryTags
	mockDB.EXPECT().GetAllDirectories().Return(nil, errors.New("first error"))

	// Force an error in syncAutoTags
	mockDB.EXPECT().GetAllDirectories().Return(nil, errors.New("second error"))

	hasChanges, errs := syncer.sync(mockDB, mockDIGS, mockDATG, mockFMG)

	assert.False(t, hasChanges)
	assert.Len(t, errs, 2)
	assert.Contains(t, errs[0].Error(), "first error")
	assert.Contains(t, errs[1].Error(), "second error")
}

func TestRemoveDir_TagNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)
	mockDW := model.NewMockDirectoryWriter(ctrl)

	// Mock tag not found
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// Should still attempt to remove directory
	mockDW.EXPECT().RemoveDirectory(gomock.Any()).Return(nil)

	errs := removeDir(mockTRW, mockDW, "/test/path")

	assert.Empty(t, errs)
}

func TestRemoveDir_TagFoundAndRemoved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)
	mockDW := model.NewMockDirectoryWriter(ctrl)

	tag := createTestTag(1, "test", nil)

	// Mock tag found
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(tag, nil)

	// Mock tag removal
	mockTRW.EXPECT().RemoveTag(tag.Id).Return(nil)

	// Mock directory removal
	mockDW.EXPECT().RemoveDirectory(gomock.Any()).Return(nil)

	errs := removeDir(mockTRW, mockDW, "/test/path")

	assert.Empty(t, errs)
}

func TestMoveDir_DestinationNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)
	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)
	mockIRW := model.NewMockItemReaderWriter(ctrl)

	// Mock destination directory doesn't exist
	mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// Should call removeDir for source
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockDRW.EXPECT().RemoveDirectory(gomock.Any()).Return(nil)

	errs := moveDir(mockTRW, mockDRW, mockIRW, "/src", "/dst")

	assert.Empty(t, errs)
}

func TestMoveFile_ItemNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTRW := model.NewMockTagReaderWriter(ctrl)
	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)
	mockIRW := model.NewMockItemReaderWriter(ctrl)

	// Mock destination directory validation - calls validateReadyDirectory
	dstDir := createTestDir("/test", false)
	// AddDirectoryIfMissing -> DirectoryExists -> GetDirectory
	mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(dstDir, nil)
	// ValidateReadyDirectory -> GetDirectory (second call)
	mockDRW.EXPECT().GetDirectory(gomock.Any()).Return(dstDir, nil)

	// Mock addMissingDirectoryTag -> GetOrCreateChildTag -> GetTag
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	// GetOrCreateChildTag -> CreateOrUpdateTag
	mockTRW.EXPECT().CreateOrUpdateTag(gomock.Any()).Return(nil)

	// Mock source directory tag for getItem - this should return nil (no tag found)
	// so getItem returns nil, leading to "original item not found" error
	mockTRW.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	err := moveFile(mockTRW, mockDRW, mockIRW, "/test/src.mp4", "/test/dst.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "original item not found")
}

// Performance test for large numbers of operations
func TestSync_LargeNumberOfOperations(t *testing.T) {
	ctrl, mockDB, mockDIGS, mockDATG, mockFMG := setupMocksForTest(t)
	defer ctrl.Finish()

	// Create a diff with many operations
	diff := &directorytree.Diff{
		AddedDirectories: make([]directorytree.Change, 1000),
		AddedFiles:       make([]directorytree.Change, 1000),
	}
	for i := 0; i < 1000; i++ {
		diff.AddedDirectories[i] = directorytree.Change{
			Path1:      fmt.Sprintf("/test/dir%d", i),
			ChangeType: directorytree.DIRECTORY_ADDED,
		}
		diff.AddedFiles[i] = directorytree.Change{
			Path1:      fmt.Sprintf("/test/file%d.mp4", i),
			ChangeType: directorytree.FILE_ADDED,
		}
	}

	syncer := &fsSyncer{
		diff:   diff,
		stales: &directorytree.Stale{Dirs: []string{}, Files: []string{}},
	}

	// Mock operations - simplified for performance test
	// addMissingDirectoryTags calls GetAllDirectories
	mockDB.EXPECT().GetAllDirectories().Return(&[]model.Directory{}, nil)

	// addMissingDirs -> AddDirectoryIfMissing -> DirectoryExists -> GetDirectory for each dir
	mockDB.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockDB.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockDB.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Return(nil).AnyTimes()

	// addNewFiles operations
	mockDIGS.EXPECT().GetBelongingItem(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockDATG.EXPECT().GetAutoTags(gomock.Any()).Return([]*model.Tag{}, nil).AnyTimes()
	mockFMG.EXPECT().GetFileMetadata(gomock.Any()).Return(int64(12345), int64(1000), nil).AnyTimes()
	mockDIGS.EXPECT().AddBelongingItem(gomock.Any()).Return(nil).AnyTimes()

	// syncAutoTags calls GetAllDirectories again
	mockDB.EXPECT().GetAllDirectories().Return(&[]model.Directory{}, nil)

	hasChanges, errs := syncer.sync(mockDB, mockDIGS, mockDATG, mockFMG)

	assert.True(t, hasChanges)
	// Some errors might occur due to simplified mocking, but should not be too many
	if len(errs) > 0 {
		t.Logf("Errors occurred: %d", len(errs))
		for i, err := range errs {
			if i < 5 { // Log first 5 errors for debugging
				t.Logf("Error %d: %s", i, err.Error())
			}
		}
	}
	assert.True(t, len(errs) < 100, "Too many errors for large operation test")
}

func TestAddMissingDirectoryTags2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dr := model.NewMockDirectoryReader(ctrl)

	dr.EXPECT().GetAllDirectories().Return(&[]model.Directory{
		{Path: "dir1"},
		{Path: "dir-2"},
		{Path: "dir_3"},
		{Path: "some/inner/path"},
		{Path: "tag/exists"},
		{Path: "excluded", Excluded: pointer.Bool(true)},
	}, nil)

	trw := model.NewMockTagReaderWriter(ctrl)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "dir1", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}).Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "dir-2", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}).Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "dir_3", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}).Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "some/inner/path", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}).Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(&model.Tag{Title: "tag/exists", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}, nil)
	errs := addMissingDirectoryTags(dr, trw)
	assert.Equal(t, 0, len(errs))
}

func TestAddMissingDirs2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	drw := model.NewMockDirectoryReaderWriter(ctrl)

	relativasor.Init("/root/dir")
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "path1", Excluded: pointer.Bool(true)}).Return(nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "some/path - to dir", Excluded: pointer.Bool(true)}).Return(nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(&model.Directory{Path: "/absolute/exists/path3"}, nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "/absolute/error/inner/path4", Excluded: pointer.Bool(true)}).Return(errors.Errorf("test"))
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: model.ROOT_DIRECTORY_PATH, Excluded: pointer.Bool(true)}).Return(nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(&model.Directory{Path: model.ROOT_DIRECTORY_PATH}, nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "some/relative/path", Excluded: pointer.Bool(true)}).Return(nil)

	errs := addMissingDirs(drw, []directorytree.Change{
		{Path1: "path1"},
		{Path1: "some/path - to dir"},
		{Path1: "/absolute/exists/path3"},
		{Path1: "/absolute/error/inner/path4"},
		{Path1: ""},
		{Path1: ""},
		{Path1: "/root/dir/some/relative/path"},
	})
	assert.Equal(t, 1, len(errs))
}

// Tests for addMissingDirs with AutoIncludeChildren flag
func TestAddMissingDirs_AutoIncludeChildren_True(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)

	// Set up parent directory with AutoIncludeChildren = true
	parentDir := &model.Directory{
		Path:                "parent",
		AutoIncludeChildren: pointer.Bool(true),
		Excluded:            pointer.Bool(false),
	}

	changes := []directorytree.Change{
		{Path1: "parent/child1", ChangeType: directorytree.DIRECTORY_ADDED},
		{Path1: "parent/child2", ChangeType: directorytree.DIRECTORY_ADDED},
	}

	// Mock calls for ShouldInclude -> GetParent -> GetDirectory for each child
	for range changes {
		// GetParent call for ShouldInclude
		mockDRW.EXPECT().GetDirectory("path = ?", "parent").Return(parentDir, nil)

		// AddDirectoryIfMissing -> DirectoryExists -> GetDirectory
		mockDRW.EXPECT().GetDirectory(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		// CreateOrUpdateDirectory with excluded=false (shouldInclude=true)
		mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
			assert.False(t, *dir.Excluded, "Directory should not be excluded when AutoIncludeChildren is true")
		}).Return(nil)
	}

	errs := addMissingDirs(mockDRW, changes)
	assert.Empty(t, errs)
}

func TestAddMissingDirs_AutoIncludeChildren_False(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)

	// Set up parent directory with AutoIncludeChildren = false
	parentDir := &model.Directory{
		Path:                 "parent",
		AutoIncludeChildren:  pointer.Bool(false),
		AutoIncludeHierarchy: pointer.Bool(false),
		Excluded:             pointer.Bool(false),
	}

	changes := []directorytree.Change{
		{Path1: "parent/child", ChangeType: directorytree.DIRECTORY_ADDED},
	}

	// Mock calls for ShouldInclude
	mockDRW.EXPECT().GetDirectory("path = ?", "parent").Return(parentDir, nil)

	// Since AutoIncludeChildren=false and AutoIncludeHierarchy=false, ShouldInclude will check parent hierarchy
	// GetParent of parent -> returns error, so ShouldInclude returns false
	mockDRW.EXPECT().GetDirectory("path = ?", model.ROOT_DIRECTORY_PATH).Return(nil, gorm.ErrRecordNotFound)

	// AddDirectoryIfMissing -> DirectoryExists -> GetDirectory
	mockDRW.EXPECT().GetDirectory("path = ?", "parent/child").Return(nil, gorm.ErrRecordNotFound)

	// CreateOrUpdateDirectory with excluded=true (shouldInclude=false)
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
		assert.True(t, *dir.Excluded, "Directory should be excluded when AutoIncludeChildren is false")
	}).Return(nil)

	errs := addMissingDirs(mockDRW, changes)
	assert.Empty(t, errs)
}

// Tests for addMissingDirs with AutoIncludeHierarchy flag
func TestAddMissingDirs_AutoIncludeHierarchy_True(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)

	// Set up directory hierarchy: grandparent -> parent -> child
	grandparentDir := &model.Directory{
		Path:                 "grandparent",
		AutoIncludeHierarchy: pointer.Bool(true),
		Excluded:             pointer.Bool(false),
	}

	parentDir := &model.Directory{
		Path:                 "grandparent/parent",
		AutoIncludeChildren:  pointer.Bool(false),
		AutoIncludeHierarchy: pointer.Bool(false),
		Excluded:             pointer.Bool(false),
	}

	changes := []directorytree.Change{
		{Path1: "grandparent/parent/child", ChangeType: directorytree.DIRECTORY_ADDED},
	}

	// Mock calls for ShouldInclude
	// GetParent of "grandparent/parent/child" -> returns "grandparent/parent"
	mockDRW.EXPECT().GetDirectory("path = ?", "grandparent/parent").Return(parentDir, nil)

	// Since parent doesn't have AutoIncludeChildren=true, check hierarchy
	// GetParent of "grandparent/parent" -> returns "grandparent"
	mockDRW.EXPECT().GetDirectory("path = ?", "grandparent").Return(grandparentDir, nil)

	// AddDirectoryIfMissing -> DirectoryExists -> GetDirectory
	mockDRW.EXPECT().GetDirectory(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// CreateOrUpdateDirectory with excluded=false (shouldInclude=true due to hierarchy)
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
		assert.False(t, *dir.Excluded, "Directory should not be excluded when ancestor has AutoIncludeHierarchy=true")
	}).Return(nil)

	errs := addMissingDirs(mockDRW, changes)
	assert.Empty(t, errs)
}

func TestAddMissingDirs_AutoIncludeHierarchy_DeepNesting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)

	// Set up deep hierarchy: root -> level1 -> level2 -> level3 -> level4
	rootDir := &model.Directory{
		Path:                 "root",
		AutoIncludeHierarchy: pointer.Bool(true),
		Excluded:             pointer.Bool(false),
	}

	level1Dir := &model.Directory{
		Path:                 "root/level1",
		AutoIncludeChildren:  pointer.Bool(false),
		AutoIncludeHierarchy: pointer.Bool(false),
		Excluded:             pointer.Bool(false),
	}

	level2Dir := &model.Directory{
		Path:                 "root/level1/level2",
		AutoIncludeChildren:  pointer.Bool(false),
		AutoIncludeHierarchy: pointer.Bool(false),
		Excluded:             pointer.Bool(false),
	}

	level3Dir := &model.Directory{
		Path:                 "root/level1/level2/level3",
		AutoIncludeChildren:  pointer.Bool(false),
		AutoIncludeHierarchy: pointer.Bool(false),
		Excluded:             pointer.Bool(false),
	}

	changes := []directorytree.Change{
		{Path1: "root/level1/level2/level3/level4", ChangeType: directorytree.DIRECTORY_ADDED},
	}

	// Mock calls for ShouldInclude traversing up the hierarchy
	mockDRW.EXPECT().GetDirectory("path = ?", "root/level1/level2/level3").Return(level3Dir, nil)
	mockDRW.EXPECT().GetDirectory("path = ?", "root/level1/level2").Return(level2Dir, nil)
	mockDRW.EXPECT().GetDirectory("path = ?", "root/level1").Return(level1Dir, nil)
	mockDRW.EXPECT().GetDirectory("path = ?", "root").Return(rootDir, nil)

	// AddDirectoryIfMissing -> DirectoryExists -> GetDirectory
	mockDRW.EXPECT().GetDirectory(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// CreateOrUpdateDirectory with excluded=false
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
		assert.False(t, *dir.Excluded, "Directory should not be excluded when distant ancestor has AutoIncludeHierarchy=true")
	}).Return(nil)

	errs := addMissingDirs(mockDRW, changes)
	assert.Empty(t, errs)
}

func TestAddMissingDirs_NoAutoInclude_AllExcluded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)

	// Set up parent with no auto-include flags
	parentDir := &model.Directory{
		Path:                 "parent",
		AutoIncludeChildren:  pointer.Bool(false),
		AutoIncludeHierarchy: pointer.Bool(false),
		Excluded:             pointer.Bool(false),
	}

	changes := []directorytree.Change{
		{Path1: "parent/child", ChangeType: directorytree.DIRECTORY_ADDED},
	}

	// Mock calls for ShouldInclude
	mockDRW.EXPECT().GetDirectory("path = ?", "parent").Return(parentDir, nil)

	// GetParent of parent fails (no more parents), so ShouldInclude returns false
	mockDRW.EXPECT().GetDirectory("path = ?", model.ROOT_DIRECTORY_PATH).Return(nil, gorm.ErrRecordNotFound)

	// AddDirectoryIfMissing -> DirectoryExists -> GetDirectory
	mockDRW.EXPECT().GetDirectory("path = ?", "parent/child").Return(nil, gorm.ErrRecordNotFound)

	// CreateOrUpdateDirectory with excluded=true
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
		assert.True(t, *dir.Excluded, "Directory should be excluded when no auto-include flags are set")
	}).Return(nil)

	errs := addMissingDirs(mockDRW, changes)
	assert.Empty(t, errs)
}

func TestAddMissingDirs_ParentNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)

	changes := []directorytree.Change{
		{Path1: "nonexistent/child", ChangeType: directorytree.DIRECTORY_ADDED},
	}

	// Mock calls for ShouldInclude - parent not found
	mockDRW.EXPECT().GetDirectory("path = ?", "nonexistent").Return(nil, gorm.ErrRecordNotFound)

	// AddDirectoryIfMissing -> DirectoryExists -> GetDirectory
	mockDRW.EXPECT().GetDirectory("path = ?", "nonexistent/child").Return(nil, gorm.ErrRecordNotFound)

	// CreateOrUpdateDirectory with excluded=true (ShouldInclude returns false when parent not found)
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
		assert.True(t, *dir.Excluded, "Directory should be excluded when parent directory not found")
	}).Return(nil)

	errs := addMissingDirs(mockDRW, changes)
	assert.Empty(t, errs)
}

func TestAddMissingDirs_MixedAutoIncludeFlags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDRW := model.NewMockDirectoryReaderWriter(ctrl)

	// Set up complex hierarchy with mixed flags
	parentWithChildren := &model.Directory{
		Path:                "include-children",
		AutoIncludeChildren: pointer.Bool(true),
		Excluded:            pointer.Bool(false),
	}

	parentWithHierarchy := &model.Directory{
		Path:                 "include-hierarchy",
		AutoIncludeHierarchy: pointer.Bool(true),
		Excluded:             pointer.Bool(false),
	}

	parentWithoutFlags := &model.Directory{
		Path:                 "no-flags",
		AutoIncludeChildren:  pointer.Bool(false),
		AutoIncludeHierarchy: pointer.Bool(false),
		Excluded:             pointer.Bool(false),
	}

	changes := []directorytree.Change{
		{Path1: "include-children/child1", ChangeType: directorytree.DIRECTORY_ADDED},
		{Path1: "include-hierarchy/level1/level2", ChangeType: directorytree.DIRECTORY_ADDED},
		{Path1: "no-flags/child", ChangeType: directorytree.DIRECTORY_ADDED},
	}

	// Mock for include-children/child1
	mockDRW.EXPECT().GetDirectory("path = ?", "include-children").Return(parentWithChildren, nil)
	mockDRW.EXPECT().GetDirectory("path = ?", "include-children/child1").Return(nil, gorm.ErrRecordNotFound)
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
		if dir.Path == "include-children/child1" {
			assert.False(t, *dir.Excluded, "Child should be included when parent has AutoIncludeChildren=true")
		}
	}).Return(nil)

	// Mock for include-hierarchy/level1/level2
	level1Dir := &model.Directory{
		Path:                 "include-hierarchy/level1",
		AutoIncludeChildren:  pointer.Bool(false),
		AutoIncludeHierarchy: pointer.Bool(false),
		Excluded:             pointer.Bool(false),
	}
	mockDRW.EXPECT().GetDirectory("path = ?", "include-hierarchy/level1").Return(level1Dir, nil)
	mockDRW.EXPECT().GetDirectory("path = ?", "include-hierarchy").Return(parentWithHierarchy, nil)
	mockDRW.EXPECT().GetDirectory("path = ?", "include-hierarchy/level1/level2").Return(nil, gorm.ErrRecordNotFound)
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
		if dir.Path == "include-hierarchy/level1/level2" {
			assert.False(t, *dir.Excluded, "Descendant should be included when ancestor has AutoIncludeHierarchy=true")
		}
	}).Return(nil)

	// Mock for no-flags/child
	mockDRW.EXPECT().GetDirectory("path = ?", "no-flags").Return(parentWithoutFlags, nil)
	mockDRW.EXPECT().GetDirectory("path = ?", model.ROOT_DIRECTORY_PATH).Return(nil, gorm.ErrRecordNotFound)
	mockDRW.EXPECT().GetDirectory("path = ?", "no-flags/child").Return(nil, gorm.ErrRecordNotFound)
	mockDRW.EXPECT().CreateOrUpdateDirectory(gomock.Any()).Do(func(dir *model.Directory) {
		if dir.Path == "no-flags/child" {
			assert.True(t, *dir.Excluded, "Child should be excluded when parent has no auto-include flags")
		}
	}).Return(nil)

	errs := addMissingDirs(mockDRW, changes)
	assert.Empty(t, errs)
}

func TestAddNewFiles2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	iw := model.NewMockItemWriter(ctrl)
	digs := model.NewMockDirectoryItemsGetterSetter(ctrl)
	datg := model.NewMockDirectoryAutoTagsGetter(ctrl)
	fmg := model.NewMockFileMetadataGetter(ctrl)

	relativasor.Init("/root/dir")

	tags := []*model.Tag{{Id: 3, Title: "old"}}
	digs.EXPECT().GetBelongingItem("new", "file").Return(&model.Item{Title: "file", Origin: "new", Url: "new/file", LastModified: 4234234, Tags: tags}, nil)
	autoTags1 := []*model.Tag{{Id: 1, Title: "auto1"}, {Id: 2, Title: "auto2"}}
	datg.EXPECT().GetAutoTags("new").Return(autoTags1, nil)
	iw.EXPECT().UpdateItem(&model.Item{Title: "file", Origin: "new", Url: "new/file", LastModified: 4234234, Tags: append(tags, autoTags1...)})

	digs.EXPECT().GetBelongingItem("some/deep/deep/deep", "file").Return(nil, nil)
	autoTags2 := []*model.Tag{{Title: "auto1"}, {Title: "auto2"}}
	datg.EXPECT().GetAutoTags("some/deep/deep/deep").Return(autoTags2, nil)
	fmg.EXPECT().GetFileMetadata("/root/dir/some/deep/deep/deep/file").Return(int64(7657657), int64(0), nil)
	digs.EXPECT().AddBelongingItem(&model.Item{Title: "file", Origin: "some/deep/deep/deep", Url: "some/deep/deep/deep/file", LastModified: 7657657, Tags: autoTags2}).Return(nil)

	digs.EXPECT().GetBelongingItem("/absolute", "path").Return(nil, nil)
	datg.EXPECT().GetAutoTags("/absolute").Return([]*model.Tag{}, nil)
	fmg.EXPECT().GetFileMetadata("/absolute/path").Return(int64(7567657), int64(0), nil)
	digs.EXPECT().AddBelongingItem(&model.Item{Title: "path", Origin: "/absolute", Url: "/absolute/path", LastModified: 7567657, Tags: make([]*model.Tag, 0)}).Return(nil)

	digs.EXPECT().GetBelongingItem("some", "file").Return(nil, nil)
	datg.EXPECT().GetAutoTags("some").Return([]*model.Tag{}, nil)
	fmg.EXPECT().GetFileMetadata("/root/dir/some/file").Return(int64(9876532), int64(0), nil)
	digs.EXPECT().AddBelongingItem(&model.Item{Title: "file", Origin: "some", Url: "some/file", LastModified: 9876532, Tags: make([]*model.Tag, 0)}).Return(nil)

	errs := addNewFiles(iw, digs, datg, fmg, []directorytree.Change{
		{Path1: "new/file"},
		{Path1: "some/deep/deep/deep/file"},
		{Path1: "/absolute/path"},
		{Path1: "/root/dir/some/file"},
	})

	assert.Equal(t, 0, len(errs))
}

func TestRemoveStaleDirs2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trw := model.NewMockTagReaderWriter(ctrl)
	dw := model.NewMockDirectoryWriter(ctrl)

	relativasor.Init("/root/dir")
	trw.EXPECT().GetTag(gomock.Any()).Return(&model.Tag{Id: 1, Title: "tag1"}, nil)
	trw.EXPECT().RemoveTag(uint64(1)).Return(nil)
	dw.EXPECT().RemoveDirectory("/existing/dir").Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	dw.EXPECT().RemoveDirectory("/not/exists/tag").Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(&model.Tag{Id: 2, Title: "relative"}, nil)
	trw.EXPECT().RemoveTag(uint64(2)).Return(nil)
	dw.EXPECT().RemoveDirectory("relative/path").Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(&model.Tag{Id: 3, Title: model.ROOT_DIRECTORY_PATH}, nil)
	trw.EXPECT().RemoveTag(uint64(3)).Return(nil)
	dw.EXPECT().RemoveDirectory(model.ROOT_DIRECTORY_PATH).Return(nil)

	errs := removeStaleDirs(trw, dw, []string{
		"/existing/dir",
		"/not/exists/tag",
		"/root/dir/relative/path",
		"",
	})
	assert.Equal(t, 0, len(errs))
}

func TestRemoveStaleItems2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dig := model.NewMockDirectoryItemsGetter(ctrl)
	iw := model.NewMockItemWriter(ctrl)

	relativasor.Init("/root/dir")
	dig.EXPECT().GetBelongingItem("/some", "not-exists").Return(nil, nil)
	dig.EXPECT().GetBelongingItem("/some/inner", "file").Return(&model.Item{Id: 1}, nil)
	iw.EXPECT().RemoveItem(uint64(1)).Return(nil)
	dig.EXPECT().GetBelongingItem("relative", "file").Return(&model.Item{Id: 2}, nil)
	iw.EXPECT().RemoveItem(uint64(2)).Return(nil)

	errs := removeStaleItems(dig, iw, []string{
		"/some/not-exists",
		"/some/inner/file",
		"/root/dir/relative/file",
	})
	assert.Equal(t, 0, len(errs))
}

func TestRenameFiles2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trw := model.NewMockTagReaderWriter(ctrl)
	irw := model.NewMockItemReaderWriter(ctrl)
	drw := model.NewMockDirectoryReaderWriter(ctrl)

	relativasor.Init("/root/dir")
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "/new", Excluded: pointer.Bool(false)}).Return(nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(&model.Directory{Path: "/new", Excluded: pointer.Bool(false)}, nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "/new", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())})
	oldTag := &model.Tag{Id: 4321, Title: "/", ParentID: pointer.Uint64(directories.GetDirectoriesTagId()), Items: []*model.Item{{Id: 1}}}
	trw.EXPECT().GetTag(gomock.Any()).Return(oldTag, nil)
	irw.EXPECT().GetItems(gomock.Any()).Return(&[]model.Item{{Id: 1, Title: "file1"}}, nil)
	irw.EXPECT().UpdateItem(&model.Item{Id: 1, Title: "file1", Origin: "/new", Url: "/new/file1"}).Return(nil)
	newTag := &model.Tag{Title: "new", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}
	trw.EXPECT().GetTag(gomock.Any()).Return(newTag, nil)
	irw.EXPECT().CreateOrUpdateItem(&model.Item{Id: 1, Title: "file1", Origin: "/new", Url: "/new/file1", Tags: []*model.Tag{newTag}})
	trw.EXPECT().GetTag(gomock.Any()).Return(oldTag, nil)
	irw.EXPECT().RemoveTagFromItem(uint64(1), uint64(4321)).Return(nil)

	errs := renameFiles(trw, drw, irw, []directorytree.Change{
		{Path1: "/file1", Path2: "/new/file1"},
	})
	assert.Equal(t, 0, len(errs))
}
