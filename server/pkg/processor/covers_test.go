package processor

import (
	"my-collection/server/pkg/model"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestRefreshMainCover(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// init mocks
	item := model.Item{Id: 0, Url: sampleMp4}
	irw := model.NewMockItemReaderWriter(ctrl)
	tempfile, err := os.CreateTemp("", "refresh-main-cover-*.png")
	assert.NoError(t, err)
	uploader := model.NewMockStorageUploader(ctrl)

	// happy path
	irw.EXPECT().GetItem(gomock.Any()).Return(&item, nil)
	irw.EXPECT().UpdateItem(gomock.Any()).Return(nil)
	uploader.EXPECT().GetFileForWriting(gomock.Any()).Return(tempfile.Name(), nil)
	uploader.EXPECT().GetStorageUrl(gomock.Any()).Return("main-cover-url")
	assert.NoError(t, refreshMainCover(irw, uploader, 0, 4.5))
	assert.Equal(t, "main-cover-url", *item.MainCoverUrl)

	// Find a way to compare PNGs and ignore metadata such as dates
	//
	// var sample4_5SecondsScreenshotPng = filepath.Join(testFiles, "sample-4_5-second-screenshot.png")
	// actualBytes, err := os.ReadFile(tempfile.Name())
	// assert.NoError(t, err)
	// expectedBytes, err := os.ReadFile(sample4_5SecondsScreenshotPng)
	// assert.NoError(t, err)
	// assert.Equal(t, expectedBytes, actualBytes)

	// error getting/updating
	irw.EXPECT().GetItem(gomock.Any()).Return(&item, errors.Errorf("test error"))
	assert.Error(t, refreshMainCover(irw, uploader, 0, 4.5))
	irw.EXPECT().GetItem(gomock.Any()).Return(&item, nil)
	irw.EXPECT().UpdateItem(gomock.Any()).Return(errors.Errorf("test error"))
	uploader.EXPECT().GetFileForWriting(gomock.Any()).Return(tempfile.Name(), nil)
	uploader.EXPECT().GetStorageUrl(gomock.Any()).Return("main-cover-url")
	assert.Error(t, refreshMainCover(irw, uploader, 0, 4.5))

	// error uploading
	irw.EXPECT().GetItem(gomock.Any()).Return(&item, nil)
	uploader.EXPECT().GetFileForWriting(gomock.Any()).Return(tempfile.Name(), errors.Errorf("test error"))
	assert.Error(t, refreshMainCover(irw, uploader, 0, 4.5))

	// bad video file
	irw.EXPECT().GetItem(gomock.Any()).Return(&item, nil)
	uploader.EXPECT().GetFileForWriting(gomock.Any()).Return(tempfile.Name(), errors.Errorf("test error"))
	assert.Error(t, refreshMainCover(irw, uploader, 0, 4.5))
}
