package gallery

import (
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("gallery")

type Gallery struct {
	*db.Database
	storage              *storage.Storage
	rootDirectory        string
	coversCount          int
	previewSceneCount    int
	previewSceneDuration int
}

func New(db *db.Database, storage *storage.Storage, rootDirectory string) *Gallery {
	return &Gallery{
		Database:             db,
		storage:              storage,
		rootDirectory:        rootDirectory,
		coversCount:          3,
		previewSceneCount:    4,
		previewSceneDuration: 3,
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

func (g *Gallery) GetFile(url string) string {
	if strings.HasPrefix(url, string(filepath.Separator)) {
		return url
	} else {
		return filepath.Join(g.rootDirectory, url)
	}
}
