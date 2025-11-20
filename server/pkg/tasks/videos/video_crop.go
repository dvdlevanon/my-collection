package video_tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"path/filepath"
	"time"

	"github.com/op/go-logging"
)

var cropLogger = logging.MustGetLogger("video-crop")

type videoCropParams struct {
	ItemId uint64  `json:"id"`
	Second float64 `json:"second"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	H      float64 `json:"height"`
	W      float64 `json:"width"`
}

func CropDesc(id uint64, title string, second float64, r model.RectFloat) string {
	return fmt.Sprintf("Cropp frame at %.2fs for %s (rect: %.0f,%.0f %.0fx%.0f)", second, title, r.X, r.Y, r.W, r.H)
}

func MarshalVideoCropParams(id uint64, second float64, r model.RectFloat) (string, error) {
	p := videoCropParams{ItemId: id, Second: second, X: r.X, Y: r.Y, W: r.W, H: r.H}
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func unmarshalVideoCropParams(params string) (videoCropParams, error) {
	var p videoCropParams
	if err := json.Unmarshal([]byte(params), &p); err != nil {
		return p, err
	}
	return p, nil
}

func CropVideoFrame(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader, params string) error {
	p, err := unmarshalVideoCropParams(params)
	if err != nil {
		return err
	}

	return cropVideoFrame(ctx, irw, uploader, p)
}

func cropVideoFrame(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader, p videoCropParams) error {
	item, err := irw.GetItem(ctx, p.ItemId)
	if err != nil {
		return err
	}

	rect := model.RectFloat{X: p.X, Y: p.Y, W: p.W, H: p.H}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	cropLogger.Infof("Cropping frame %d [second: %s] [rect: %s]", item.Id, p.Second, rect)

	relativeFile := fmt.Sprintf("frames/%s_%s.png", filepath.Base(item.Url), time.Now().Format("2006-01-02_15_04_05"))
	storageFile, err := uploader.GetFileForWriting(relativeFile)
	if err != nil {
		cropLogger.Errorf("Error getting cropped file from storage %v", err)
		return err
	}

	if err := ffmpeg.CropScreenshot(videoFile, p.Second, rect, storageFile); err != nil {
		cropLogger.Errorf("Error cropping frame for item %d, error %v", item.Id, err)
		return err
	}

	cropLogger.Infof("Frame cropped to %s", relativeFile)
	return nil
}
