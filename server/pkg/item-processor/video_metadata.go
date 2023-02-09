package itemprocessor

import (
	"my-collection/server/pkg/ffmpeg"
)

func (p itemProcessorImpl) EnqueueAllItemsVideoMetadata() error {
	items, err := p.gallery.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *items {
		if item.DurationSeconds != 0 {
			continue
		}

		p.EnqueueItemVideoMetadata(item.Id)
	}

	return nil
}

func (p itemProcessorImpl) EnqueueItemVideoMetadata(id uint64) {
	p.queue <- task{taskType: REFRESH_METADATA_TASK, id: id}
}

func (p itemProcessorImpl) refreshItemMetadata(id uint64) error {
	item, err := p.gallery.GetItem(id)
	if err != nil {
		return err
	}

	if item.DurationSeconds != 0 {
		return nil
	}

	videoFile := p.gallery.GetFile(item.Url)
	logger.Infof("Refreshing video metadata for item %d  [videoFile: %s]", item.Id, videoFile)

	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		logger.Errorf("Error getting duration of a video %s", videoFile)
		return err
	}

	rawVideoMetadata, err := ffmpeg.GetVideoMetadata(videoFile)
	if err != nil {
		logger.Errorf("Error getting video metadata of %s", videoFile)
		return err
	}

	item.DurationSeconds = duration
	item.Width = rawVideoMetadata.Width
	item.Height = rawVideoMetadata.Height
	item.CodecName = rawVideoMetadata.CodecName
	return p.gallery.UpdateItem(item)
}
