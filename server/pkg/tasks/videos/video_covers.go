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

var vcLogger = logging.MustGetLogger("video-covers")

type videoCoversParams struct {
	ItemId uint64 `json:"id"`
	Count  int    `json:"count"`
}

func CoversDesc(id uint64, title string, count int) string {
	return fmt.Sprintf("Extract %d covers for %s", count, title)
}

func MarshalVideoCoversParams(id uint64, count int) (string, error) {
	p := videoCoversParams{ItemId: id, Count: count}
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func unmarshalVideoCoversParams(params string) (videoCoversParams, error) {
	var p videoCoversParams
	if err := json.Unmarshal([]byte(params), &p); err != nil {
		return p, err
	}
	return p, nil
}

func RefreshVideoCovers(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader, params string) error {
	p, err := unmarshalVideoCoversParams(params)
	if err != nil {
		return err
	}

	return refreshVideoCovers(ctx, irw, uploader, p)
}

func refreshVideoCovers(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader, p videoCoversParams) error {
	if p.Count == 0 {
		return nil
	}

	item, err := irw.GetItem(ctx, p.ItemId)
	if err != nil {
		return err
	}

	item.Covers = make([]model.Cover, 0)

	for i := 1; i <= int(p.Count); i++ {
		if err := refreshItemCover(ctx, irw, uploader, item, p.Count, i); err != nil {
			return err
		}
	}

	return nil
}

func getDurationForItem(item *model.Item, videoFile string) (float64, error) {
	if items.IsSubItem(item) || items.IsHighlight(item) {
		return item.DurationSeconds, nil
	}

	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		vcLogger.Errorf("Error getting duration of a video %s", videoFile)
		return 0, err
	}

	return duration, err
}

func refreshItemCover(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader,
	item *model.Item, coversCount int, coverNumber int) error {
	videoFile := relativasor.GetAbsoluteFile(item.Url)
	vcLogger.Infof("Setting cover for item %d [coverNumber: %d] [videoFile: %s]", item.Id, coverNumber, videoFile)

	duration, err := getDurationForItem(item, videoFile)
	if err != nil {
		return err
	}

	relativeFile := fmt.Sprintf("covers/%d/%d.png", item.Id, coverNumber)
	storageFile, err := uploader.GetFileForWriting(relativeFile)
	if err != nil {
		vcLogger.Errorf("Error getting new cover file from storage %v", err)
		return err
	}

	startOffset := 0.0
	if items.IsSubItem(item) || items.IsHighlight(item) {
		startOffset = item.StartPosition
	}

	screenshotSecond := startOffset + ((duration / float64(coversCount+1)) * float64(coverNumber))
	if err := ffmpeg.TakeScreenshot(videoFile, screenshotSecond, storageFile); err != nil {
		vcLogger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.Covers = append(item.Covers, model.Cover{
		Url: uploader.GetStorageUrl(relativeFile),
	})

	return irw.UpdateItem(ctx, item)
}
