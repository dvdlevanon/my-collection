package video_tasks

import (
	"context"
	"encoding/json"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPreviewDesc(t *testing.T) {
	result := PreviewDesc(123, "test.mp4", 5, 10)
	assert.Equal(t, "Generate preview with 5 scenes (10s each) for test.mp4", result)
}

func TestMarshalVideoPreviewParams(t *testing.T) {
	params, err := MarshalVideoPreviewParams(123, 5, 10)
	assert.NoError(t, err)

	var p videoPreviewParams
	err = json.Unmarshal([]byte(params), &p)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), p.ItemId)
	assert.Equal(t, 5, p.SceneCount)
	assert.Equal(t, 10, p.SceneDuration)
}

func TestUnmarshalVideoPreviewParams(t *testing.T) {
	p := videoPreviewParams{ItemId: 456, SceneCount: 8, SceneDuration: 15}
	data, _ := json.Marshal(p)

	result, err := unmarshalVideoPreviewParams(string(data))
	assert.NoError(t, err)
	assert.Equal(t, uint64(456), result.ItemId)
	assert.Equal(t, 8, result.SceneCount)
	assert.Equal(t, 15, result.SceneDuration)
}

func TestUnmarshalVideoPreviewParams_InvalidJSON(t *testing.T) {
	_, err := unmarshalVideoPreviewParams("invalid json")
	assert.Error(t, err)
}

func TestRefreshVideoPreview_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	err := RefreshVideoPreview(ctx, mockIRW, mockUploader, "invalid json")
	assert.Error(t, err)
}

func TestRefreshVideoPreview_GetItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(nil, assert.AnError)

	params, _ := MarshalVideoPreviewParams(123, 5, 10)
	err := RefreshVideoPreview(ctx, mockIRW, mockUploader, params)
	assert.Error(t, err)
}

func TestRefreshVideoPreview_PreviewFromStartPosition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	item := &model.Item{
		Id:          123,
		PreviewMode: items.PREVIEW_FROM_START_POSITION,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)

	params, _ := MarshalVideoPreviewParams(123, 5, 10)
	err := RefreshVideoPreview(ctx, mockIRW, mockUploader, params)
	assert.NoError(t, err)
}

func TestRefreshVideoPreview_ZeroSceneCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	item := &model.Item{
		Id: 123,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)

	params, _ := MarshalVideoPreviewParams(123, 0, 10)
	err := RefreshVideoPreview(ctx, mockIRW, mockUploader, params)
	assert.NoError(t, err)
}
