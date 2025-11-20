package video_tasks

import (
	"context"
	"encoding/json"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMetadataDesc(t *testing.T) {
	result := MetadataDesc(123, "test.mp4")
	assert.Equal(t, "Update metadata for test.mp4", result)
}

func TestMarshalVideoMetadataParams(t *testing.T) {
	params, err := MarshalVideoMetadataParams(123)
	assert.NoError(t, err)

	var p videoMetadataParams
	err = json.Unmarshal([]byte(params), &p)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), p.ItemId)
}

func TestUnmarshalVideoMetadataParams(t *testing.T) {
	p := videoMetadataParams{ItemId: 456}
	data, _ := json.Marshal(p)

	result, err := unmarshalVideoMetadataParams(string(data))
	assert.NoError(t, err)
	assert.Equal(t, uint64(456), result.ItemId)
}

func TestUnmarshalVideoMetadataParams_InvalidJSON(t *testing.T) {
	_, err := unmarshalVideoMetadataParams("invalid json")
	assert.Error(t, err)
}

func TestUpdateVideoMetadata_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	err := UpdateVideoMetadata(ctx, mockIRW, "invalid json")
	assert.Error(t, err)
}

func TestUpdateVideoMetadata_GetItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(nil, assert.AnError)

	params, _ := MarshalVideoMetadataParams(123)
	err := UpdateVideoMetadata(ctx, mockIRW, params)
	assert.Error(t, err)
}

func TestUpdateVideoMetadata_SubItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	mainItemId := uint64(100)
	subItem := &model.Item{
		Id:            123,
		MainItemId:    &mainItemId,
		StartPosition: 10.0,
		EndPosition:   20.0,
	}

	mainItem := &model.Item{
		Id:              100,
		Width:           1920,
		Height:          1080,
		VideoCodecName:  "h264",
		AudioCodecName:  "aac",
		DurationSeconds: 100.0,
	}

	mockIRW.EXPECT().GetItem(gomock.Any(), uint64(123)).Return(subItem, nil)
	mockIRW.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(mainItem, nil)
	mockIRW.EXPECT().UpdateItem(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, item *model.Item) error {
		assert.Equal(t, 10.0, item.DurationSeconds) // end - start
		assert.Equal(t, 1920, item.Width)
		assert.Equal(t, 1080, item.Height)
		assert.Equal(t, "h264", item.VideoCodecName)
		assert.Equal(t, "aac", item.AudioCodecName)
		return nil
	})

	params, _ := MarshalVideoMetadataParams(123)
	err := UpdateVideoMetadata(ctx, mockIRW, params)
	assert.NoError(t, err)
}

func TestUpdateVideoMetadata_SubItem_GetMainItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	mainItemId := uint64(100)
	subItem := &model.Item{
		Id:         123,
		MainItemId: &mainItemId,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(subItem, nil)
	mockIRW.EXPECT().GetItem(ctx, gomock.Any()).Return(nil, assert.AnError)

	params, _ := MarshalVideoMetadataParams(123)
	err := UpdateVideoMetadata(ctx, mockIRW, params)
	assert.Error(t, err)
}

func TestUpdateVideoMetadata_Highlight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	mainItemId := uint64(100)
	highlightItem := &model.Item{
		Id:                   123,
		HighlightParentItemId: &mainItemId,
		StartPosition:        10.0,
		EndPosition:          20.0,
	}

	mainItem := &model.Item{
		Id:              100,
		Width:           1920,
		Height:          1080,
		VideoCodecName:  "h264",
		AudioCodecName:  "aac",
		DurationSeconds: 100.0,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(highlightItem, nil)
	mockIRW.EXPECT().GetItem(ctx, gomock.Any()).Return(mainItem, nil)
	mockIRW.EXPECT().UpdateItem(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, item *model.Item) error {
		assert.Equal(t, 10.0, item.DurationSeconds) // end - start
		return nil
	})

	params, _ := MarshalVideoMetadataParams(123)
	err := UpdateVideoMetadata(ctx, mockIRW, params)
	assert.NoError(t, err)
}

func TestUpdateVideoMetadata_UpdateItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	ctx := context.Background()

	mainItemId := uint64(100)
	subItem := &model.Item{
		Id:         123,
		MainItemId: &mainItemId,
	}

	mainItem := &model.Item{
		Id: 100,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(subItem, nil)
	mockIRW.EXPECT().GetItem(ctx, gomock.Any()).Return(mainItem, nil)
	mockIRW.EXPECT().UpdateItem(ctx, gomock.Any()).Return(assert.AnError)

	params, _ := MarshalVideoMetadataParams(123)
	err := UpdateVideoMetadata(ctx, mockIRW, params)
	assert.Error(t, err)
}

func TestUpdateNonMainItemMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIR := model.NewMockItemReader(ctrl)
	ctx := context.Background()

	mainItemId := uint64(100)
	subItem := &model.Item{
		Id:         123,
		MainItemId: &mainItemId,
		StartPosition: 10.0,
		EndPosition:   20.0,
	}

	mainItem := &model.Item{
		Id:              100,
		Width:           1920,
		Height:          1080,
		VideoCodecName:  "h264",
		AudioCodecName:  "aac",
	}

	mockIR.EXPECT().GetItem(ctx, gomock.Any()).Return(mainItem, nil)

	err := updateNonMainItemMetadata(ctx, mockIR, subItem)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, subItem.DurationSeconds)
	assert.Equal(t, 1920, subItem.Width)
	assert.Equal(t, 1080, subItem.Height)
	assert.Equal(t, "h264", subItem.VideoCodecName)
	assert.Equal(t, "aac", subItem.AudioCodecName)
}

func TestUpdateNonMainItemMetadata_GetItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIR := model.NewMockItemReader(ctrl)
	ctx := context.Background()

	mainItemId := uint64(100)
	subItem := &model.Item{
		Id:         123,
		MainItemId: &mainItemId,
	}

	mockIR.EXPECT().GetItem(ctx, gomock.Any()).Return(nil, assert.AnError)

	err := updateNonMainItemMetadata(ctx, mockIR, subItem)
	assert.Error(t, err)
}

