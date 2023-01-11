package gallery

import (
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"path/filepath"
	"strings"
	"time"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("gallery")

type Gallery struct {
	*db.Database
	storage       *storage.Storage
	rootDirectory string
}

func New(db *db.Database, storage *storage.Storage, rootDirectory string) *Gallery {
	return &Gallery{
		Database:      db,
		storage:       storage,
		rootDirectory: rootDirectory,
	}
}

func (g *Gallery) CreateItem(item *model.Item) error {
	item.Url = g.getRelativePath(item.Url)
	return g.Database.CreateItem(item)
}

func (g *Gallery) CreateOrUpdateItem(item *model.Item) error {
	item.Url = g.getRelativePath(item.Url)
	return g.Database.CreateOrUpdateItem(item)
}

func (g *Gallery) getRelativePath(url string) string {
	if !strings.HasPrefix(url, g.rootDirectory) {
		return url
	}
	return strings.TrimPrefix(strings.TrimPrefix(url, g.rootDirectory), string(filepath.Separator))
}

func (g *Gallery) GetItemAbsolutePath(url string) string {
	if strings.HasPrefix(url, string(filepath.Separator)) {
		return url
	}
	return filepath.Join(g.rootDirectory, url)
}

func (g *Gallery) RefreshItemsPreview() error {
	items, err := g.GetAllItems()

	if err != nil {
		return err
	}

	startMillis := time.Now().UnixMilli()
	logger.Infof("Start refreshing preview of %d items", len(*items))
	errorsCounter := 0

	for _, item := range *items {
		if len(item.Previews) > 0 {
			continue
		}

		if err := g.refreshItemsPreview(&item); err != nil {
			errorsCounter++
		}
	}

	logger.Infof("Done refreshing preview of %d items in %dms - %d errors", len(*items), time.Now().UnixMilli()-startMillis, errorsCounter)

	return nil
}

func (g *Gallery) refreshItemsPreview(item *model.Item) error {
	videoFile := g.GetItemAbsolutePath(item.Url)
	logger.Infof("Setting preview for item %d - file %s", item.Id, videoFile)

	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		return err
	}

	outputFile := g.storage.GetFile(fmt.Sprintf("%d.png", item.Id))
	if err := ffmpeg.TakeScreenshot(videoFile, uint64(duration/2), outputFile); err != nil {
		logger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.Previews = append(item.Previews, model.Preview{
		Url: filepath.Base(outputFile),
	})

	return g.UpdateItem(item)
}
