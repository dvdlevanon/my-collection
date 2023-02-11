package gallery

import (
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
	"k8s.io/utils/pointer"
)

var logger = logging.MustGetLogger("gallery")

type Gallery struct {
	*db.Database
	storage              *storage.Storage
	rootDirectory        string
	CoversCount          int
	PreviewSceneCount    int
	PreviewSceneDuration int
	AutomaticProcessing  bool
}

func New(db *db.Database, storage *storage.Storage, rootDirectory string) *Gallery {
	return &Gallery{
		Database:             db,
		storage:              storage,
		rootDirectory:        rootDirectory,
		CoversCount:          3,
		PreviewSceneCount:    4,
		PreviewSceneDuration: 3,
		AutomaticProcessing:  false,
	}
}

func (g *Gallery) CreateOrUpdateDirectory(directory *model.Directory) error {
	directory.Path = g.getRelativePath(directory.Path)
	return g.Database.CreateOrUpdateDirectory(directory)
}

func (g *Gallery) ExcludeDirectory(path string) error {
	path = g.getRelativePath(path)

	directory, err := g.Database.GetDirectory(path)
	if err != nil {
		return err
	}

	if *directory.Excluded {
		return nil
	}

	directory.Excluded = pointer.Bool(true)
	return g.Database.CreateOrUpdateDirectory(directory)
}

func (g *Gallery) CreateOrUpdateItem(item *model.Item) error {
	item.Url = g.getRelativePath(item.Url)
	item.Origin = g.getRelativePath(item.Origin)
	return g.Database.CreateOrUpdateItem(item)
}

func (g *Gallery) getRelativePath(url string) string {
	if !strings.HasPrefix(url, g.rootDirectory) {
		return url
	}

	if url == g.rootDirectory {
		return url
	}

	relativePath := strings.TrimPrefix(url, g.rootDirectory)
	return strings.TrimPrefix(relativePath, string(filepath.Separator))
}

func (g *Gallery) GetFile(url string) string {
	if strings.HasPrefix(url, string(filepath.Separator)) {
		return url
	} else {
		return filepath.Join(g.rootDirectory, url)
	}
}

func (g *Gallery) GetItemsOfTag(tag *model.Tag) (*[]model.Item, error) {
	itemIds := make([]uint64, 0)
	for _, item := range tag.Items {
		itemIds = append(itemIds, item.Id)
	}

	items, err := g.GetItems(itemIds)
	if err != nil {
		logger.Errorf("Error getting files of tag %t", err)
		return nil, err
	}

	return items, nil
}

func (g *Gallery) tagExists(tag *model.Tag, tags []*model.Tag) bool {
	for _, t := range tags {
		if tag.Id == t.Id {
			return true
		}
	}

	return false
}

func (g *Gallery) SetDirectoryTags(directory *model.Directory) error {
	existingDirectory, err := g.GetDirectory("path = ?", directory.Path)
	if err != nil {
		logger.Errorf("Error getting exising directory %s %t", directory.Path, err)
		return err
	}

	for _, tag := range existingDirectory.Tags {
		if g.tagExists(tag, directory.Tags) {
			continue
		}

		if err := g.RemoveTagFromDirectory(directory.Path, tag.Id); err != nil {
			logger.Warningf("Unable to remove tag %d from directory %s - %t",
				directory.Path, tag.Id, err)
		}
	}

	return g.CreateOrUpdateDirectory(directory)
}
