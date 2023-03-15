package items

import (
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"
	"path/filepath"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("items")

func FileExists(item model.Item) bool {
	path := relativasor.GetAbsoluteFile(filepath.Join(item.Origin, item.Title))
	_, err := os.Stat(path)
	return err == nil
}

func TitleFromFileName(path string) string {
	return filepath.Base(path)
}

func BuildItemFromPath(path string) (*model.Item, error) {
	title := TitleFromFileName(path)
	file, err := os.Stat(path)
	if err != nil {
		logger.Errorf("Error getting file stat %s - %t", file, err)
		return nil, err
	}

	return &model.Item{
		Title:        title,
		Origin:       relativasor.GetRelativePath(filepath.Dir(path)),
		Url:          path,
		LastModified: file.ModTime().UnixMilli(),
	}, nil
}

func TagExists(tags []*model.Tag, tag *model.Tag) bool {
	for _, t := range tags {
		if tag.Id == t.Id {
			return true
		}
	}

	return false
}

func EnsureItemHaveTags(iw model.ItemWriter, item *model.Item, tags []*model.Tag) error {
	missingTags := make([]*model.Tag, 0)
	for _, tag := range tags {
		if !TagExists(item.Tags, tag) {
			missingTags = append(missingTags, tag)
		}
	}

	if len(missingTags) == 0 {
		return nil
	}

	item.Tags = append(item.Tags, missingTags...)
	if err := iw.UpdateItem(item); err != nil {
		return err
	}

	return nil
}

func HasSingleTag(item *model.Item, tag *model.Tag) bool {
	return len(item.Tags) == 1 && item.Tags[0].Id == tag.Id
}

func RemoveItemAndItsAssociations(iw model.ItemWriter, item *model.Item) {
	if err := iw.RemoveItem(item.Id); err != nil {
		logger.Errorf("Unable to remove item %d - %s", item.Id, err)
	}
}
