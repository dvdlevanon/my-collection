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
	coversCount   int
}

func New(db *db.Database, storage *storage.Storage, rootDirectory string) *Gallery {
	return &Gallery{
		Database:      db,
		storage:       storage,
		rootDirectory: rootDirectory,
		coversCount:   1,
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

func (g *Gallery) RefreshItemsCovers() error {
	items, err := g.GetAllItems()

	if err != nil {
		return err
	}

	startMillis := time.Now().UnixMilli()
	logger.Infof("Start refreshing covers of %d items", len(*items))
	errorsCounter := 0

	for _, item := range *items {
		if len(item.Covers) == g.coversCount {
			continue
		}

		item.Covers = make([]model.Cover, 0)

		for i := 1; i <= int(g.coversCount); i++ {
			if err := g.refreshItemCover(&item, i); err != nil {
				errorsCounter++
			}
		}
	}

	logger.Infof("Done refreshing covers of %d items in %dms - %d errors", len(*items), time.Now().UnixMilli()-startMillis, errorsCounter)

	return nil
}

func (g *Gallery) refreshItemCover(item *model.Item, coverNumber int) error {
	videoFile := g.GetItemAbsolutePath(item.Url)
	logger.Infof("Setting cover for item %d [coverNumber: %d] [videoFile: %s]", item.Id, videoFile)

	duration, err := ffmpeg.GetDurationInSeconds(videoFile)
	if err != nil {
		return err
	}

	relativeFile := fmt.Sprintf("covers/%d/%d.png", item.Id, coverNumber)
	storageFile, err := g.storage.GetFileForWriting(relativeFile)
	if err != nil {
		return err
	}

	screenshotSecond := (int(duration) / (g.coversCount + 1)) * coverNumber
	if err := ffmpeg.TakeScreenshot(videoFile, screenshotSecond, storageFile); err != nil {
		logger.Errorf("Error taking screenshot for item %d, error %v", item.Id, err)
		return err
	}

	item.Covers = append(item.Covers, model.Cover{
		Url: relativeFile,
	})

	return g.UpdateItem(item)
}
