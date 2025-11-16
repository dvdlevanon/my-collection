package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (d *databaseImpl) CreateOrUpdateDirectory(directory *model.Directory) error {
	err := d.create(directory)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		return d.update(directory)
	}

	return err
}

func (d *databaseImpl) UpdateDirectory(directory *model.Directory) error {
	return d.update(directory)
}

func (d *databaseImpl) RemoveDirectory(path string) error {
	return d.delete(model.Directory{Path: path})
}

func (d *databaseImpl) RemoveTagFromDirectory(direcotryPath string, tagId uint64) error {
	return d.deleteAssociation(model.Directory{Path: direcotryPath}, model.Tag{Id: tagId}, "Tags")
}

func (d *databaseImpl) getDirectoryModel() *gorm.DB {
	tagsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	return d.db.Model(model.Directory{}).Preload("Tags", tagsPreloading)
}

func (d *databaseImpl) GetDirectory(conds ...interface{}) (*model.Directory, error) {
	directory := &model.Directory{}
	err := d.handleError(d.getDirectoryModel().First(directory, conds...).Error)
	return directory, err
}

func (d *databaseImpl) GetDirectories(conds ...interface{}) (*[]model.Directory, error) {
	var directories []model.Directory
	err := d.handleError(d.getDirectoryModel().Find(&directories, conds...).Error)
	return &directories, err
}

func (d *databaseImpl) GetAllDirectories() (*[]model.Directory, error) {
	return d.GetDirectories()
}
