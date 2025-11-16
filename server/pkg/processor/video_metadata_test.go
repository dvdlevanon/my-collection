package processor

import (
	"context"
	"my-collection/server/pkg/model"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestRefreshVideoMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// init mocks
	item := model.Item{Id: 0, Url: sampleMp4}
	irw := model.NewMockItemReaderWriter(ctrl)

	// happy path
	irw.EXPECT().GetItem(context.Background(), gomock.Any()).Return(&item, nil)
	irw.EXPECT().UpdateItem(context.Background(), gomock.Any()).Return(nil)
	assert.NoError(t, refreshItemMetadata(context.Background(), irw, 0))
	assert.Equal(t, 5.568, item.DurationSeconds)
	assert.Equal(t, 560, item.Width)
	assert.Equal(t, 320, item.Height)
	assert.Equal(t, "h264", item.VideoCodecName)
	assert.Equal(t, "aac", item.AudioCodecName)

	// invalid duration
	irw.EXPECT().GetItem(context.Background(), gomock.Any()).Return(&model.Item{Url: sample3SecondsScreenshotPng}, nil)
	assert.Error(t, refreshItemMetadata(context.Background(), irw, 0))

	// no video
	irw.EXPECT().GetItem(context.Background(), gomock.Any()).Return(&model.Item{Url: sampleNoVideoMp4}, nil)
	assert.Error(t, refreshItemMetadata(context.Background(), irw, 0))

	// no audio
	irw.EXPECT().GetItem(context.Background(), gomock.Any()).Return(&model.Item{Url: sampleNoAudioMp4}, nil)
	assert.Error(t, refreshItemMetadata(context.Background(), irw, 0))

	// error getting/updating item
	irw.EXPECT().GetItem(context.Background(), gomock.Any()).Return(&item, errors.Errorf("test error"))
	assert.Error(t, refreshItemMetadata(context.Background(), irw, 0))
	irw.EXPECT().UpdateItem(context.Background(), gomock.Any()).Return(errors.Errorf("test error"))
	irw.EXPECT().GetItem(context.Background(), gomock.Any()).Return(&item, nil)
	assert.Error(t, refreshItemMetadata(context.Background(), irw, 0))
}
