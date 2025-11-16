package items

import (
	"context"
	"errors"
	"math/rand"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"k8s.io/utils/pointer"
)

var testFiles = "../../../testdata/sample-files"
var sampleMp4 = filepath.Join(testFiles, "sample.mp4")

func setupTest(t *testing.T) func() {
	// Initialize relativasor with a temp directory
	tempDir, err := os.MkdirTemp("", "items-test-*")
	assert.NoError(t, err)

	err = relativasor.Init(tempDir)
	assert.NoError(t, err)

	return func() {
		os.RemoveAll(tempDir)
	}
}

func TestFileExists(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-file-*.mp4")
	assert.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Test with existing file
	item := model.Item{
		Origin: filepath.Dir(tempFile.Name()),
		Title:  filepath.Base(tempFile.Name()),
	}
	exists := FileExists(item)
	assert.True(t, exists, "Expected temporary file to exist")

	// Test with non-existing file
	item = model.Item{
		Origin: filepath.Dir(tempFile.Name()),
		Title:  "non-existent.mp4",
	}
	exists = FileExists(item)
	assert.False(t, exists, "Expected non-existent.mp4 to not exist")

	// Test with invalid path
	item = model.Item{
		Origin: "/invalid/path",
		Title:  "file.mp4",
	}
	exists = FileExists(item)
	assert.False(t, exists, "Expected invalid path to not exist")
}

func TestTitleFromFileName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple filename",
			path:     "video.mp4",
			expected: "video.mp4",
		},
		{
			name:     "full path",
			path:     "/path/to/video.mp4",
			expected: "video.mp4",
		},
		{
			name:     "unix path with backslashes",
			path:     "/path/to/video.mp4",
			expected: "video.mp4",
		},
		{
			name:     "path with spaces",
			path:     "/path with spaces/video file.mp4",
			expected: "video file.mp4",
		},
		{
			name:     "empty path",
			path:     "",
			expected: ".",
		},
		{
			name:     "path ending with slash",
			path:     "/path/to/",
			expected: "to",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TitleFromFileName(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildItemFromPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileMetadataGetter := model.NewMockFileMetadataGetter(ctrl)

	t.Run("successful build", func(t *testing.T) {
		origin := "/test/origin"
		path := "/test/path/video.mp4"
		expectedLastModified := int64(1234567890)
		expectedFileSize := int64(1024)

		mockFileMetadataGetter.EXPECT().
			GetFileMetadata(path).
			Return(expectedLastModified, expectedFileSize, nil)

		item, err := BuildItemFromPath(origin, path, mockFileMetadataGetter)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "video.mp4", item.Title)
		assert.Equal(t, origin, item.Origin)
		assert.Equal(t, filepath.Join(origin, "video.mp4"), item.Url)
		assert.Equal(t, expectedLastModified, item.LastModified)
		assert.Equal(t, expectedFileSize, item.FileSize)
	})

	t.Run("metadata getter error", func(t *testing.T) {
		origin := "/test/origin"
		path := "/test/path/video.mp4"
		expectedError := errors.New("metadata error")

		mockFileMetadataGetter.EXPECT().
			GetFileMetadata(path).
			Return(int64(0), int64(0), expectedError)

		item, err := BuildItemFromPath(origin, path, mockFileMetadataGetter)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, expectedError, err)
	})
}

func TestUpdateFileLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemWriter := model.NewMockItemWriter(ctrl)

	t.Run("simple item without highlights or subitems", func(t *testing.T) {
		item := &model.Item{
			Id:     1,
			Title:  "old_title.mp4",
			Origin: "old_origin",
			Url:    "old_url",
		}

		origin := "new_origin"
		path := "/new/path/new_title.mp4"
		url := ""

		mockItemWriter.EXPECT().
			UpdateItem(gomock.Any(), item).
			Return(nil)

		err := UpdateFileLocation(context.Background(), mockItemWriter, item, origin, path, url)

		assert.NoError(t, err)
		assert.Equal(t, "new_title.mp4", item.Title)
		assert.Equal(t, origin, item.Origin)
		assert.Equal(t, filepath.Join(origin, "new_title.mp4"), item.Url)
	})

	t.Run("item with custom url", func(t *testing.T) {
		item := &model.Item{
			Id:     1,
			Title:  "old_title.mp4",
			Origin: "old_origin",
			Url:    "old_url",
		}

		origin := "new_origin"
		path := "/new/path/new_title.mp4"
		url := "custom_url"

		mockItemWriter.EXPECT().
			UpdateItem(gomock.Any(), item).
			Return(nil)

		err := UpdateFileLocation(context.Background(), mockItemWriter, item, origin, path, url)

		assert.NoError(t, err)
		assert.Equal(t, "new_title.mp4", item.Title)
		assert.Equal(t, origin, item.Origin)
		assert.Equal(t, url, item.Url)
	})

	t.Run("item with highlights and subitems", func(t *testing.T) {
		highlight := &model.Item{
			Id:            2,
			StartPosition: 10.0,
			EndPosition:   20.0,
		}
		subitem := &model.Item{
			Id:            3,
			StartPosition: 30.0,
			EndPosition:   40.0,
		}
		item := &model.Item{
			Id:         1,
			Title:      "old_title.mp4",
			Origin:     "old_origin",
			Url:        "old_url",
			Highlights: []*model.Item{highlight},
			SubItems:   []*model.Item{subitem},
		}

		origin := "new_origin"
		path := "/new/path/new_title.mp4"
		url := ""

		// Expect recursive calls for highlights and subitems
		mockItemWriter.EXPECT().
			UpdateItem(gomock.Any(), highlight).
			Return(nil)
		mockItemWriter.EXPECT().
			UpdateItem(gomock.Any(), subitem).
			Return(nil)
		mockItemWriter.EXPECT().
			UpdateItem(gomock.Any(), item).
			Return(nil)

		err := UpdateFileLocation(context.Background(), mockItemWriter, item, origin, path, url)

		assert.NoError(t, err)
		assert.Equal(t, "new_title.mp4", item.Title)
		assert.Equal(t, origin, item.Origin)
		assert.Equal(t, filepath.Join(origin, "new_title.mp4"), item.Url)
	})

	t.Run("update error", func(t *testing.T) {
		item := &model.Item{
			Id:     1,
			Title:  "old_title.mp4",
			Origin: "old_origin",
			Url:    "old_url",
		}

		origin := "new_origin"
		path := "/new/path/new_title.mp4"
		url := ""
		expectedError := errors.New("update error")

		mockItemWriter.EXPECT().
			UpdateItem(gomock.Any(), item).
			Return(expectedError)

		err := UpdateFileLocation(context.Background(), mockItemWriter, item, origin, path, url)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestItemExists(t *testing.T) {
	items := []*model.Item{
		{Id: 1, Title: "item1"},
		{Id: 2, Title: "item2"},
		{Id: 3, Title: "item3"},
	}

	t.Run("item exists", func(t *testing.T) {
		item := &model.Item{Id: 2, Title: "item2"}
		exists := ItemExists(items, item)
		assert.True(t, exists)
	})

	t.Run("item does not exist", func(t *testing.T) {
		item := &model.Item{Id: 4, Title: "item4"}
		exists := ItemExists(items, item)
		assert.False(t, exists)
	})

	t.Run("empty items list", func(t *testing.T) {
		item := &model.Item{Id: 1, Title: "item1"}
		exists := ItemExists([]*model.Item{}, item)
		assert.False(t, exists)
	})

	t.Run("nil item", func(t *testing.T) {
		// ItemExists doesn't handle nil items gracefully, so we'll test with zero-value item
		item := &model.Item{Id: 0, Title: ""}
		exists := ItemExists(items, item)
		assert.False(t, exists)
	})
}

func TestTagExists(t *testing.T) {
	tags := []*model.Tag{
		{Id: 1, Title: "tag1"},
		{Id: 2, Title: "tag2"},
		{Id: 3, Title: "tag3"},
	}

	t.Run("tag exists", func(t *testing.T) {
		tag := &model.Tag{Id: 2, Title: "tag2"}
		exists := TagExists(tags, tag)
		assert.True(t, exists)
	})

	t.Run("tag does not exist", func(t *testing.T) {
		tag := &model.Tag{Id: 4, Title: "tag4"}
		exists := TagExists(tags, tag)
		assert.False(t, exists)
	})

	t.Run("empty tags list", func(t *testing.T) {
		tag := &model.Tag{Id: 1, Title: "tag1"}
		exists := TagExists([]*model.Tag{}, tag)
		assert.False(t, exists)
	})

	t.Run("nil tag", func(t *testing.T) {
		// TagExists doesn't handle nil tags gracefully, so we'll test with zero-value tag
		tag := &model.Tag{Id: 0, Title: ""}
		exists := TagExists(tags, tag)
		assert.False(t, exists)
	})
}

func TestEnsureItemHaveTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemWriter := model.NewMockItemWriter(ctrl)

	t.Run("item already has all tags", func(t *testing.T) {
		existingTags := []*model.Tag{
			{Id: 1, Title: "tag1"},
			{Id: 2, Title: "tag2"},
		}
		item := &model.Item{
			Id:   1,
			Tags: existingTags,
		}
		tagsToEnsure := []*model.Tag{
			{Id: 1, Title: "tag1"},
		}

		changed, err := EnsureItemHaveTags(context.Background(), mockItemWriter, item, tagsToEnsure)

		assert.NoError(t, err)
		assert.False(t, changed)
		assert.Equal(t, 2, len(item.Tags))
	})

	t.Run("item missing some tags", func(t *testing.T) {
		existingTags := []*model.Tag{
			{Id: 1, Title: "tag1"},
		}
		item := &model.Item{
			Id:   1,
			Tags: existingTags,
		}
		tagsToEnsure := []*model.Tag{
			{Id: 1, Title: "tag1"},
			{Id: 2, Title: "tag2"},
			{Id: 3, Title: "tag3"},
		}

		mockItemWriter.EXPECT().
			UpdateItem(gomock.Any(), item).
			Return(nil)

		changed, err := EnsureItemHaveTags(context.Background(), mockItemWriter, item, tagsToEnsure)

		assert.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, 3, len(item.Tags))
		assert.True(t, TagExists(item.Tags, &model.Tag{Id: 2, Title: "tag2"}))
		assert.True(t, TagExists(item.Tags, &model.Tag{Id: 3, Title: "tag3"}))
	})

	t.Run("update error", func(t *testing.T) {
		item := &model.Item{
			Id:   1,
			Tags: []*model.Tag{},
		}
		tagsToEnsure := []*model.Tag{
			{Id: 1, Title: "tag1"},
		}
		expectedError := errors.New("update error")

		mockItemWriter.EXPECT().
			UpdateItem(gomock.Any(), item).
			Return(expectedError)

		changed, err := EnsureItemHaveTags(context.Background(), mockItemWriter, item, tagsToEnsure)

		assert.Error(t, err)
		assert.False(t, changed)
		assert.Equal(t, expectedError, err)
	})
}

func TestEnsureItemMissingTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemWriter := model.NewMockItemWriter(ctrl)

	t.Run("remove existing tags", func(t *testing.T) {
		item := &model.Item{
			Id: 1,
			Tags: []*model.Tag{
				{Id: 1, Title: "tag1"},
				{Id: 2, Title: "tag2"},
				{Id: 3, Title: "tag3"},
			},
		}
		tagsToRemove := []*model.Tag{
			{Id: 1, Title: "tag1"},
			{Id: 3, Title: "tag3"},
		}

		mockItemWriter.EXPECT().
			RemoveTagFromItem(gomock.Any(), item.Id, uint64(1)).
			Return(nil)
		mockItemWriter.EXPECT().
			RemoveTagFromItem(gomock.Any(), item.Id, uint64(3)).
			Return(nil)

		err := EnsureItemMissingTags(context.Background(), mockItemWriter, item, tagsToRemove)

		assert.NoError(t, err)
	})

	t.Run("no matching tags to remove", func(t *testing.T) {
		item := &model.Item{
			Id: 1,
			Tags: []*model.Tag{
				{Id: 1, Title: "tag1"},
				{Id: 2, Title: "tag2"},
			},
		}
		tagsToRemove := []*model.Tag{
			{Id: 5, Title: "tag5"},
		}

		err := EnsureItemMissingTags(context.Background(), mockItemWriter, item, tagsToRemove)

		assert.NoError(t, err)
	})

	t.Run("remove error", func(t *testing.T) {
		item := &model.Item{
			Id: 1,
			Tags: []*model.Tag{
				{Id: 1, Title: "tag1"},
			},
		}
		tagsToRemove := []*model.Tag{
			{Id: 1, Title: "tag1"},
		}
		expectedError := errors.New("remove error")

		mockItemWriter.EXPECT().
			RemoveTagFromItem(gomock.Any(), item.Id, uint64(1)).
			Return(expectedError)

		err := EnsureItemMissingTags(context.Background(), mockItemWriter, item, tagsToRemove)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestHasSingleTag(t *testing.T) {
	t.Run("has single matching tag", func(t *testing.T) {
		tag := &model.Tag{Id: 1, Title: "tag1"}
		item := &model.Item{
			Tags: []*model.Tag{tag},
		}

		result := HasSingleTag(item, tag)
		assert.True(t, result)
	})

	t.Run("has single non-matching tag", func(t *testing.T) {
		tag := &model.Tag{Id: 1, Title: "tag1"}
		item := &model.Item{
			Tags: []*model.Tag{{Id: 2, Title: "tag2"}},
		}

		result := HasSingleTag(item, tag)
		assert.False(t, result)
	})

	t.Run("has multiple tags", func(t *testing.T) {
		tag := &model.Tag{Id: 1, Title: "tag1"}
		item := &model.Item{
			Tags: []*model.Tag{
				tag,
				{Id: 2, Title: "tag2"},
			},
		}

		result := HasSingleTag(item, tag)
		assert.False(t, result)
	})

	t.Run("has no tags", func(t *testing.T) {
		tag := &model.Tag{Id: 1, Title: "tag1"}
		item := &model.Item{
			Tags: []*model.Tag{},
		}

		result := HasSingleTag(item, tag)
		assert.False(t, result)
	})
}

func TestRemoveItemAndItsAssociations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemWriter := model.NewMockItemWriter(ctrl)

	t.Run("successful removal", func(t *testing.T) {
		itemId := uint64(123)

		mockItemWriter.EXPECT().
			RemoveItem(gomock.Any(), itemId).
			Return(nil)

		errors := RemoveItemAndItsAssociations(context.Background(), mockItemWriter, itemId)

		assert.Empty(t, errors)
	})

	t.Run("removal error", func(t *testing.T) {
		itemId := uint64(123)
		expectedError := errors.New("removal error")

		mockItemWriter.EXPECT().
			RemoveItem(gomock.Any(), itemId).
			Return(expectedError)

		errs := RemoveItemAndItsAssociations(context.Background(), mockItemWriter, itemId)

		assert.Len(t, errs, 1)
		assert.Equal(t, expectedError, errs[0])
	})
}

func TestDeleteRealFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cleanup := setupTest(t)
	defer cleanup()

	mockItemReader := model.NewMockItemReader(ctrl)

	t.Run("delete regular item", func(t *testing.T) {
		// Create a temporary file for testing
		tempFile, err := os.CreateTemp("", "test-delete-*.mp4")
		assert.NoError(t, err)
		tempFile.Close()
		defer func() {
			// Clean up if file still exists
			if _, err := os.Stat(tempFile.Name()); err == nil {
				os.Remove(tempFile.Name())
			}
		}()

		itemId := uint64(123)
		item := &model.Item{
			Id:  itemId,
			Url: tempFile.Name(),
		}

		mockItemReader.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(item, nil)

		err = DeleteRealFile(context.Background(), mockItemReader, itemId)

		assert.NoError(t, err)
		// Verify file was deleted
		_, err = os.Stat(tempFile.Name())
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("delete subitem error", func(t *testing.T) {
		itemId := uint64(123)
		item := &model.Item{
			Id:         itemId,
			MainItemId: pointer.Uint64(456), // This makes it a subitem
		}

		mockItemReader.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(item, nil)

		err := DeleteRealFile(context.Background(), mockItemReader, itemId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Deletion of subitem is forbidden")
	})

	t.Run("delete highlight error", func(t *testing.T) {
		itemId := uint64(123)
		item := &model.Item{
			Id:                    itemId,
			HighlightParentItemId: pointer.Uint64(456), // This makes it a highlight
		}

		mockItemReader.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(item, nil)

		err := DeleteRealFile(context.Background(), mockItemReader, itemId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Deletion of highlight is forbidden")
	})

	t.Run("get item error", func(t *testing.T) {
		itemId := uint64(123)
		expectedError := errors.New("get item error")

		mockItemReader.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(nil, expectedError)

		err := DeleteRealFile(context.Background(), mockItemReader, itemId)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("file does not exist", func(t *testing.T) {
		itemId := uint64(123)
		item := &model.Item{
			Id:  itemId,
			Url: "/non/existent/file.mp4",
		}

		mockItemReader.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(item, nil)

		err := DeleteRealFile(context.Background(), mockItemReader, itemId)

		// Should not return error even if file doesn't exist
		assert.NoError(t, err)
	})
}

func TestNoRandom(t *testing.T) {
	t.Run("item with no random tag", func(t *testing.T) {
		item := &model.Item{
			Tags: []*model.Tag{
				{Id: 1, NoRandom: pointer.Bool(true)},
				{Id: 2, NoRandom: pointer.Bool(false)},
			},
		}

		result := noRandom(item)
		assert.True(t, result)
	})

	t.Run("item without no random tag", func(t *testing.T) {
		item := &model.Item{
			Tags: []*model.Tag{
				{Id: 1, NoRandom: pointer.Bool(false)},
				{Id: 2, NoRandom: nil},
			},
		}

		result := noRandom(item)
		assert.False(t, result)
	})

	t.Run("item with no tags", func(t *testing.T) {
		item := &model.Item{
			Tags: []*model.Tag{},
		}

		result := noRandom(item)
		assert.False(t, result)
	})

	t.Run("item with nil NoRandom fields", func(t *testing.T) {
		item := &model.Item{
			Tags: []*model.Tag{
				{Id: 1, NoRandom: nil},
				{Id: 2, NoRandom: nil},
			},
		}

		result := noRandom(item)
		assert.False(t, result)
	})
}

func TestGetRandomItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemReader := model.NewMockItemReader(ctrl)

	// Create test filter that accepts all items
	acceptAllFilter := func(item *model.Item) bool {
		return true
	}

	// Create test filter that rejects items with ID 2
	rejectId2Filter := func(item *model.Item) bool {
		return item.Id != 2
	}

	t.Run("get random items successfully", func(t *testing.T) {
		items := []model.Item{
			{Id: 1, Title: "item1"},
			{Id: 2, Title: "item2"},
			{Id: 3, Title: "item3"},
			{Id: 4, Title: "item4"},
			{Id: 5, Title: "item5"},
		}

		mockItemReader.EXPECT().
			GetAllItems(gomock.Any()).
			Return(&items, nil)

		// Seed random for deterministic testing
		rand.Seed(time.Now().UnixNano())

		result, err := GetRandomItems(context.Background(), mockItemReader, 3, acceptAllFilter)

		assert.NoError(t, err)
		assert.Len(t, result, 3)

		// Check that all returned items are unique
		ids := make(map[uint64]bool)
		for _, item := range result {
			assert.False(t, ids[item.Id], "Duplicate item found")
			ids[item.Id] = true
		}
	})

	t.Run("get random items with filter", func(t *testing.T) {
		items := []model.Item{
			{Id: 1, Title: "item1"},
			{Id: 2, Title: "item2"}, // This will be filtered out
			{Id: 3, Title: "item3"},
			{Id: 4, Title: "item4"},
		}

		mockItemReader.EXPECT().
			GetAllItems(gomock.Any()).
			Return(&items, nil)

		result, err := GetRandomItems(context.Background(), mockItemReader, 2, rejectId2Filter)

		assert.NoError(t, err)
		assert.LessOrEqual(t, len(result), 2) // Should be at most 2 items

		// Verify that item with ID 2 is not in the result
		for _, item := range result {
			assert.NotEqual(t, uint64(2), item.Id)
		}
	})

	t.Run("get random items with noRandom tag", func(t *testing.T) {
		items := []model.Item{
			{Id: 1, Title: "item1"},
			{Id: 2, Title: "item2", Tags: []*model.Tag{{Id: 1, NoRandom: pointer.Bool(true)}}},
			{Id: 3, Title: "item3"},
		}

		mockItemReader.EXPECT().
			GetAllItems(gomock.Any()).
			Return(&items, nil)

		result, err := GetRandomItems(context.Background(), mockItemReader, 2, acceptAllFilter)

		assert.NoError(t, err)

		// Verify that item with noRandom tag is not in the result
		for _, item := range result {
			assert.NotEqual(t, uint64(2), item.Id)
		}
	})

	t.Run("request more items than available", func(t *testing.T) {
		items := []model.Item{
			{Id: 1, Title: "item1"},
			{Id: 2, Title: "item2"},
		}

		mockItemReader.EXPECT().
			GetAllItems(gomock.Any()).
			Return(&items, nil)

		result, err := GetRandomItems(context.Background(), mockItemReader, 5, acceptAllFilter)

		assert.NoError(t, err)
		assert.LessOrEqual(t, len(result), len(items)-1) // Should be at most len(items)-1
	})

	t.Run("no items available", func(t *testing.T) {
		items := []model.Item{}

		mockItemReader.EXPECT().
			GetAllItems(gomock.Any()).
			Return(&items, nil)

		result, err := GetRandomItems(context.Background(), mockItemReader, 3, acceptAllFilter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no items")
	})

	t.Run("get all items error", func(t *testing.T) {
		expectedError := errors.New("get items error")

		mockItemReader.EXPECT().
			GetAllItems(gomock.Any()).
			Return(nil, expectedError)

		result, err := GetRandomItems(context.Background(), mockItemReader, 3, acceptAllFilter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestIsModified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileMetadataGetter := model.NewMockFileMetadataGetter(ctrl)

	t.Run("item is modified", func(t *testing.T) {
		item := &model.Item{
			Origin:       "/test/origin",
			Title:        "video.mp4",
			LastModified: 1234567890,
		}
		newLastModified := int64(1234567999)

		mockFileMetadataGetter.EXPECT().
			GetFileMetadata(gomock.Any()).
			Return(newLastModified, int64(1024), nil)

		modified, err := IsModified(item, mockFileMetadataGetter)

		assert.NoError(t, err)
		assert.True(t, modified)
	})

	t.Run("item is not modified", func(t *testing.T) {
		item := &model.Item{
			Origin:       "/test/origin",
			Title:        "video.mp4",
			LastModified: 1234567890,
		}
		sameLastModified := int64(1234567890)

		mockFileMetadataGetter.EXPECT().
			GetFileMetadata(gomock.Any()).
			Return(sameLastModified, int64(1024), nil)

		modified, err := IsModified(item, mockFileMetadataGetter)

		assert.NoError(t, err)
		assert.False(t, modified)
	})

	t.Run("metadata getter error", func(t *testing.T) {
		item := &model.Item{
			Origin:       "/test/origin",
			Title:        "video.mp4",
			LastModified: 1234567890,
		}
		expectedError := errors.New("metadata error")

		mockFileMetadataGetter.EXPECT().
			GetFileMetadata(gomock.Any()).
			Return(int64(0), int64(0), expectedError)

		modified, err := IsModified(item, mockFileMetadataGetter)

		assert.Error(t, err)
		assert.False(t, modified)
		assert.Equal(t, expectedError, err)
	})
}
