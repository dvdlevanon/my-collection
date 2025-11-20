package video_tasks

import (
	"context"
	"encoding/json"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCoversDesc(t *testing.T) {
	result := CoversDesc(123, "test.mp4", 5)
	assert.Equal(t, "Extract 5 covers for test.mp4", result)
}

func TestMarshalVideoCoversParams(t *testing.T) {
	params, err := MarshalVideoCoversParams(123, 5)
	assert.NoError(t, err)

	var p videoCoversParams
	err = json.Unmarshal([]byte(params), &p)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), p.ItemId)
	assert.Equal(t, 5, p.Count)
}

func TestUnmarshalVideoCoversParams(t *testing.T) {
	p := videoCoversParams{ItemId: 456, Count: 10}
	data, _ := json.Marshal(p)

	result, err := unmarshalVideoCoversParams(string(data))
	assert.NoError(t, err)
	assert.Equal(t, uint64(456), result.ItemId)
	assert.Equal(t, 10, result.Count)
}

func TestUnmarshalVideoCoversParams_InvalidJSON(t *testing.T) {
	_, err := unmarshalVideoCoversParams("invalid json")
	assert.Error(t, err)
}

func TestRefreshVideoCovers_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	err := RefreshVideoCovers(ctx, mockIRW, mockUploader, "invalid json")
	assert.Error(t, err)
}

func TestRefreshVideoCovers_ZeroCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	params, _ := MarshalVideoCoversParams(123, 0)
	err := RefreshVideoCovers(ctx, mockIRW, mockUploader, params)
	assert.NoError(t, err)
}

func TestRefreshVideoCovers_GetItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(nil, assert.AnError)

	params, _ := MarshalVideoCoversParams(123, 5)
	err := RefreshVideoCovers(ctx, mockIRW, mockUploader, params)
	assert.Error(t, err)
}

func TestGetDurationForItem_SubItem(t *testing.T) {
	item := &model.Item{
		Id:              123,
		MainItemId:      func() *uint64 { id := uint64(100); return &id }(),
		DurationSeconds: 50.0,
	}

	duration, err := getDurationForItem(item, "/path/to/video.mp4")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, duration)
}

func TestGetDurationForItem_Highlight(t *testing.T) {
	item := &model.Item{
		Id:                   123,
		HighlightParentItemId: func() *uint64 { id := uint64(100); return &id }(),
		DurationSeconds:       30.0,
	}

	duration, err := getDurationForItem(item, "/path/to/video.mp4")
	assert.NoError(t, err)
	assert.Equal(t, 30.0, duration)
}

