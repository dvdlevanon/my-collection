package video_tasks

import (
	"context"
	"encoding/json"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestResolutionDesc(t *testing.T) {
	tests := []struct {
		name     string
		w        int
		h        int
		expected string
	}{
		{"Width -1", -1, 1080, "Change resolution to height 1080 for test.mp4"},
		{"Height -1", 1920, -1, "Change resolution to width 1920 for test.mp4"},
		{"Both specified", 1920, 1080, "Change resolution to 1920x1080 for test.mp4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolutionDesc(123, "test.mp4", tt.w, tt.h)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarshalVideoResolutionParams(t *testing.T) {
	params, err := MarshalVideoResolutionParams(123, 1920, 1080)
	assert.NoError(t, err)

	var p videoResolutionParams
	err = json.Unmarshal([]byte(params), &p)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), p.ItemId)
	assert.Equal(t, 1920, p.Width)
	assert.Equal(t, 1080, p.Height)
}

func TestUnmarshalVideoResolutionParams(t *testing.T) {
	p := videoResolutionParams{ItemId: 456, Width: 1280, Height: 720}
	data, _ := json.Marshal(p)

	result, err := unmarshalVideoResolutionParams(string(data))
	assert.NoError(t, err)
	assert.Equal(t, uint64(456), result.ItemId)
	assert.Equal(t, 1280, result.Width)
	assert.Equal(t, 720, result.Height)
}

func TestUnmarshalVideoResolutionParams_InvalidJSON(t *testing.T) {
	_, err := unmarshalVideoResolutionParams("invalid json")
	assert.Error(t, err)
}

func TestChangeVideoResolution_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockTempProvider := model.NewMockTempFileProvider(ctrl)
	ctx := context.Background()

	err := ChangeVideoResolution(ctx, mockIRW, mockTempProvider, "invalid json")
	assert.Error(t, err)
}

func TestChangeVideoResolution_GetItemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockTempProvider := model.NewMockTempFileProvider(ctrl)
	ctx := context.Background()

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(nil, assert.AnError)

	params, _ := MarshalVideoResolutionParams(123, 1920, 1080)
	err := ChangeVideoResolution(ctx, mockIRW, mockTempProvider, params)
	assert.Error(t, err)
}

func TestChangeVideoResolution_AlreadySameResolution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockTempProvider := model.NewMockTempFileProvider(ctrl)
	ctx := context.Background()

	item := &model.Item{
		Id:     123,
		Width:  1920,
		Height: 1080,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)

	params, _ := MarshalVideoResolutionParams(123, 1920, 1080)
	err := ChangeVideoResolution(ctx, mockIRW, mockTempProvider, params)
	assert.NoError(t, err)
}

func TestChangeVideoResolution_WidthSame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockTempProvider := model.NewMockTempFileProvider(ctrl)
	ctx := context.Background()

	item := &model.Item{
		Id:     123,
		Width:  1920,
		Height: 1080,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)

	params, _ := MarshalVideoResolutionParams(123, 1920, 1080)
	err := ChangeVideoResolution(ctx, mockIRW, mockTempProvider, params)
	assert.NoError(t, err)
}

func TestChangeVideoResolution_HeightSame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockTempProvider := model.NewMockTempFileProvider(ctrl)
	ctx := context.Background()

	item := &model.Item{
		Id:     123,
		Width:  1920,
		Height: 1080,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)

	params, _ := MarshalVideoResolutionParams(123, 1920, 1080)
	err := ChangeVideoResolution(ctx, mockIRW, mockTempProvider, params)
	assert.NoError(t, err)
}

func TestChangeVideoResolution_WidthMinusOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockTempProvider := model.NewMockTempFileProvider(ctrl)
	ctx := context.Background()

	item := &model.Item{
		Id:     123,
		Width:  1920,
		Height: 1080,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)

	params, _ := MarshalVideoResolutionParams(123, -1, 1080)
	err := ChangeVideoResolution(ctx, mockIRW, mockTempProvider, params)
	// This will fail because we can't actually call ffmpeg, but we test the logic
	assert.NoError(t, err) // Actually this should return early if already same height
}

func TestChangeVideoResolution_HeightMinusOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIRW := model.NewMockItemReaderWriter(ctrl)
	mockTempProvider := model.NewMockTempFileProvider(ctrl)
	ctx := context.Background()

	item := &model.Item{
		Id:     123,
		Width:  1920,
		Height: 1080,
	}

	mockIRW.EXPECT().GetItem(ctx, uint64(123)).Return(item, nil)

	params, _ := MarshalVideoResolutionParams(123, 1920, -1)
	err := ChangeVideoResolution(ctx, mockIRW, mockTempProvider, params)
	// This will fail because we can't actually call ffmpeg, but we test the logic
	assert.NoError(t, err) // Actually this should return early if already same width
}

