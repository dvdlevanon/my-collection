package db

import (
	"context"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (d *databaseImpl) CreateOrUpdateDirectory(ctx context.Context, directory *model.Directory) error {
	err := d.create(ctx, directory)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		return d.update(ctx, directory)
	}

	return err
}

func (d *databaseImpl) UpdateDirectory(ctx context.Context, directory *model.Directory) error {
	return d.update(ctx, directory)
}

func (d *databaseImpl) RemoveDirectory(ctx context.Context, path string) error {
	return d.delete(ctx, model.Directory{Path: path})
}

func (d *databaseImpl) RemoveTagFromDirectory(ctx context.Context, direcotryPath string, tagId uint64) error {
	return d.deleteAssociation(ctx, model.Directory{Path: direcotryPath}, model.Tag{Id: tagId}, "Tags")
}

func (d *databaseImpl) getDirectoryModel(ctx context.Context) *gorm.DB {
	tagsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	return d.db.WithContext(ctx).Model(model.Directory{}).Preload("Tags", tagsPreloading)
}

func (d *databaseImpl) GetDirectory(ctx context.Context, conds ...interface{}) (*model.Directory, error) {
	directory := &model.Directory{}
	err := d.handleError(d.getDirectoryModel(ctx).First(directory, conds...).Error)
	return directory, err
}

func (d *databaseImpl) GetDirectories(ctx context.Context, conds ...interface{}) (*[]model.Directory, error) {
	var directories []model.Directory
	err := d.handleError(d.getDirectoryModel(ctx).Find(&directories, conds...).Error)
	return &directories, err
}

func (d *databaseImpl) GetAllDirectories(ctx context.Context) (*[]model.Directory, error) {
	return d.GetDirectories(ctx)
}
