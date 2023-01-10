package gallery

import (
	"fmt"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"path/filepath"
	"strings"
)

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
	return g.Database.CreateItem(g.normalizeUrl(item))
}

func (g *Gallery) CreateOrUpdateItem(item *model.Item) error {
	return g.Database.CreateOrUpdateItem(g.normalizeUrl(item))
}

func (g *Gallery) normalizeUrl(item *model.Item) *model.Item {
	item.Url = strings.TrimPrefix(item.Url, g.rootDirectory)
	item.Url = strings.TrimPrefix(item.Url, string(filepath.Separator))
	return item
}

func (g *Gallery) Watch(itemId uint64) error {
	item, err := g.GetItem(itemId)

	if err != nil {
		return err
	}

	fmt.Printf("Watching %s/%s\n", g.rootDirectory, item.Url)

	return nil
}

func (g *Gallery) RefreshItemsPreview() error {
	items, err := g.GetAllItems()

	if err != nil {
		return err
	}

	for _, item := range *items {
		if len(item.Previews) > 0 {
			continue
		}

		fmt.Printf("Setting preview for %v %v\n", item.Id, len(item.Previews))

		videoFile := filepath.Join(g.rootDirectory, item.Url)
		duration, err := ffmpeg.GetDurationInSeconds(videoFile)

		if err != nil {
			return err
		}

		outputFile := g.storage.GetFile(fmt.Sprintf("%d.png", item.Id))

		if err := ffmpeg.TakeScreenshot(videoFile, uint64(duration/2), outputFile); err != nil {
			fmt.Printf("ERROR %v\n", err)
			continue
		}

		item.Previews = append(item.Previews, model.Preview{
			Url: filepath.Base(outputFile),
		})

		if err = g.UpdateItem(&item); err != nil {
			return err
		}
	}

	return nil
}
