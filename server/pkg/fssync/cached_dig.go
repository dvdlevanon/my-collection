package fssync

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/model"
)

func NewCachedDig(tr model.TagReader, ir model.ItemReader) (*CachedDig, error) {
	tags, err := tr.GetAllTags()
	if err != nil {
		return nil, err
	}

	tagsMap := make(map[string]model.Tag)
	for _, tag := range *tags {
		if tag.ParentID == nil {
			continue
		}
		if *tag.ParentID == directories.GetDirectoriesTagId() {
			tagsMap[tag.Title] = tag
		}
	}

	items, err := ir.GetAllItems()
	if err != nil {
		return nil, err
	}

	itemsMap := make(map[uint64]model.Item)
	for _, item := range *items {
		itemsMap[item.Id] = item
	}

	return &CachedDig{
		tags:  tagsMap,
		items: itemsMap,
	}, nil
}

type CachedDig struct {
	tags  map[string]model.Tag
	items map[uint64]model.Item
}

func (d *CachedDig) GetBelongingItems(path string) (*[]model.Item, error) {
	return newFsDirectory(directories.NormalizeDirectoryPath(path)).getItems(d)
}

func (d *CachedDig) GetBelongingItem(path string, filename string) (*model.Item, error) {
	return newFsDirectory(directories.NormalizeDirectoryPath(path)).getItem(d, filename)
}

func (d *CachedDig) GetDirectoryTag(path string) (*model.Tag, error) {
	tag, ok := d.tags[path]
	if !ok {
		return nil, nil
	}
	return &tag, nil
}

func (d *CachedDig) GetItems(ids []uint64) (*[]model.Item, error) {
	result := make([]model.Item, 0)
	for _, id := range ids {
		item, ok := d.items[id]
		if !ok {
			continue
		}
		result = append(result, item)
	}

	return &result, nil
}

func (d *CachedDig) GetItemByTitle(tag *model.Tag, title string) (*model.Item, error) {
	return nil, nil
}
