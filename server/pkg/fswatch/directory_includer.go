package fswatch

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
	processor "my-collection/server/pkg/processor"
	"my-collection/server/pkg/utils"
)

func newDirectoryIncluder(trustFileExtenssion bool, directory *model.Directory,
	trw model.TagReaderWriter, drw model.DirectoryReaderWriter,
	irw model.ItemReaderWriter, processor processor.Processor) (*directoryIncluder, error) {
	title := directories.DirectoryNameToTag(directory.Path)
	tag, err := tags.GetOrCreateChildTag(trw, directories.DIRECTORIES_TAG_ID, title)
	if err != nil {
		return nil, err
	}

	concreteTags, err := tags.GetOrCreateTags(trw, directories.BuildDirectoryTags(directory))
	if err != nil {
		return nil, err
	}

	return &directoryIncluder{
		trustFileExtenssion: trustFileExtenssion,
		directory:           directory,
		tag:                 tag,
		concreteTags:        concreteTags,
		tr:                  trw,
		irw:                 irw,
		drw:                 drw,
		processor:           processor,
	}, nil
}

type directoryIncluder struct {
	trustFileExtenssion bool
	directory           *model.Directory
	tag                 *model.Tag
	concreteTags        []*model.Tag
	tr                  model.TagReader
	irw                 model.ItemReaderWriter
	drw                 model.DirectoryReaderWriter
	processor           processor.Processor
	filesCount          int
	dirsCount           int
	anyItemAdded        bool
}

func (d *directoryIncluder) process() error {
	if err := d.processFiles(); err != nil {
		return err
	}

	if err := d.reloadTag(); err != nil {
		return err
	}

	d.removeDeletedFiles()
	// removeDeletedSubdirs

	return nil
}

func (d *directoryIncluder) processFiles() error {
	files, err := directories.GetDirectoryFiles(d.directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := directories.GetDirectoryFile(d.directory, file.Name())
		if file.IsDir() {
			handleError(d.handleChildDirectory(path))
			d.dirsCount = d.dirsCount + 1
		} else {
			handleError(d.handleChildFile(path))
			d.filesCount = d.filesCount + 1
		}
	}

	return nil
}

func (d *directoryIncluder) handleChildDirectory(path string) error {
	return directories.AddExcludedDirectoryIfMissing(d.drw, path)
}

func (d *directoryIncluder) handleChildFile(path string) error {
	newItem, err := items.BuildItemFromPath(path)
	if err != nil {
		return err
	}

	existingItem, err := tags.GetItemByTitle(d.irw, d.tag, newItem.Title)
	if err != nil {
		return err
	}

	if existingItem == nil && !utils.IsVideo(d.trustFileExtenssion, path) {
		return nil
	}

	if existingItem != nil {
		return d.processExistingItem(existingItem, newItem)
	} else {
		return d.processNewItem(newItem)
	}
}

func (d *directoryIncluder) processExistingItem(existingItem *model.Item, newItem *model.Item) error {
	if existingItem.LastModified != newItem.LastModified {
		d.processItem(existingItem)
	}

	return items.EnsureItemHaveTags(d.irw, existingItem, d.concreteTags)
}

func (d *directoryIncluder) processNewItem(item *model.Item) error {
	item.Tags = append(d.concreteTags, d.tag)

	if err := d.irw.CreateOrUpdateItem(item); err != nil {
		return err
	}

	d.anyItemAdded = true
	d.processItem(item)
	return nil
}

func (d *directoryIncluder) processItem(item *model.Item) {
	if d.processor.IsAutomaticProcessing() {
		d.processor.EnqueueItemVideoMetadata(item.Id)
		d.processor.EnqueueItemCovers(item.Id)
		d.processor.EnqueueItemPreview(item.Id)
	}
}

func (d *directoryIncluder) reloadTag() error {
	if !d.anyItemAdded {
		return nil
	}

	tag, err := d.tr.GetTag(d.tag.Id)
	if err != nil {
		return err
	}

	d.tag = tag
	return nil
}

func (d *directoryIncluder) removeDeletedFiles() error {
	belongingItems, err := tags.GetItems(d.irw, d.tag)
	if err != nil {
		return err
	}

	for _, item := range *belongingItems {
		if item.Origin != d.directory.Path {
			continue
		}

		if items.FileExists(item) {
			continue
		}

		handleError(d.irw.RemoveItem(item.Id))
		d.filesCount = d.filesCount - 1
	}

	return nil
}
