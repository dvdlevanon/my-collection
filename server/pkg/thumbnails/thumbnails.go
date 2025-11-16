package thumbnails

import (
	"context"
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"my-collection/server/pkg/utils"
	"path/filepath"
	"time"
)

func New(tr model.TagReader, tiw model.TagImageWriter, storage *storage.Storage, thumbnailWidth int, thumbnailHeight int) *Thumbnails {
	return &Thumbnails{
		tr:              tr,
		tiw:             tiw,
		storage:         storage,
		thumbnailWidth:  thumbnailWidth,
		thumbnailHeight: thumbnailHeight,
	}
}

type Thumbnails struct {
	tr              model.TagReader
	tiw             model.TagImageWriter
	storage         *storage.Storage
	thumbnailWidth  int
	thumbnailHeight int
}

func (t *Thumbnails) Run(ctx context.Context) error {
	ctx = utils.ContextWithSubject(ctx, "thumnails")
	for {
		select {
		case <-time.After(1 * time.Minute):
			if err := t.processThumbnails(ctx); err != nil {
				utils.LogError("Error processing thumbnails", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (t *Thumbnails) processThumbnails(ctx context.Context) error {
	tags, err := t.tr.GetAllTags(ctx)
	if err != nil {
		return err
	}

	for _, tag := range *tags {
		for _, image := range tag.Images {
			if image.ThumbnailUrlRect == image.ThumbnailRect {
				continue
			}

			if err := t.ProcessThumbnail(ctx, image); err != nil {
				utils.LogError("Error processing thumbnail", err)
			}
		}
	}

	return nil
}

func (t *Thumbnails) ProcessThumbnail(ctx context.Context, image *model.TagImage) error {
	relativeFile := filepath.Join("thumbnails", fmt.Sprint(image.TagId), fmt.Sprintf("%d.png", image.Id))
	storageFile, err := t.storage.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	if err := t.extractThumbnail(image, storageFile); err != nil {
		return err
	}

	image.ThumbnailUrl = t.storage.GetStorageUrl(relativeFile)
	image.ThumbnailUrlRect = image.ThumbnailRect
	image.ThumbnailUrlNonce = time.Now().UnixNano()
	return t.tiw.UpdateTagImage(ctx, image)
}

func (t *Thumbnails) extractThumbnail(image *model.TagImage, outputFile string) error {
	imageFile := t.storage.GetFile(image.Url)

	return utils.ExtractImage(imageFile, outputFile, image.ThumbnailRect.X, image.ThumbnailRect.Y, image.ThumbnailRect.H, image.ThumbnailRect.W,
		t.thumbnailWidth, t.thumbnailHeight)
}
