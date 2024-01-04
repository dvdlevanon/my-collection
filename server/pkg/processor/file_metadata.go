package processor

import (
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"os"

	"github.com/go-errors/errors"
)

func refreshFileMetadata(irw model.ItemReaderWriter, id uint64) error {
	item, err := irw.GetItem(id)
	if err != nil {
		return err
	}

	if err := updateFileMetadata(item); err != nil {
		return err
	}

	return irw.UpdateItem(item)
}

func updateFileMetadata(item *model.Item) error {
	path := relativasor.GetAbsoluteFile(item.Url)
	logger.Infof("Refreshing file metadata for item %d  [file: %s]", item.Id, path)

	file, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	item.LastModified = file.ModTime().UnixMilli()
	item.FileSize = file.Size()
	return nil
}
