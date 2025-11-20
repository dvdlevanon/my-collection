package video_tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"

	"github.com/op/go-logging"
)

var vmLogger = logging.MustGetLogger("file-metadata")

type videoMetadataParams struct {
	ItemId uint64 `json:"id,omitempty"`
}

func MetadataDesc(id uint64, title string) string {
	return fmt.Sprintf("Update metadata for %s", title)
}

func MarshalVideoMetadataParams(id uint64) (string, error) {
	p := videoMetadataParams{ItemId: id}
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func unmarshalVideoMetadataParams(params string) (videoMetadataParams, error) {
	var p videoMetadataParams
	if err := json.Unmarshal([]byte(params), &p); err != nil {
		return p, err
	}
	return p, nil
}

func UpdateVideoMetadata(ctx context.Context, irw model.ItemReaderWriter, params string) error {
	p, err := unmarshalVideoMetadataParams(params)
	if err != nil {
		return err
	}

	return updateVideoMetadata(ctx, irw, p)
}

func updateVideoMetadata(ctx context.Context, irw model.ItemReaderWriter, p videoMetadataParams) error {
	item, err := irw.GetItem(ctx, p.ItemId)
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
	vmLogger.Infof("Refreshing video metadata for item %d  [videoFile: %s]", item.Id, videoFile)

	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		vmLogger.Errorf("Error getting duration of a video %s", videoFile)
		return err
	}

	rawVideoMetadata, err := ffmpeg.GetVideoMetadata(videoFile)
	if err != nil {
		vmLogger.Errorf("Error getting video metadata of %s", videoFile)
		return err
	}

	rawAudioMetadata, err := ffmpeg.GetAudioMetadata(videoFile)
	if err != nil {
		vmLogger.Errorf("Error getting audio metadata of %s", videoFile)
		return err
	}

	item.DurationSeconds = duration
	item.Width = rawVideoMetadata.Width
	item.Height = rawVideoMetadata.Height
	item.VideoCodecName = rawVideoMetadata.CodecName
	item.AudioCodecName = rawAudioMetadata.CodecName
	return nil
}
