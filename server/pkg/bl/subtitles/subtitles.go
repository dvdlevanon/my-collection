package subtitles

import (
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/srt"
	"my-collection/server/pkg/utils"
	"os"
	"path/filepath"
)

var ErrSubtitileNotFound = fmt.Errorf("subtitle not found")

func GetSubtitle(ir model.ItemReader, itemId uint64) (model.Subtitle, error) {
	item, err := ir.GetItem(itemId)
	if err != nil {
		return model.Subtitle{}, err
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	videoDir := filepath.Dir(videoFile)
	subtitlePath := lookForSubtitles(videoDir, utils.BaseRemoveExtension(videoFile))

	if subtitlePath == "" {
		return model.Subtitle{}, ErrSubtitileNotFound
	}

	return srt.LoadFile(subtitlePath)
}

func lookForSubtitles(dir string, movieName string) string {
	specificSubtitle := filepath.Join(dir, fmt.Sprintf("%s.srt", movieName))
	if _, err := os.Stat(specificSubtitle); err == nil {
		return specificSubtitle
	}

	srtFiles, err := filepath.Glob(filepath.Join(dir, "*.srt"))
	if err != nil || len(srtFiles) == 0 {
		return ""
	}

	return srtFiles[0]
}
