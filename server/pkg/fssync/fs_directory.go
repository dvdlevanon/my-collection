package fssync

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type BelongingItemDbReader interface {
	GetDirectoryTag(path string) (*model.Tag, error)
	GetItems(ids []uint64) (*[]model.Item, error)
	GetItemByTitle(tag *model.Tag, title string) (*model.Item, error)
}

func wrapDb(tr model.TagReader, ir model.ItemReader) BelongingItemDbReader {
	return &dbBelongingItemDbReader{
		tr: tr,
		ir: ir,
	}
}

type dbBelongingItemDbReader struct {
	tr model.TagReader
	ir model.ItemReader
}

func (d *dbBelongingItemDbReader) GetDirectoryTag(path string) (*model.Tag, error) {
	return tags.GetChildTag(d.tr, directories.GetDirectoriesTagId(), path)
}

func (d *dbBelongingItemDbReader) GetItemByTitle(tag *model.Tag, title string) (*model.Item, error) {
	return tags.GetItemByTitle(d.ir, tag, title)
}

func (d *dbBelongingItemDbReader) GetItems(ids []uint64) (*[]model.Item, error) {
	if len(ids) == 0 {
		result := make([]model.Item, 0)
		return &result, nil
	}

	return d.ir.GetItems(ids)
}

func newFsDirectory(path string) *FsDirectory {
	return &FsDirectory{
		path: path,
	}
}

type FsDirectory struct {
	path string
}

func (d *FsDirectory) getTag(tr BelongingItemDbReader) (*model.Tag, error) {
	tag, err := tr.GetDirectoryTag(d.path)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tag, nil
}

func (d *FsDirectory) removeItem(tr model.TagReader, iw model.ItemWriter, item *model.Item) error {
	tag, err := d.getTag(wrapDb(tr, nil))
	if err != nil {
		return err
	}
	if tag == nil {
		return errors.Errorf("'directory tag' not found %s", d.path)
	}

	return iw.RemoveTagFromItem(item.Id, tag.Id)
}

func (d *FsDirectory) addItem(tr model.TagReader, iw model.ItemWriter, item *model.Item) error {
	tag, err := d.getTag(wrapDb(tr, nil))
	if err != nil {
		return err
	}
	if tag == nil {
		return errors.Errorf("'directory tag' not found %s", d.path)
	}

	item.Tags = append(item.Tags, tag)
	return iw.CreateOrUpdateItem(item)
}

func (d *FsDirectory) getItem(tr BelongingItemDbReader, filename string) (*model.Item, error) {
	tag, err := d.getTag(tr)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, nil
	}

	return tr.GetItemByTitle(tag, items.TitleFromFileName(filename))
}

func (d *FsDirectory) getItems(tr BelongingItemDbReader) (*[]model.Item, error) {
	tag, err := d.getTag(tr)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		empty := make([]model.Item, 0)
		return &empty, nil
	}

	itemIds := make([]uint64, 0)
	for _, item := range tag.Items {
		itemIds = append(itemIds, item.Id)
	}

	items, err := tr.GetItems(itemIds)
	if err != nil {
		return nil, err
	}

	return items, nil
}
