package video_tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
)

type videoResolutionParams struct {
	ItemId uint64 `json:"id,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

func ResolutionDesc(id uint64, title string, w int, h int) string {
	if w == -1 {
		return fmt.Sprintf("Change resolution to height %d for %s", h, title)
	} else if h == -1 {
		return fmt.Sprintf("Change resolution to width %d for %s", w, title)
	}
	return fmt.Sprintf("Change resolution to %dx%d for %s", w, h, title)
}

func MarshalVideoResolutionParams(id uint64, w int, h int) (string, error) {
	p := videoResolutionParams{ItemId: id, Width: w, Height: h}
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func unmarshalVideoResolutionParams(params string) (videoResolutionParams, error) {
	var p videoResolutionParams
	if err := json.Unmarshal([]byte(params), &p); err != nil {
		return p, err
	}
	return p, nil
}

func ChangeVideoResolution(ctx context.Context, irw model.ItemReaderWriter, tempProvider model.TempFileProvider, params string) error {
	p, err := unmarshalVideoResolutionParams(params)
	if err != nil {
		return err
	}

	return changeVideoResolution(ctx, irw, tempProvider, p)
}

func changeVideoResolution(ctx context.Context, irw model.ItemReaderWriter, tempProvider model.TempFileProvider, p videoResolutionParams) error {
	item, err := irw.GetItem(ctx, p.ItemId)
	if err != nil {
		return err
	}

	widthSame := p.Width == -1 || item.Width == p.Width
	heightSame := p.Height == -1 || item.Height == p.Height

	if widthSame && heightSame {
		return nil
	}

	videoFile := relativasor.GetAbsoluteFile(item.Url)
	tempFile := fmt.Sprintf("%s.mp4", tempProvider.GetTempFile())
	if err := ffmpeg.ChangeVideoResolution(videoFile, tempFile, p.Width, p.Height); err != nil {
		return err
	}

	return updateVideoMetadata(ctx, irw, videoMetadataParams{ItemId: p.ItemId})
}
