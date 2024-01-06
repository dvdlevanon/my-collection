package processor

import (
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"time"

	"k8s.io/utils/pointer"
)

func refreshItemCovers(irw model.ItemReaderWriter, uploader model.StorageUploader, id uint64, coversCount int) error {
	item, err := irw.GetItem(id)
	if err != nil {
		return err
	}

	item.Covers = make([]model.Cover, 0)

	for i := 1; i <= int(coversCount); i++ {
		if err := refreshItemCover(irw, uploader, item, coversCount, i); err != nil {
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
		logger.Errorf("Error getting duration of a video %s", videoFile)
		return 0, err
	}

	return duration, err
}

func refreshItemCover(irw model.ItemReaderWriter, uploader model.StorageUploader,
	item *model.Item, coversCount int, coverNumber int) error {
	videoFile := relativasor.GetAbsoluteFile(item.Url)
	logger.Infof("Setting cover for item %d [coverNumber: %d] [videoFile: %s]", item.Id, coverNumber, videoFile)

	duration, err := getDurationForItem(item, videoFile)
	if err != nil {
		return err
	}

	relativeFile := fmt.Sprintf("covers/%d/%d.png", item.Id, coverNumber)
	storageFile, err := uploader.GetFileForWriting(relativeFile)
	if err != nil {
		logger.Errorf("Error getting new cover file from storage %v", err)
		return err
	}

	startOffset := 0.0
	if items.IsSubItem(item) || items.IsHighlight(item) {
		startOffset = item.StartPosition
	}

	screenshotSecond := startOffset + ((duration / float64(coversCount+1)) * float64(coverNumber))
	if err := ffmpeg.TakeScreenshot(videoFile, screenshotSecond, storageFile); err != nil {
		logger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.Covers = append(item.Covers, model.Cover{
		Url: uploader.GetStorageUrl(relativeFile),
	})

	return irw.UpdateItem(item)
}

func refreshMainCover(irw model.ItemReaderWriter, uploader model.StorageUploader, id uint64, second float64) error {
	item, err := irw.GetItem(id)
	if err != nil {
		return err
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	logger.Infof("Setting main cover for item %d [second: %s]", item.Id, second)

	relativeFile := fmt.Sprintf("main-covers/%d/main.png", item.Id)
	storageFile, err := uploader.GetFileForWriting(relativeFile)
	if err != nil {
		logger.Errorf("Error getting main cover file from storage %v", err)
		return err
	}

	if err := ffmpeg.TakeScreenshot(videoFile, second, storageFile); err != nil {
		logger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.MainCoverSecond = second
	item.MainCoverUrl = pointer.String(uploader.GetStorageUrl(relativeFile))
	item.MainCoverNonce = time.Now().UnixNano()
	return irw.UpdateItem(item)
}
