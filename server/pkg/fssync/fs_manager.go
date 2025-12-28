package fssync

import (
	"context"
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/bl/fs"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/directorytree"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"my-collection/server/pkg/utils"
	"os"
	"time"

	"github.com/go-errors/errors"
)

var logger = utils.MustGetLogger("fsmanager")

func NewFsManager(ctx context.Context, db db.Database, filesFilter directorytree.FilesFilter, checkInterval time.Duration) (*FsManager, error) {
	if err := directories.AddRootDirectory(ctx, db); err != nil {
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
	utils.PushSender
	filesFilter   directorytree.FilesFilter
	checkInterval time.Duration
	db            db.Database
	changeChannel chan bool
}

func (f *FsManager) Watch(ctx context.Context) error {
	ctx = utils.ContextWithSubject(ctx, "fs_manager")

	if err := f.Sync(ctx); err != nil {
		utils.LogError(ctx, logger, "Error in FS Watch", err)
	}

	for {
		select {
		case <-f.changeChannel:
			if err := f.Sync(ctx); err != nil {
				utils.LogError(ctx, logger, "Error in FS Watch", err)
			}
		case <-ctx.Done():
			return nil
		case <-time.After(f.checkInterval):
			if err := f.Sync(ctx); err != nil {
				utils.LogError(ctx, logger, "Error in FS Watch", err)
			}
		}
	}
}

func (f *FsManager) DirectoryChanged() {
	select {
	case f.changeChannel <- true:
	default:
	}
}

func (f *FsManager) GetBelongingItems(ctx context.Context, path string) (*[]model.Item, error) {
	return newFsDirectory(directories.NormalizeDirectoryPath(path)).getItems(ctx, wrapDb(f.db, f.db))
}

func (f *FsManager) GetBelongingItem(ctx context.Context, path string, filename string) (*model.Item, error) {
	return newFsDirectory(directories.NormalizeDirectoryPath(path)).getItem(ctx, wrapDb(f.db, f.db), filename)
}

func (f *FsManager) AddBelongingItem(ctx context.Context, item *model.Item) error {
	return newFsDirectory(item.Origin).addItem(ctx, f.db, f.db, item)
}

func (f *FsManager) GetAutoTags(ctx context.Context, path string) ([]*model.Tag, error) {
	dirs, err := fs.GetDirHierarchy(ctx, f.db, path)
	if err != nil {
		return nil, err
	}
	if len(dirs) == 0 {
		return nil, nil
	}

	title := directories.DirectoryNameToTag(dirs[len(dirs)-1].Path)
	result := make([]*model.Tag, 0)
	for _, dir := range dirs {
		tags, err := tags.GetOrCreateTags(ctx, f.db, directories.BuildDirectoryTags(dir, title))
		if err != nil {
			return nil, err
		}
		result = append(result, tags...)
	}

	return result, nil
}

func (f *FsManager) GetFileMetadata(path string) (int64, int64, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, 0, errors.Wrap(err, 1)
	}

	return file.ModTime().UnixMilli(), file.Size(), nil
}

func (f *FsManager) runSync(ctx context.Context) (bool, error) {
	dig, err := NewCachedDig(ctx, f.db, f.db)
	if err != nil {
		return false, err
	}
	fsSync, err := newFsSyncer(ctx, relativasor.GetRootDirectory(), f.db, dig, f.filesFilter)
	if err != nil {
		return false, err
	}

	hasChanges, errors := fsSync.sync(ctx, f.db, f, f, f)

	if len(errors) > 0 {
		logger.Errorf(ctx, "FS Sync finished with %d errors", len(errors))
		for _, err := range errors {
			utils.LogError(ctx, logger, "Error in FS Sync", err)
		}

		return hasChanges, errors[0]
	}

	return hasChanges, nil
}

func (f *FsManager) Sync(ctx context.Context) error {
	var lastError error
	hasAnyChange := false

	hasChanges := true
	for hasChanges {
		var err error
		hasChanges, err = f.runSync(ctx)
		if err != nil {
			lastError = err
		}
		if hasChanges {
			hasAnyChange = hasChanges
		}
	}

	if hasAnyChange {
		f.Push(model.PushMessage{MessageType: model.PUSH_FS_CHANGE, Payload: ""})
	}

	return lastError
}
