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

func (g *Gallery) GetTagAvailableAnnotations(tagId uint64) ([]model.TagAnnotation, error) {
	tag, err := g.GetTag(tagId)
	if err != nil {
		return nil, err
	}

	availableAnnotations := make(map[uint64]model.TagAnnotation)
	for _, child := range tag.Children {
		annotations, err := g.GetTagAnnotations(child.Id)
		if err != nil {
			return nil, err
		}

		for _, annotation := range annotations {
			availableAnnotations[annotation.Id] = annotation
		}
	}

	result := make([]model.TagAnnotation, 0, len(availableAnnotations))
	for _, v := range availableAnnotations {
		result = append(result, v)
	}

	return result, nil
}
