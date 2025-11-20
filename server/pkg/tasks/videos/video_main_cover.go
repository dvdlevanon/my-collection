package video_tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"time"

	"k8s.io/utils/ptr"
)

type videoMainCoverParams struct {
	ItemId uint64  `json:"id"`
	Second float64 `json:"second"`
}

func MainCoverDesc(id uint64, title string, second float64) string {
	return fmt.Sprintf("Set main cover at %.2fs for %s", second, title)
}

func MarshalMainCoverParams(id uint64, second float64) (string, error) {
	p := videoMainCoverParams{ItemId: id, Second: second}
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func unmarshalMainCoverParams(params string) (videoMainCoverParams, error) {
	var p videoMainCoverParams
	if err := json.Unmarshal([]byte(params), &p); err != nil {
		return p, err
	}
	return p, nil
}

func UpdateVideoMainCover(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader, params string) error {
	p, err := unmarshalMainCoverParams(params)
	if err != nil {
		return err
	}

	return updateVideoMainCover(ctx, irw, uploader, p)
}

func updateVideoMainCover(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader, p videoMainCoverParams) error {
	item, err := irw.GetItem(ctx, p.ItemId)
	if err != nil {
		return err
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	vcLogger.Infof("Setting main cover for item %d [second: %s]", item.Id, p.Second)

	relativeFile := fmt.Sprintf("main-covers/%d/main.png", item.Id)
	storageFile, err := uploader.GetFileForWriting(relativeFile)
	if err != nil {
		vcLogger.Errorf("Error getting main cover file from storage %v", err)
		return err
	}

	if err := ffmpeg.TakeScreenshot(videoFile, p.Second, storageFile); err != nil {
		vcLogger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.MainCoverSecond = p.Second
	item.MainCoverUrl = ptr.To(uploader.GetStorageUrl(relativeFile))
	item.MainCoverNonce = time.Now().UnixNano()
	return irw.UpdateItem(ctx, item)
}
