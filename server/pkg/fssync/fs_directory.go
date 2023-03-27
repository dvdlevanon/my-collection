package fssync

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	"os"
	"strings"

	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

func newFsDirectory(path string) *FsDirectory {
	return &FsDirectory{
		path: path,
	}
}

type FsDirectory struct {
	path string
}

func (d *FsDirectory) getTagTitle() string {
	name := strings.ReplaceAll(d.path, string(os.PathSeparator), "_")
	title := directories.DirectoryNameToTag(name)
	return title
}

func (d *FsDirectory) getTag(tr model.TagReader) (*model.Tag, error) {
	tag, err := tags.GetChildTag(tr, directories.DIRECTORIES_TAG_ID, d.getTagTitle())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tag, nil
}

func (d *FsDirectory) removeItem(tr model.TagReader, iw model.ItemWriter, item *model.Item) error {
	tag, err := d.getTag(tr)
	if err != nil {
		return err
	}
	if tag == nil {
		return errors.Errorf("'directory tag' not found %s", d.getTagTitle())
	}

	return iw.RemoveTagFromItem(item.Id, tag.Id)
}

func (d *FsDirectory) addItem(tr model.TagReader, iw model.ItemWriter, item *model.Item) error {
	tag, err := d.getTag(tr)
	if err != nil {
		return err
	}
	if tag == nil {
		return errors.Errorf("'directory tag' not found %s", d.getTagTitle())
	}

	item.Tags = append(item.Tags, tag)
	return iw.CreateOrUpdateItem(item)
}

func (d *FsDirectory) getItem(tr model.TagReader, ir model.ItemReader, filename string) (*model.Item, error) {
	tag, err := d.getTag(tr)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, nil
	}

	return tags.GetItemByTitle(ir, tag, items.TitleFromFileName(filename))
}

func (d *FsDirectory) getItems(tr model.TagReader, ir model.ItemReader) (*[]model.Item, error) {
	tag, err := d.getTag(tr)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		empty := make([]model.Item, 0)
		return &empty, nil
	}

	items, err := tags.GetItems(ir, tag)
	if err != nil {
		return nil, err
	}

	return items, nil
}
