package fssync

import (
	"context"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/directorytree"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/utils"
	"os"
	"time"

	"github.com/go-errors/errors"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fsmanager")

func NewFsManager(db *db.Database, filesFilter directorytree.FilesFilter, checkInterval time.Duration) (*FsManager, error) {
	if err := directories.AddRootDirectory(db); err != nil {
		return nil, err
	}

	return &FsManager{
		filesFilter:   filesFilter,
		checkInterval: checkInterval,
		db:            db,
		changeChannel: make(chan bool),
	}, nil
}

type FsManager struct {
	filesFilter   directorytree.FilesFilter
	checkInterval time.Duration
	db            *db.Database
	changeChannel chan bool
}

func (f *FsManager) Watch(ctx context.Context) {
	f.Sync()
	for {
		select {
		case <-f.changeChannel:
			f.Sync()
		case <-ctx.Done():
			return
		case <-time.After(f.checkInterval):
			f.Sync()
		}
	}
}

func (f *FsManager) DirectoryChanged() {
	select {
	case f.changeChannel <- true:
	default:
	}
}

func (f *FsManager) GetBelongingItems(path string) (*[]model.Item, error) {
	return newFsDirectory(directories.NormalizeDirectoryPath(path)).getItems(f.db, f.db)
}

func (f *FsManager) GetBelongingItem(path string, filename string) (*model.Item, error) {
	return newFsDirectory(directories.NormalizeDirectoryPath(path)).getItem(f.db, f.db, filename)
}

func (f *FsManager) AddBelongingItem(item *model.Item) error {
	return newFsDirectory(item.Origin).addItem(f.db, f.db, item)
}

func (f *FsManager) GetConcreteTags(path string) ([]*model.Tag, error) {
	directory, err := directories.GetDirectory(f.db, path)
	if err != nil {
		return nil, err
	}

	return tags.GetOrCreateTags(f.db, directories.BuildDirectoryTags(directory))
}

func (f *FsManager) GetFileMetadata(path string) (int64, int64, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, 0, errors.Wrap(err, 1)
	}

	return file.ModTime().UnixMilli(), file.Size(), nil
}

func (f *FsManager) Sync() error {
	hasChanges := true
	var allErrors []error

	for hasChanges {
		fsSync, err := newFsSyncer(relativasor.GetRootDirectory(), f.db, f, f.filesFilter)
		if err != nil {
			return err
		}

		var errors []error
		hasChanges, errors = fsSync.sync(f.db, f, f, f)

		for _, err := range errors {
			utils.LogError(err)
		}
		allErrors = append(allErrors, errors...)
	}

	if len(allErrors) > 0 {
		return allErrors[0]
	}

	return nil
}
