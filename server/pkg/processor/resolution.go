package processor

import (
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
)

func changeResolution(irw model.ItemReaderWriter, tempProvider model.TempFileProvider, id uint64, newResolution string) error {
	item, err := irw.GetItem(id)
	if err != nil {
		return err
	}

	resolution, err := ffmpeg.ResolutionFromString(newResolution)
	if err != nil {
		return err
	}

	widthSame := resolution.Width == -1 || item.Width == resolution.Width
	heightSame := resolution.Height == -1 || item.Height == resolution.Height

	if widthSame && heightSame {
		// Probably double enqueue
		return nil
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	tempFile := fmt.Sprintf("%s.mp4", tempProvider.GetTempFile())
	if err := ffmpeg.ChangeVideoResolution(videoFile, tempFile, newResolution); err != nil {
		return err
	}

	return refreshItemMetadata(irw, id)
}
