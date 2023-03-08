package itemprocessor

import (
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"os"
)

func (p itemProcessorImpl) EnqueueAllItemsPreview(force bool) error {
	items, err := p.gallery.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *items {
		if !force && item.PreviewUrl != "" {
			continue
		}

		p.EnqueueItemPreview(item.Id)
	}

	return nil
}

func (p itemProcessorImpl) EnqueueItemPreview(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_PREVIEW_TASK, IdParam: id})
}

func (p itemProcessorImpl) refreshItemPreview(id uint64) error {
	item, err := p.gallery.GetItem(id)
	if err != nil {
		return err
	}

	logger.Infof("Setting preview for item %d [videoFile: %s] [count: %d] [duration: %d]",
		item.Id, item.Url, p.gallery.PreviewSceneCount, p.gallery.PreviewSceneDuration)

	videoParts, err := p.getPreviewParts(item)
	defer func() {
		for _, file := range videoParts {
			os.Remove(file)
		}
	}()

	if err != nil {
		return err
	}

	relativeFile := fmt.Sprintf("previews/%d/preview.mp4", item.Id)
	storageFile, err := p.storage.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err := ffmpeg.JoinVideoFiles(videoParts, storageFile); err != nil {
		logger.Errorf("Error joining video files for item %d, error %v", item.Id, err)
		return err
	}

	tempFile := fmt.Sprintf("%s.mp4", p.storage.GetTempFile())
	if err := ffmpeg.OptimizeVideoForPreview(storageFile, tempFile); err != nil {
		logger.Errorf("Error optimizing video file for item %d, error %v", item.Id, err)
		return err
	}

	item.PreviewUrl = p.storage.GetStorageUrl(relativeFile)
	return p.gallery.UpdateItem(item)
}

func (p itemProcessorImpl) getPreviewParts(item *model.Item) ([]string, error) {
	videoFile := p.relativasor.GetAbsoluteFile(item.Url)
	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for i := 1; i <= int(p.gallery.PreviewSceneCount); i++ {
		startSecond := (int(duration) / (p.gallery.PreviewSceneCount + 1)) * i
		tempFile := fmt.Sprintf("%s.mp4", p.storage.GetTempFile())
		result = append(result, tempFile)

		if err := ffmpeg.ExtractPartOfVideo(videoFile, startSecond, p.gallery.PreviewSceneDuration, tempFile); err != nil {
			logger.Errorf("Error extracting part of video for item %d, error %v", item.Id, err)
			return nil, err
		}
	}

	return result, nil
}
