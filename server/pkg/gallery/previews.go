package gallery

import (
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"os"
	"time"
)

func (g *Gallery) RefreshItemsPreview() error {
	items, err := g.GetAllItems()
	if err != nil {
		return err
	}

	startMillis := time.Now().UnixMilli()
	logger.Infof("Start refreshing preview of %d items", len(*items))
	errorsCounter := 0

	for _, item := range *items {
		if item.PreviewUrl != "" {
			continue
		}

		if err := g.refreshItemPreview(&item); err != nil {
			errorsCounter++
		}
	}

	logger.Infof("Done refreshing previews of %d items in %dms - %d errors", len(*items), time.Now().UnixMilli()-startMillis, errorsCounter)
	return nil
}

func (g *Gallery) refreshItemPreview(item *model.Item) error {
	logger.Infof("Setting preview for item %d [videoFile: %s] [count: %d] [duration: %d]",
		item.Id, item.Url, g.previewSceneCount, g.previewSceneDuration)

	videoParts, err := g.getPreviewParts(item)
	defer func() {
		for _, file := range videoParts {
			os.Remove(file)
		}
	}()

	if err != nil {
		return err
	}

	relativeFile := fmt.Sprintf("previews/%d/preview.mp4", item.Id)
	storageFile, err := g.storage.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err := ffmpeg.JoinVideoFiles(videoParts, storageFile); err != nil {
		logger.Errorf("Error joining video files for item %d, error %v", item.Id, err)
		return err
	}

	tempFile := fmt.Sprintf("%s.mp4", g.storage.GetTempFile())
	if err := ffmpeg.OptimizeVideoForPreview(storageFile, tempFile); err != nil {
		logger.Errorf("Error optimizing video file for item %d, error %v", item.Id, err)
		return err
	}

	item.PreviewUrl = g.storage.GetStorageUrl(relativeFile)
	return g.UpdateItem(item)
}

func (g *Gallery) getPreviewParts(item *model.Item) ([]string, error) {
	videoFile := g.GetFile(item.Url)
	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for i := 1; i <= int(g.previewSceneCount); i++ {
		startSecond := (int(duration) / (g.previewSceneCount + 1)) * i
		tempFile := fmt.Sprintf("%s.mp4", g.storage.GetTempFile())
		result = append(result, tempFile)

		if err := ffmpeg.ExtractPartOfVideo(videoFile, startSecond, g.previewSceneDuration, tempFile); err != nil {
			logger.Errorf("Error extracting part of video for item %d, error %v", item.Id, err)
			return nil, err
		}
	}

	return result, nil
}
