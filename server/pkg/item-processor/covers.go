package itemprocessor

import (
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"

	"k8s.io/utils/pointer"
)

func (p itemProcessorImpl) EnqueueAllItemsCovers() error {
	items, err := p.gallery.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *items {
		if len(item.Covers) == p.gallery.CoversCount {
			continue
		}

		p.EnqueueItemCovers(item.Id)
	}

	return nil
}

func (p itemProcessorImpl) EnqueueMainCover(id uint64, second float64) {
	p.queue <- task{taskType: SET_MAIN_COVER, id: id, floatParam: second}
}

func (p itemProcessorImpl) EnqueueItemCovers(id uint64) {
	p.queue <- task{taskType: REFRESH_COVER_TASK, id: id}
}

func (p itemProcessorImpl) refreshItemCovers(id uint64) error {
	item, err := p.gallery.GetItem(id)
	if err != nil {
		return err
	}

	if len(item.Covers) == p.gallery.CoversCount {
		return nil
	}

	item.Covers = make([]model.Cover, 0)

	for i := 1; i <= int(p.gallery.CoversCount); i++ {
		if err := p.refreshItemCover(item, i); err != nil {
			return err
		}
	}

	return nil
}

func (p itemProcessorImpl) refreshItemCover(item *model.Item, coverNumber int) error {
	videoFile := p.gallery.GetFile(item.Url)
	logger.Infof("Setting cover for item %d [coverNumber: %d] [videoFile: %s]", item.Id, coverNumber, videoFile)

	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		logger.Errorf("Error getting duration of a video %s", videoFile)
		return err
	}

	relativeFile := fmt.Sprintf("covers/%d/%d.png", item.Id, coverNumber)
	storageFile, err := p.storage.GetFileForWriting(relativeFile)
	if err != nil {
		logger.Errorf("Error getting new cover file from storage %v", err)
		return err
	}

	screenshotSecond := (int(duration) / (p.gallery.CoversCount + 1)) * coverNumber
	if err := ffmpeg.TakeScreenshot(videoFile, float64(screenshotSecond), storageFile); err != nil {
		logger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.Covers = append(item.Covers, model.Cover{
		Url: p.storage.GetStorageUrl(relativeFile),
	})

	return p.gallery.UpdateItem(item)
}

func (p itemProcessorImpl) setMainCover(id uint64, second float64) error {
	item, err := p.gallery.GetItem(id)
	if err != nil {
		return err
	}

	videoFile := p.gallery.GetFile(item.Url)
	logger.Infof("Setting main cover for item %d [second: %s]", item.Id, second)

	relativeFile := fmt.Sprintf("main-covers/%d/main.png", item.Id)
	storageFile, err := p.storage.GetFileForWriting(relativeFile)
	if err != nil {
		logger.Errorf("Error getting main cover file from storage %v", err)
		return err
	}

	if err := ffmpeg.TakeScreenshot(videoFile, second, storageFile); err != nil {
		logger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.MainCoverUrl = pointer.String(p.storage.GetStorageUrl(relativeFile))
	return p.gallery.UpdateItem(item)
}
