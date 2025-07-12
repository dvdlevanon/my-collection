package processor

import (
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"path/filepath"
	"time"
)

func cropFrame(irw model.ItemReaderWriter, uploader model.StorageUploader, id uint64, second float64, rectStr string) error {
	item, err := irw.GetItem(id)
	if err != nil {
		return err
	}

	rect, err := model.DeserializeRectFloat(rectStr)
	if err != nil {
		return fmt.Errorf("unable to deserialize rect %s %w", rectStr, err)
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	logger.Infof("Cropping frame %d [second: %s] [rect: %s]", item.Id, second, rect)

	relativeFile := fmt.Sprintf("frames/%s_%s.png", filepath.Base(item.Url), time.Now().Format("2006-01-02_15_04_05"))
	storageFile, err := uploader.GetFileForWriting(relativeFile)
	if err != nil {
		logger.Errorf("Error getting cropped file from storage %v", err)
		return err
	}

	if err := ffmpeg.CropScreenshot(videoFile, second, rect, storageFile); err != nil {
		logger.Errorf("Error cropping frame for item %d, error %v", item.Id, err)
		return err
	}

	logger.Infof("Frame cropped to %s", relativeFile)
	return nil
}
