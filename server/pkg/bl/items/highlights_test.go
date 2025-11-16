package items

import (
	"context"
	"errors"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"k8s.io/utils/pointer"
)

func TestInitHighlights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReaderWriter := model.NewMockTagReaderWriter(ctrl)

	t.Run("highlights tag already exists", func(t *testing.T) {
		existingTag := &model.Tag{
			Id:       123,
			Title:    "Highlights",
			ParentID: nil,
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), gomock.Any()).
			Return(existingTag, nil)

		err := InitHighlights(context.Background(), mockTagReaderWriter)

		assert.NoError(t, err)
		// Verify that the global highlightsTag was updated
		assert.Equal(t, existingTag.Id, GetHighlightsTagId())
	})

	t.Run("highlights tag does not exist - create new", func(t *testing.T) {
		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("tag not found"))
		mockTagReaderWriter.EXPECT().
			CreateOrUpdateTag(gomock.Any(), gomock.Any()).
			Return(nil)

		err := InitHighlights(context.Background(), mockTagReaderWriter)

		assert.NoError(t, err)
	})

	t.Run("error creating highlights tag", func(t *testing.T) {
		expectedError := errors.New("create tag error")

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("tag not found"))
		mockTagReaderWriter.EXPECT().
			CreateOrUpdateTag(gomock.Any(), gomock.Any()).
			Return(expectedError)

		err := InitHighlights(context.Background(), mockTagReaderWriter)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestBuildHighlightUrl(t *testing.T) {
	tests := []struct {
		name          string
		origin        string
		startPosition float64
		endPosition   float64
		expected      string
	}{
		{
			name:          "basic highlight url",
			origin:        "test/video.mp4",
			startPosition: 10.5,
			endPosition:   20.5,
			expected:      "test/video.mp4-10.500000-20.500000",
		},
		{
			name:          "zero start position",
			origin:        "video.mp4",
			startPosition: 0.0,
			endPosition:   15.0,
			expected:      "video.mp4-0.000000-15.000000",
		},
		{
			name:          "fractional positions",
			origin:        "path/to/video.mp4",
			startPosition: 1.23456,
			endPosition:   7.89012,
			expected:      "path/to/video.mp4-1.234560-7.890120",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildHighlightUrl(tt.origin, tt.startPosition, tt.endPosition)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildHighlight(t *testing.T) {
	originalItem := &model.Item{
		Id:              1,
		Title:           "Test Video",
		Origin:          "test/origin",
		Url:             "test/url/video.mp4",
		Width:           1920,
		Height:          1080,
		VideoCodecName:  "h264",
		AudioCodecName:  "aac",
		LastModified:    123456789,
		DurationSeconds: 100.0,
	}

	startPosition := 10.0
	endPosition := 30.0
	highlightId := uint64(456)

	result := buildHighlight(originalItem, startPosition, endPosition, highlightId)

	assert.Equal(t, originalItem.Title, result.Title)
	assert.Equal(t, "test/origin-10.000000-30.000000", result.Origin)
	assert.Equal(t, originalItem.Url, result.Url)
	assert.Equal(t, startPosition, result.StartPosition)
	assert.Equal(t, endPosition, result.EndPosition)
	assert.Equal(t, originalItem.Width, result.Width)
	assert.Equal(t, originalItem.Height, result.Height)
	assert.Equal(t, endPosition-startPosition, result.DurationSeconds)
	assert.Equal(t, originalItem.VideoCodecName, result.VideoCodecName)
	assert.Equal(t, originalItem.AudioCodecName, result.AudioCodecName)
	assert.Equal(t, originalItem.LastModified, result.LastModified)
	assert.Equal(t, PREVIEW_FROM_START_POSITION, result.PreviewMode)
	assert.Len(t, result.Tags, 1)
	assert.Equal(t, highlightId, result.Tags[0].Id)
}

func TestMakeHighlight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemReaderWriter := model.NewMockItemReaderWriter(ctrl)

	t.Run("successful highlight creation", func(t *testing.T) {
		itemId := uint64(123)
		startPosition := 10.0
		endPosition := 30.0
		highlightId := uint64(456)

		originalItem := &model.Item{
			Id:              itemId,
			Title:           "Test Video",
			Origin:          "test/origin",
			Url:             "test/url/video.mp4",
			Width:           1920,
			Height:          1080,
			VideoCodecName:  "h264",
			AudioCodecName:  "aac",
			LastModified:    123456789,
			DurationSeconds: 100.0,
			Highlights:      []*model.Item{},
		}

		mockItemReaderWriter.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(originalItem, nil)
		mockItemReaderWriter.EXPECT().
			CreateOrUpdateItem(gomock.Any(), gomock.Any()).
			Return(nil)
		mockItemReaderWriter.EXPECT().
			UpdateItem(gomock.Any(), originalItem).
			Return(nil)

		result, err := MakeHighlight(context.Background(), mockItemReaderWriter, itemId, startPosition, endPosition, highlightId)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, originalItem.Title, result.Title)
		assert.Equal(t, startPosition, result.StartPosition)
		assert.Equal(t, endPosition, result.EndPosition)
		assert.Equal(t, endPosition-startPosition, result.DurationSeconds)
		assert.Len(t, result.Tags, 1)
		assert.Equal(t, highlightId, result.Tags[0].Id)

		// Verify the original item was updated with the highlight
		assert.Len(t, originalItem.Highlights, 1)
		assert.Equal(t, result, originalItem.Highlights[0])
	})

	t.Run("error getting item", func(t *testing.T) {
		itemId := uint64(123)
		startPosition := 10.0
		endPosition := 30.0
		highlightId := uint64(456)
		expectedError := errors.New("item not found")

		mockItemReaderWriter.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(nil, expectedError)

		result, err := MakeHighlight(context.Background(), mockItemReaderWriter, itemId, startPosition, endPosition, highlightId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("error creating highlight item", func(t *testing.T) {
		itemId := uint64(123)
		startPosition := 10.0
		endPosition := 30.0
		highlightId := uint64(456)
		expectedError := errors.New("create highlight error")

		originalItem := &model.Item{
			Id:         itemId,
			Title:      "Test Video",
			Highlights: []*model.Item{},
		}

		mockItemReaderWriter.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(originalItem, nil)
		mockItemReaderWriter.EXPECT().
			CreateOrUpdateItem(gomock.Any(), gomock.Any()).
			Return(expectedError)

		result, err := MakeHighlight(context.Background(), mockItemReaderWriter, itemId, startPosition, endPosition, highlightId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("error updating original item", func(t *testing.T) {
		itemId := uint64(123)
		startPosition := 10.0
		endPosition := 30.0
		highlightId := uint64(456)
		expectedError := errors.New("update item error")

		originalItem := &model.Item{
			Id:         itemId,
			Title:      "Test Video",
			Highlights: []*model.Item{},
		}

		mockItemReaderWriter.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(originalItem, nil)
		mockItemReaderWriter.EXPECT().
			CreateOrUpdateItem(gomock.Any(), gomock.Any()).
			Return(nil)
		mockItemReaderWriter.EXPECT().
			UpdateItem(gomock.Any(), originalItem).
			Return(expectedError)

		result, err := MakeHighlight(context.Background(), mockItemReaderWriter, itemId, startPosition, endPosition, highlightId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestIsHighlight(t *testing.T) {
	t.Run("item is highlight", func(t *testing.T) {
		item := &model.Item{
			Id:                    1,
			HighlightParentItemId: pointer.Uint64(123),
		}

		result := IsHighlight(item)
		assert.True(t, result)
	})

	t.Run("item is not highlight", func(t *testing.T) {
		item := &model.Item{
			Id:                    1,
			HighlightParentItemId: nil,
		}

		result := IsHighlight(item)
		assert.False(t, result)
	})
}

func TestGetHighlightsTagId(t *testing.T) {
	// Set up the highlights tag with a known ID
	highlightTagId := uint64(789)
	highlightsTag.Id = highlightTagId

	result := GetHighlightsTagId()
	assert.Equal(t, highlightTagId, result)
}

// Integration test to verify the full workflow
func TestHighlightsIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagReaderWriter := model.NewMockTagReaderWriter(ctrl)
	mockItemReaderWriter := model.NewMockItemReaderWriter(ctrl)

	t.Run("full highlights workflow", func(t *testing.T) {
		// Step 1: Initialize highlights
		existingHighlightsTag := &model.Tag{
			Id:       999,
			Title:    "Highlights",
			ParentID: nil,
		}

		mockTagReaderWriter.EXPECT().
			GetTag(gomock.Any(), gomock.Any()).
			Return(existingHighlightsTag, nil)

		err := InitHighlights(context.Background(), mockTagReaderWriter)
		assert.NoError(t, err)

		// Step 2: Create a highlight
		itemId := uint64(123)
		startPosition := 15.0
		endPosition := 45.0
		highlightId := existingHighlightsTag.Id

		originalItem := &model.Item{
			Id:              itemId,
			Title:           "Original Video",
			Origin:          "videos/original",
			Url:             "videos/original/video.mp4",
			Width:           1920,
			Height:          1080,
			VideoCodecName:  "h264",
			AudioCodecName:  "aac",
			LastModified:    123456789,
			DurationSeconds: 120.0,
			Highlights:      []*model.Item{},
		}

		mockItemReaderWriter.EXPECT().
			GetItem(gomock.Any(), itemId).
			Return(originalItem, nil)
		mockItemReaderWriter.EXPECT().
			CreateOrUpdateItem(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, item *model.Item) error {
				// Verify the highlight was built correctly
				assert.Equal(t, originalItem.Title, item.Title)
				assert.Equal(t, "videos/original-15.000000-45.000000", item.Origin)
				assert.Equal(t, originalItem.Url, item.Url)
				assert.Equal(t, startPosition, item.StartPosition)
				assert.Equal(t, endPosition, item.EndPosition)
				assert.Equal(t, endPosition-startPosition, item.DurationSeconds)
				assert.Equal(t, PREVIEW_FROM_START_POSITION, item.PreviewMode)
				assert.Len(t, item.Tags, 1)
				assert.Equal(t, highlightId, item.Tags[0].Id)
				return nil
			})
		mockItemReaderWriter.EXPECT().
			UpdateItem(gomock.Any(), originalItem).
			Return(nil)

		highlight, err := MakeHighlight(context.Background(), mockItemReaderWriter, itemId, startPosition, endPosition, highlightId)

		assert.NoError(t, err)
		assert.NotNil(t, highlight)

		// Step 3: Verify highlight properties
		// Note: The highlight created by buildHighlight doesn't automatically set HighlightParentItemId
		// This would typically be set when the highlight is associated with its parent in the database
		assert.Equal(t, GetHighlightsTagId(), highlightId)

		// Verify original item was updated
		assert.Len(t, originalItem.Highlights, 1)
		assert.Equal(t, highlight, originalItem.Highlights[0])
	})
}
