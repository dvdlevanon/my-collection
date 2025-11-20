package video_tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"

	"github.com/op/go-logging"
)

var pLogger = logging.MustGetLogger("video-preview")

type videoPreviewParams struct {
	ItemId        uint64 `json:"id"`
	SceneCount    int    `json:"count"`
	SceneDuration int    `json:"duration"`
}

func PreviewDesc(id uint64, title string, sceneCount int, sceneDuration int) string {
	return fmt.Sprintf("Generate preview with %d scenes (%ds each) for %s", sceneCount, sceneDuration, title)
}

func MarshalVideoPreviewParams(id uint64, sceneCount int, sceneDuration int) (string, error) {
	p := videoPreviewParams{ItemId: id, SceneCount: sceneCount, SceneDuration: sceneDuration}
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func unmarshalVideoPreviewParams(params string) (videoPreviewParams, error) {
	var p videoPreviewParams
	if err := json.Unmarshal([]byte(params), &p); err != nil {
		return p, err
	}
	return p, nil
}

func RefreshVideoPreview(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader, params string) error {
	p, err := unmarshalVideoPreviewParams(params)
	if err != nil {
		return err
	}

	return refreshVideoPreview(ctx, irw, uploader, p)
}

func refreshVideoPreview(ctx context.Context, irw model.ItemReaderWriter, uploader model.StorageUploader, p videoPreviewParams) error {
	item, err := irw.GetItem(ctx, p.ItemId)
	if err != nil {
		return err
	}

	if item.PreviewMode == items.PREVIEW_FROM_START_POSITION {
		return nil
	}

	if p.SceneCount == 0 {
		return nil
	}

	pLogger.Infof("Setting preview for item %d [videoFile: %s] [count: %d] [duration: %d]",
		item.Id, item.Url, p.SceneCount, p.SceneDuration)

	videoParts, err := getPreviewParts(uploader, item, p.SceneCount, p.SceneDuration)
	defer func() {
		for _, file := range videoParts {
			os.Remove(file)
		}
	}()

	if err != nil {
		return err
	}

	relativeFile := fmt.Sprintf("previews/%d/preview.mp4", item.Id)
	storageFile, err := uploader.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err := ffmpeg.JoinVideoFiles(videoParts, storageFile); err != nil {
		pLogger.Errorf("Error joining video files for item %d, error %v", item.Id, err)
		return err
	}

	tempFile := fmt.Sprintf("%s.mp4", uploader.GetTempFile())
	if err := ffmpeg.OptimizeVideoForPreview(storageFile, tempFile); err != nil {
		pLogger.Errorf("Error optimizing video file for item %d, error %v", item.Id, err)
		return err
	}

	item.PreviewUrl = uploader.GetStorageUrl(relativeFile)
	return irw.UpdateItem(ctx, item)
}

func getPreviewParts(uploader model.StorageUploader, item *model.Item,
	previewSceneCount int, previewSceneDuration int) ([]string, error) {
	videoFile := relativasor.GetAbsoluteFile(item.Url)
	duration, err := getDurationForItem(item, videoFile)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for i := 1; i <= int(previewSceneCount); i++ {
		startOffset := 0.0
		if items.IsSubItem(item) || items.IsHighlight(item) {
			startOffset = item.StartPosition
		}

		startSecond := startOffset + ((duration / float64(previewSceneCount+1)) * float64(i))
		tempFile := fmt.Sprintf("%s.mp4", uploader.GetTempFile())
		result = append(result, tempFile)

		if err := ffmpeg.ExtractPartOfVideo(videoFile, startSecond, previewSceneDuration, tempFile); err != nil {
			pLogger.Errorf("Error extracting part of video for item %d, error %v", item.Id, err)
			return nil, err
		}
	}

	return result, nil
}
