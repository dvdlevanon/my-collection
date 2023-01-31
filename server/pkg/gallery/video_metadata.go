package gallery

import (
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"time"
)

func (g *Gallery) RefreshItemsVideoMetadata() error {
	items, err := g.GetAllItems()
	if err != nil {
		return err
	}

	startMillis := time.Now().UnixMilli()
	logger.Infof("Start refreshing video metadata of %d items", len(*items))
	errorsCounter := 0

	for _, item := range *items {
		if item.DurationSeconds != 0 {
			continue
		}

		if err := g.refreshItemMetadata(&item); err != nil {
			errorsCounter++
		}
	}

	logger.Infof("Done refreshing video metadata of %d items in %dms - %d errors", len(*items), time.Now().UnixMilli()-startMillis, errorsCounter)
	return nil
}

func (g *Gallery) refreshItemMetadata(item *model.Item) error {
	videoFile := g.GetFile(item.Url)
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
	return g.UpdateItem(item)
}
