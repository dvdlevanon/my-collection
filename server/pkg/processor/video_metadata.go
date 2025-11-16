package processor

import (
	"context"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
)

func refreshItemMetadata(ctx context.Context, irw model.ItemReaderWriter, id uint64) error {
	item, err := irw.GetItem(ctx, id)
	if err != nil {
		return err
	}

	if items.IsSubItem(item) || items.IsHighlight(item) {
		if err := updateNonMainItemMetadata(ctx, irw, item); err != nil {
			return err
		}
	} else {
		if err := updateMainItemMetadata(item); err != nil {
			return err
		}
	}

	return irw.UpdateItem(ctx, item)
}

func updateNonMainItemMetadata(ctx context.Context, ir model.ItemReader, item *model.Item) error {
	item.DurationSeconds = item.EndPosition - item.StartPosition

	mainItem, err := ir.GetItem(ctx, item.MainItemId)
	if err != nil {
		return err
	}

	item.Width = mainItem.Width
	item.Height = mainItem.Height
	item.VideoCodecName = mainItem.VideoCodecName
	item.AudioCodecName = mainItem.AudioCodecName
	return nil
}

func updateMainItemMetadata(item *model.Item) error {
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
	return nil
}
