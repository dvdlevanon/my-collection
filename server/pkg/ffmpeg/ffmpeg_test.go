package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDuration(t *testing.T) {
	duration, err := GetDurationInSeconds("test.mp4")
	assert.NoError(t, err)
	assert.Equal(t, duration, 5)
}

func TestGetDurationOfMissingFile(t *testing.T) {
	_, err := GetDurationInSeconds("missing.mp4")
	assert.Error(t, err)
}

func TestTakeScreenshot(t *testing.T) {
	err := TakeScreenshot("test.mp4", 3, "test-screenshot.png")
	assert.NoError(t, err)
}
