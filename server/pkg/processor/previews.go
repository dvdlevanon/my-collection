package processor

import (
	"context"
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"
)

func refreshItemPreview(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader,
	previewSceneCount int, previewSceneDuration int, id uint64) error {
	item, err := irw.GetItem(ctx, id)
	if err != nil {
		return err
	}

	if item.PreviewMode == items.PREVIEW_FROM_START_POSITION {
		return nil
	}

	logger.Infof("Setting preview for item %d [videoFile: %s] [count: %d] [duration: %d]",
		item.Id, item.Url, previewSceneCount, previewSceneDuration)

	videoParts, err := getPreviewParts(uploader, item, previewSceneCount, previewSceneDuration)
	defer func() {
		for _, file := range videoParts {
			os.Remove(file)
		}
	}()

	if err != nil {
		return err
	}

	relativeFile := fmt.Sprintf("previews/%d/preview.mp4", item.Id)
	storageFile, err := uploader.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err := ffmpeg.JoinVideoFiles(videoParts, storageFile); err != nil {
		logger.Errorf("Error joining video files for item %d, error %v", item.Id, err)
		return err
	}

	tempFile := fmt.Sprintf("%s.mp4", uploader.GetTempFile())
	if err := ffmpeg.OptimizeVideoForPreview(storageFile, tempFile); err != nil {
		logger.Errorf("Error optimizing video file for item %d, error %v", item.Id, err)
		return err
	}

	item.PreviewUrl = uploader.GetStorageUrl(relativeFile)
	return irw.UpdateItem(ctx, item)
}

func getPreviewParts(uploader model.StorageUploader, item *model.Item,
	previewSceneCount int, previewSceneDuration int) ([]string, error) {
	videoFile := relativasor.GetAbsoluteFile(item.Url)
	duration, err := getDurationForItem(item, videoFile)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for i := 1; i <= int(previewSceneCount); i++ {
		startOffset := 0.0
		if items.IsSubItem(item) || items.IsHighlight(item) {
			startOffset = item.StartPosition
		}

		startSecond := startOffset + ((duration / float64(previewSceneCount+1)) * float64(i))
		tempFile := fmt.Sprintf("%s.mp4", uploader.GetTempFile())
		result = append(result, tempFile)

		if err := ffmpeg.ExtractPartOfVideo(videoFile, startSecond, previewSceneDuration, tempFile); err != nil {
			logger.Errorf("Error extracting part of video for item %d, error %v", item.Id, err)
			return nil, err
		}
	}

	return result, nil
}
