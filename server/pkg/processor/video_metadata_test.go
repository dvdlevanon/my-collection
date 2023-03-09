package processor

import (
	"my-collection/server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshVideoMetadata(t *testing.T) {
	// happy path
	irw := newTestSingleItemReaderWriter(&model.Item{Id: 0, Url: sampleMp4}, false, false)
	assert.NoError(t, refreshItemMetadata(&irw, 0))
	updated, _ := irw.GetItem()
	assert.Equal(t, 5, updated.DurationSeconds)
	assert.Equal(t, 560, updated.Width)
	assert.Equal(t, 320, updated.Height)
	assert.Equal(t, "h264", updated.VideoCodecName)
	assert.Equal(t, "aac", updated.AudioCodecName)

	// invalid duration
	irw.UpdateItem(&model.Item{Id: 0, Url: sample3SecondsScreenshotPng})
	assert.Error(t, refreshItemMetadata(&irw, 0))

	// no video
	irw.UpdateItem(&model.Item{Id: 0, Url: sampleNoVideoMp4})
	assert.Error(t, refreshItemMetadata(&irw, 0))

	// no audio
	irw.UpdateItem(&model.Item{Id: 0, Url: sampleNoAudioMp4})
	assert.Error(t, refreshItemMetadata(&irw, 0))

	// error getting/updating item
	irwErrorGet := newTestSingleItemReaderWriter(&model.Item{Id: 0, Url: sampleMp4}, true, false)
	irwErrorSet := newTestSingleItemReaderWriter(&model.Item{Id: 0, Url: sampleMp4}, true, false)
	assert.Error(t, refreshItemMetadata(&irwErrorGet, 0))
	assert.Error(t, refreshItemMetadata(&irwErrorSet, 0))
}
