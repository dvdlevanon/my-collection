package fssync

import (
	"context"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type BelongingItemDbReader interface {
	GetDirectoryTag(ctx context.Context, path string) (*model.Tag, error)
	GetItems(ctx context.Context, ids []uint64) (*[]model.Item, error)
	GetItemByTitle(ctx context.Context, tag *model.Tag, title string) (*model.Item, error)
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

func (d *dbBelongingItemDbReader) GetDirectoryTag(ctx context.Context, path string) (*model.Tag, error) {
	return tags.GetChildTag(ctx, d.tr, directories.GetDirectoriesTagId(), path)
}

func (d *dbBelongingItemDbReader) GetItemByTitle(ctx context.Context, tag *model.Tag, title string) (*model.Item, error) {
	return tags.GetItemByTitle(ctx, d.ir, tag, title)
}

func (d *dbBelongingItemDbReader) GetItems(ctx context.Context, ids []uint64) (*[]model.Item, error) {
	if len(ids) == 0 {
		result := make([]model.Item, 0)
		return &result, nil
	}

	return d.ir.GetItems(ctx, ids)
}

func newFsDirectory(path string) *FsDirectory {
	return &FsDirectory{
		path: path,
	}
}

type FsDirectory struct {
	path string
}

func (d *FsDirectory) getTag(ctx context.Context, tr BelongingItemDbReader) (*model.Tag, error) {
	tag, err := tr.GetDirectoryTag(ctx, d.path)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tag, nil
}

func (d *FsDirectory) removeItem(ctx context.Context, tr model.TagReader, iw model.ItemWriter, item *model.Item) error {
	tag, err := d.getTag(ctx, wrapDb(tr, nil))
	if err != nil {
		return err
	}
	if tag == nil {
		return errors.Errorf("'directory tag' not found %s", d.path)
	}

	return iw.RemoveTagFromItem(ctx, item.Id, tag.Id)
}

func (d *FsDirectory) addItem(ctx context.Context, tr model.TagReader, iw model.ItemWriter, item *model.Item) error {
	tag, err := d.getTag(ctx, wrapDb(tr, nil))
	if err != nil {
		return err
	}
	if tag == nil {
		return errors.Errorf("'directory tag' not found %s", d.path)
	}

	item.Tags = append(item.Tags, tag)
	return iw.CreateOrUpdateItem(ctx, item)
}

func (d *FsDirectory) getItem(ctx context.Context, tr BelongingItemDbReader, filename string) (*model.Item, error) {
	tag, err := d.getTag(ctx, tr)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, nil
	}

	return tr.GetItemByTitle(ctx, tag, items.TitleFromFileName(filename))
}

func (d *FsDirectory) getItems(ctx context.Context, tr BelongingItemDbReader) (*[]model.Item, error) {
	tag, err := d.getTag(ctx, tr)
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

	items, err := tr.GetItems(ctx, itemIds)
	if err != nil {
		return nil, err
	}

	return items, nil
}
