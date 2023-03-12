package processor

import (
	"my-collection/server/pkg/model"
	"path/filepath"

	"github.com/go-errors/errors"
)

var testFiles = "../../test-files"
var sampleMp4 = filepath.Join(testFiles, "sample.mp4")
var sampleNoAudioMp4 = filepath.Join(testFiles, "sample-no-audio.mp4")
var sampleNoVideoMp4 = filepath.Join(testFiles, "sample-no-video.mp4")
var sample3SecondsScreenshotPng = filepath.Join(testFiles, "sample-3-second-screenshot.png")
var sample4_5SecondsScreenshotPng = filepath.Join(testFiles, "sample-4_5-second-screenshot.png")

func newTestSingleItemReaderWriter(item *model.Item, errorGet bool, errorSet bool) testSingleItemReaderWriter {
	return testSingleItemReaderWriter{
		item:     item,
		errorGet: errorGet,
		errorSet: errorSet,
	}
}

type testSingleItemReaderWriter struct {
	item     *model.Item
	errorGet bool
	errorSet bool
}

func (t *testSingleItemReaderWriter) GetAllItems() (*[]model.Item, error)                 { return nil, nil }
func (t *testSingleItemReaderWriter) CreateOrUpdateItem(item *model.Item) error           { return nil }
func (t *testSingleItemReaderWriter) RemoveItem(itemId uint64) error                      { return nil }
func (t *testSingleItemReaderWriter) RemoveTagFromItem(itemId uint64, tagId uint64) error { return nil }
func (t *testSingleItemReaderWriter) GetItems(conds ...interface{}) (*[]model.Item, error) {
	return nil, nil
}

func (t *testSingleItemReaderWriter) GetItem(conds ...interface{}) (*model.Item, error) {
	if t.errorGet {
		return nil, errors.Errorf("test error")
	}

	return t.item, nil
}

func (t *testSingleItemReaderWriter) UpdateItem(item *model.Item) error {
	if t.errorSet {
		return errors.Errorf("test error")
	}

	t.item = item
	return nil
}

func newTestStorageUploader(fileForWritting string, storageUrl string, tempfile string, errorWriting bool) testStorageUploader {
	return testStorageUploader{
		fileForWritting: fileForWritting,
		storageUrl:      storageUrl,
		tempfile:        tempfile,
		errorWriting:    errorWriting,
	}
}

type testStorageUploader struct {
	fileForWritting string
	storageUrl      string
	tempfile        string
	errorWriting    bool
}

func (t *testStorageUploader) GetFileForWriting(name string) (string, error) {
	if t.errorWriting {
		return "", errors.Errorf("test error")
	}

	return t.fileForWritting, nil
}

func (t *testStorageUploader) GetStorageUrl(name string) string {
	return t.storageUrl
}

func (t *testStorageUploader) GetTempFile() string {
	return t.tempfile
}
