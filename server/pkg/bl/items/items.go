package items

import (
	"math/rand"
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

func BuildItemFromPath(origin string, path string, flmg model.FileLastModifiedGetter) (*model.Item, error) {
	lastModified, err := flmg.GetLastModified(path)
	if err != nil {
		return nil, err
	}

	title := TitleFromFileName(path)
	return &model.Item{
		Title:        title,
		Origin:       origin,
		Url:          filepath.Join(origin, title),
		LastModified: lastModified,
	}, nil
}

func UpdateFileLocation(iw model.ItemWriter, item *model.Item, origin string, path string) error {
	title := TitleFromFileName(path)
	item.Origin = origin
	item.Title = title
	item.Url = filepath.Join(origin, title)
	return iw.UpdateItem(item)
}

func ItemExists(items []*model.Item, item *model.Item) bool {
	for _, i := range items {
		if item.Id == i.Id {
			return true
		}
	}

	return false
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

func RemoveItemAndItsAssociations(iw model.ItemWriter, item *model.Item) []error {
	errors := make([]error, 0)
	if err := iw.RemoveItem(item.Id); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func GetRandomItems(ir model.ItemReader, count int) ([]*model.Item, error) {
	allItems, err := ir.GetAllItems()
	if err != nil {
		return nil, err
	}

	randomItems := make([]*model.Item, 0)
	for i := 0; i < count; i++ {
		chosenItem := &((*allItems)[rand.Intn(len(*allItems))])

		for ItemExists(randomItems, chosenItem) {
			chosenItem = &((*allItems)[rand.Intn(len(*allItems))])
		}

		randomItems = append(randomItems, chosenItem)
	}

	return randomItems, nil
}

func IsModified(item *model.Item, flmg model.FileLastModifiedGetter) (bool, error) {
	path := relativasor.GetAbsoluteFile(filepath.Join(item.Origin, item.Title))
	lastModified, err := flmg.GetLastModified(path)
	if err != nil {
		return false, err
	}

	return item.LastModified != lastModified, nil
}
