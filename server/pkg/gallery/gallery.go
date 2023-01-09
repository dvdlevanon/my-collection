package gallery

import (
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"path/filepath"
	"strings"
)

type Gallery struct {
	*db.Database
	rootDirectory string
}

func New(db *db.Database, rootDirectory string) *Gallery {
	return &Gallery{
		Database:      db,
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
