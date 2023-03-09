package processor

import (
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
)

func refreshItemMetadata(irw model.ItemReaderWriter, id uint64) error {
	item, err := irw.GetItem(id)
	if err != nil {
		return err
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
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
	return irw.UpdateItem(item)
}
