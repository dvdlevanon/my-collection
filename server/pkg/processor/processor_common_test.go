package processor

import (
	"path/filepath"
)

var testFiles = "../../test-files"
var sampleMp4 = filepath.Join(testFiles, "sample.mp4")
var sampleNoAudioMp4 = filepath.Join(testFiles, "sample-no-audio.mp4")
var sampleNoVideoMp4 = filepath.Join(testFiles, "sample-no-video.mp4")
var sample3SecondsScreenshotPng = filepath.Join(testFiles, "sample-3-second-screenshot.png")
var sample4_5SecondsScreenshotPng = filepath.Join(testFiles, "sample-4_5-second-screenshot.png")
