package ffmpeg

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFiles = "../../../testdata/sample-files"
var sampleMp4 = filepath.Join(testFiles, "sample.mp4")
var sample3SecondsScreenshotPng = filepath.Join(testFiles, "sample-3-second-screenshot.png")

func TestGetDuration(t *testing.T) {
	duration, err := GetDurationInSeconds(filepath.Join(testFiles, "sample.mp4"))
	assert.NoError(t, err)
	assert.Equal(t, duration, 5.568)
}

func TestGetDurationOfMissingFile(t *testing.T) {
	_, err := GetDurationInSeconds("missing.mp4")
	assert.Error(t, err)
}

func TestTakeScreenshot(t *testing.T) {
	err := TakeScreenshot(sampleMp4, 3, sample3SecondsScreenshotPng)
	assert.NoError(t, err)
}
