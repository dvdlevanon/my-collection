package video_tasks

import (
	"context"
	"encoding/json"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMainCoverDesc(t *testing.T) {
	result := MainCoverDesc(123, "test.mp4", 5.5)
	assert.Equal(t, "Set main cover at 5.50s for test.mp4", result)
}

func TestMarshalMainCoverParams(t *testing.T) {
	params, err := MarshalMainCoverParams(123, 5.5)
	assert.NoError(t, err)

	var p videoMainCoverParams
	err = json.Unmarshal([]byte(params), &p)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), p.ItemId)
	assert.Equal(t, 5.5, p.Second)
}

func TestUnmarshalMainCoverParams(t *testing.T) {
	p := videoMainCoverParams{ItemId: 456, Second: 10.5}
	data, _ := json.Marshal(p)

	result, err := unmarshalMainCoverParams(string(data))
	assert.NoError(t, err)
	assert.Equal(t, uint64(456), result.ItemId)
	assert.Equal(t, 10.5, result.Second)
}

func TestUnmarshalMainCoverParams_InvalidJSON(t *testing.T) {
	_, err := unmarshalMainCoverParams("invalid json")
	assert.Error(t, err)
}

func TestUpdateVideoMainCover_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	err := UpdateVideoMainCover(ctx, mockIRW, mockUploader, "invalid json")
	assert.Error(t, err)
}

func TestUpdateVideoMainCover_GetItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(nil, assert.AnError)

	params, _ := MarshalMainCoverParams(123, 5.5)
	err := UpdateVideoMainCover(ctx, mockIRW, mockUploader, params)
	assert.Error(t, err)
}

