package utils

import (
	"math"
	"os"
	"strings"

	"github.com/h2non/filetype"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("utils")

func IsVideo(trustFileExtenssion bool, path string) bool {
	if trustFileExtenssion {
		path = strings.ToLower(path)
		return strings.HasSuffix(path, ".avi") ||
			strings.HasSuffix(path, ".mkv") ||
			strings.HasSuffix(path, ".mpg") ||
			strings.HasSuffix(path, ".mpeg") ||
			strings.HasSuffix(path, ".wmv") ||
			strings.HasSuffix(path, ".mp4")
	}

	file, err := os.Open(path)
	if err != nil {
		logger.Errorf("Error opening file for reading %s - %t", file, err)
		return false
	}

	stat, err := file.Stat()
	if err != nil {
		logger.Errorf("Error getting stats of file %s - %t", path, err)
		return false
	}

	header := make([]byte, int(math.Max(float64(stat.Size())-1, 1024)))
	_, err = file.Read(header)
	if err != nil {
		logger.Errorf("Error reading from file %s - %t", path, err)
		return false
	}

	return filetype.IsVideo(header)
}
