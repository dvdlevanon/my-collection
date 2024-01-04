package fssync

import (
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

func NewFsManager(db *db.Database, filesFilter directorytree.FilesFilter, periodicCheckDuration time.Duration) (*FsManager, error) {
	return &FsManager{
		filesFilter:           filesFilter,
		periodicCheckDuration: periodicCheckDuration,
		db:                    db,
		changeChannel:         make(chan string),
	}, nil
}

type FsManager struct {
	filesFilter           directorytree.FilesFilter
	periodicCheckDuration time.Duration
	db                    *db.Database
	changeChannel         chan string
}

func (f *FsManager) Watch() {
	for {
		select {
		case <-f.changeChannel:
			f.Sync()
		case <-time.After(f.periodicCheckDuration):
			f.Sync()
		}
	}
}

func (f *FsManager) DirectoryChanged(path string) {
	f.changeChannel <- path
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
	fsSync, err := newFsSyncer(relativasor.GetRootDirectory(), f.db, f, f.filesFilter)
	if err != nil {
		return err
	}

	errors := fsSync.sync(f.db, f, f, f)

	for _, err := range errors {
		utils.LogError(err)
	}

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}
