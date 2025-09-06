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

func getSubtitleByName(dir string, subtitleName string) (string, error) {
	path := filepath.Join(dir, subtitleName)
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func lookForSubtitles(dir string, movieName string) (string, error) {
	specificSubtitle := filepath.Join(dir, fmt.Sprintf("%s.srt", movieName))
	if _, err := os.Stat(specificSubtitle); err == nil {
		return specificSubtitle, nil
	}

	srtFiles, err := filepath.Glob(filepath.Join(dir, "*.srt"))
	if err != nil || len(srtFiles) == 0 {
		return "", err
	}

	return srtFiles[0], nil
}

func lookForAvailableSubtitles(dir string) ([]string, error) {
	names := make([]string, 0)
	srtFiles, err := filepath.Glob(filepath.Join(dir, "*.srt"))
	if err != nil {
		return nil, err
	}

	for _, srt := range srtFiles {
		names = append(names, filepath.Base(srt))
	}

	return names, nil
}

func GetSubtitle(ir model.ItemReader, itemId uint64, subtitleName string) (model.Subtitle, error) {
	item, err := ir.GetItem(itemId)
	if err != nil {
		return model.Subtitle{}, err
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	videoDir := filepath.Dir(videoFile)

	var subtitlePath string
	if subtitleName == "" {
		subtitlePath, err = lookForSubtitles(videoDir, utils.BaseRemoveExtension(videoFile))
	} else {
		subtitlePath, err = getSubtitleByName(videoDir, subtitleName)
	}
	if err != nil {
		return model.Subtitle{}, err
	}
	if subtitlePath == "" {
		return model.Subtitle{}, ErrSubtitileNotFound
	}

	return srt.LoadFile(subtitlePath)
}

func GetAvailableNames(ir model.ItemReader, itemId uint64) ([]string, error) {
	item, err := ir.GetItem(itemId)
	if err != nil {
		return nil, err
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	videoDir := filepath.Dir(videoFile)
	return lookForAvailableSubtitles(videoDir)
}
