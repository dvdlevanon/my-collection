package itemprocessor

import (
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
)

func (p itemProcessorImpl) EnqueueAllItemsVideoMetadata(force bool) error {
	items, err := p.db.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *items {
		if !force && item.DurationSeconds != 0 {
			continue
		}

		p.EnqueueItemVideoMetadata(item.Id)
	}

	return nil
}

func (p itemProcessorImpl) EnqueueItemVideoMetadata(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_METADATA_TASK, IdParam: id})
}

func (p itemProcessorImpl) refreshItemMetadata(id uint64) error {
	item, err := p.db.GetItem(id)
	if err != nil {
		return err
	}

	videoFile := p.relativasor.GetAbsoluteFile(item.Url)
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

	rawAudioMetadata, err := ffmpeg.GetAudioMetadata(videoFile)
	if err != nil {
		logger.Errorf("Error getting audio metadata of %s", videoFile)
		return err
	}

	item.DurationSeconds = duration
	item.Width = rawVideoMetadata.Width
	item.Height = rawVideoMetadata.Height
	item.VideoCodecName = rawVideoMetadata.CodecName
	item.AudioCodecName = rawAudioMetadata.CodecName
	return p.db.UpdateItem(item)
}
