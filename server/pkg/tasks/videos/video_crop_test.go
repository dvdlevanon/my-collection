package video_tasks

import (
	"context"
	"encoding/json"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCropDesc(t *testing.T) {
	rect := model.RectFloat{X: 10, Y: 20, W: 100, H: 200}
	result := CropDesc(123, "test.mp4", 5.5, rect)
	assert.Contains(t, result, "Cropp frame at 5.50s")
	assert.Contains(t, result, "test.mp4")
	assert.Contains(t, result, "10,20")
	assert.Contains(t, result, "100x200")
}

func TestMarshalVideoCropParams(t *testing.T) {
	rect := model.RectFloat{X: 10, Y: 20, W: 100, H: 200}
	params, err := MarshalVideoCropParams(123, 5.5, rect)
	assert.NoError(t, err)

	var p videoCropParams
	err = json.Unmarshal([]byte(params), &p)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), p.ItemId)
	assert.Equal(t, 5.5, p.Second)
	assert.Equal(t, 10.0, p.X)
	assert.Equal(t, 20.0, p.Y)
	assert.Equal(t, 100.0, p.W)
	assert.Equal(t, 200.0, p.H)
}

func TestUnmarshalVideoCropParams(t *testing.T) {
	p := videoCropParams{
		ItemId: 456,
		Second: 10.5,
		X:      15.0,
		Y:      25.0,
		W:      200.0,
		H:      300.0,
	}
	data, _ := json.Marshal(p)

	result, err := unmarshalVideoCropParams(string(data))
	assert.NoError(t, err)
	assert.Equal(t, uint64(456), result.ItemId)
	assert.Equal(t, 10.5, result.Second)
	assert.Equal(t, 15.0, result.X)
	assert.Equal(t, 25.0, result.Y)
	assert.Equal(t, 200.0, result.W)
	assert.Equal(t, 300.0, result.H)
}

func TestUnmarshalVideoCropParams_InvalidJSON(t *testing.T) {
	_, err := unmarshalVideoCropParams("invalid json")
	assert.Error(t, err)
}

func TestCropVideoFrame_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	err := CropVideoFrame(ctx, mockIRW, mockUploader, "invalid json")
	assert.Error(t, err)
}

func TestCropVideoFrame_GetItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockUploader := model.NewMockStorageUploader(ctrl)
	ctx := context.Background()

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(nil, assert.AnError)

	rect := model.RectFloat{X: 10, Y: 20, W: 100, H: 200}
	params, _ := MarshalVideoCropParams(123, 5.5, rect)
	err := CropVideoFrame(ctx, mockIRW, mockUploader, params)
	assert.Error(t, err)
}

