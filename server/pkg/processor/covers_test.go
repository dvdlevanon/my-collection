package processor

import (
	"my-collection/server/pkg/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshMainCover(t *testing.T) {
	// happy path
	irw := newTestSingleItemReaderWriter(&model.Item{Id: 0, Url: sampleMp4}, false, false)
	tempfile, err := os.CreateTemp("", "refresh-main-cover-*.png")
	assert.NoError(t, err)
	uploader := newTestStorageUploader(tempfile.Name(), "main-cover-url", "", false)
	assert.NoError(t, refreshMainCover(&irw, &uploader, 0, 4.5))
	assert.Equal(t, uploader.storageUrl, *irw.item.MainCoverUrl)
	actualBytes, err := os.ReadFile(uploader.fileForWritting)
	assert.NoError(t, err)
	expectedBytes, err := os.ReadFile(sample4_5SecondsScreenshotPng)
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, actualBytes)

	// error getting/updating
	irw.errorGet = true
	irw.errorSet = false
	assert.Error(t, refreshMainCover(&irw, &uploader, 0, 4.5))
	irw.errorGet = false
	irw.errorSet = true
	assert.Error(t, refreshMainCover(&irw, &uploader, 0, 4.5))
	irw.errorGet = false
	irw.errorSet = false

	// error uploading
	uploader.errorWriting = true
	assert.Error(t, refreshMainCover(&irw, &uploader, 0, 4.5))
	uploader.errorWriting = false

	// bad video file
	irw.UpdateItem(&model.Item{Id: 0, Url: "fsdafsa"})
	assert.Error(t, refreshMainCover(&irw, &uploader, 0, 4.5))
}
