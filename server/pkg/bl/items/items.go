package items

import (
	"fmt"
	"math/rand"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("items")

const PREVIEW_FROM_START_POSITION = "start-position" //items-util.js

type ItemsFilter func(item *model.Item) bool

func FileExists(item model.Item) bool {
	path := relativasor.GetAbsoluteFile(filepath.Join(item.Origin, item.Title))
	_, err := os.Stat(path)
	return err == nil
}

func TitleFromFileName(path string) string {
	return filepath.Base(path)
}

func BuildItemFromPath(origin string, path string, fmdg model.FileMetadataGetter) (*model.Item, error) {
	lastModified, fileSize, err := fmdg.GetFileMetadata(path)
	if err != nil {
		return nil, err
	}

	title := TitleFromFileName(path)
	return &model.Item{
		Title:        title,
		Origin:       origin,
		Url:          filepath.Join(origin, title),
		LastModified: lastModified,
		FileSize:     fileSize,
	}, nil
}

func UpdateFileLocation(iw model.ItemWriter, item *model.Item, origin string, path string, url string) error {
	title := TitleFromFileName(path)
	item.Origin = origin
	item.Title = title

	if url == "" {
		item.Url = filepath.Join(origin, title)
	} else {
		item.Url = url
	}

	for _, highlight := range item.Highlights {
		highlightOrigin := buildHighlightUrl(origin, highlight.StartPosition, highlight.EndPosition)
		UpdateFileLocation(iw, highlight, highlightOrigin, path, item.Url)
	}

	for _, subitem := range item.SubItems {
		subitemOrigin := buildSubItemOrigin(origin, subitem.StartPosition, subitem.EndPosition)
		UpdateFileLocation(iw, subitem, subitemOrigin, path, item.Url)
	}

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

func EnsureItemHaveTags(iw model.ItemWriter, item *model.Item, tags []*model.Tag) (bool, error) {
	missingTags := make([]*model.Tag, 0)
	for _, tag := range tags {
		if !TagExists(item.Tags, tag) {
			missingTags = append(missingTags, tag)
		}
	}

	if len(missingTags) == 0 {
		return false, nil
	}

	item.Tags = append(item.Tags, missingTags...)
	if err := iw.UpdateItem(item); err != nil {
		return false, err
	}

	return true, nil
}

func EnsureItemMissingTags(iw model.ItemWriter, item *model.Item, tags []*model.Tag) error {
	for _, tagToRemove := range tags {
		for _, tag := range item.Tags {
			if tag.Id == tagToRemove.Id {
				if err := iw.RemoveTagFromItem(item.Id, tag.Id); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func HasSingleTag(item *model.Item, tag *model.Tag) bool {
	return len(item.Tags) == 1 && item.Tags[0].Id == tag.Id
}

func RemoveItemAndItsAssociations(iw model.ItemWriter, itemId uint64) []error {
	errors := make([]error, 0)
	if err := iw.RemoveItem(itemId); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func DeleteRealFile(ir model.ItemReader, itemId uint64) error {
	item, err := ir.GetItem(itemId)
	if err != nil {
		return err
	}

	if IsSubItem(item) {
		return errors.Errorf("Deletion of subitem is forbidden")
	}

	if IsHighlight(item) {
		return errors.Errorf("Deletion of highlight is forbidden")
	}

	file := relativasor.GetAbsoluteFile(item.Url)
	logger.Infof("About to delete real file %s", file)
	if err := os.Remove(file); err != nil {
		logger.Warningf("Unable to delete file %s, %s", file, err)
	}

	return nil
}

func noRandom(item *model.Item) bool {
	for _, tag := range item.Tags {
		if tag.NoRandom != nil && *tag.NoRandom {
			return true
		}
	}

	return false
}

func GetRandomItems(ir model.ItemReader, count int, filter ItemsFilter) ([]*model.Item, error) {
	allItems, err := ir.GetAllItems()
	if err != nil {
		return nil, err
	}
	if len(*allItems) == 0 {
		return nil, fmt.Errorf("no items")
	}
	if count >= len(*allItems) {
		count = len(*allItems) - 1
	}

	randomItems := make([]*model.Item, 0)
	for i := 0; i < count; i++ {
		chosenItem := &((*allItems)[rand.Intn(len(*allItems))])

		for ItemExists(randomItems, chosenItem) || !filter(chosenItem) || noRandom(chosenItem) {
			chosenItem = &((*allItems)[rand.Intn(len(*allItems))])
		}

		randomItems = append(randomItems, chosenItem)
	}

	return randomItems, nil
}

func IsModified(item *model.Item, fmg model.FileMetadataGetter) (bool, error) {
	path := relativasor.GetAbsoluteFile(filepath.Join(item.Origin, item.Title))
	lastModified, _, err := fmg.GetFileMetadata(path)
	if err != nil {
		return false, err
	}

	return item.LastModified != lastModified, nil
}
