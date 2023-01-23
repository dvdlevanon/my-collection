package gallery

import (
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"time"
)

func (g *Gallery) RefreshItemsCovers() error {
	items, err := g.GetAllItems()
	if err != nil {
		return err
	}

	startMillis := time.Now().UnixMilli()
	logger.Infof("Start refreshing covers of %d items", len(*items))
	errorsCounter := 0

	for _, item := range *items {
		if len(item.Covers) == g.coversCount {
			continue
		}

		item.Covers = make([]model.Cover, 0)

		for i := 1; i <= int(g.coversCount); i++ {
			if err := g.refreshItemCover(&item, i); err != nil {
				errorsCounter++
			}
		}
	}

	logger.Infof("Done refreshing covers of %d items in %dms - %d errors", len(*items), time.Now().UnixMilli()-startMillis, errorsCounter)
	return nil
}

func (g *Gallery) refreshItemCover(item *model.Item, coverNumber int) error {
	videoFile := g.GetFile(item.Url)
	logger.Infof("Setting cover for item %d [coverNumber: %d] [videoFile: %s]", item.Id, coverNumber, videoFile)

	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		logger.Errorf("Error getting duration of a video %s", videoFile)
		return err
	}

	relativeFile := fmt.Sprintf("covers/%d/%d.png", item.Id, coverNumber)
	storageFile, err := g.storage.GetFileForWriting(relativeFile)
	if err != nil {
		logger.Errorf("Error getting new cover file from storage %v", err)
		return err
	}

	screenshotSecond := (int(duration) / (g.coversCount + 1)) * coverNumber
	if err := ffmpeg.TakeScreenshot(videoFile, screenshotSecond, storageFile); err != nil {
		logger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.Covers = append(item.Covers, model.Cover{
		Url: g.storage.GetStorageUrl(relativeFile),
	})

	return g.UpdateItem(item)
}
