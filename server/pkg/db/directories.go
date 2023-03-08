package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (d *Database) CreateOrUpdateDirectory(directory *model.Directory) error {
	err := d.create(directory)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		return d.update(directory)
	}

	return err
}

func (d *Database) GetDirectory(conds ...interface{}) (*model.Directory, error) {
	directory := &model.Directory{}

	tagsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	err := d.db.Model(directory).
		Preload("Tags", tagsPreloading).
		First(directory, conds...).Error
	return directory, err
}

func (d *Database) GetAllDirectories() (*[]model.Directory, error) {
	var directories []model.Directory

	tagsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	err := d.db.Model(&model.Directory{}).Preload("Tags", tagsPreloading).Find(&directories).Error
	if err != nil {
		err = errors.Wrap(err, 0)
	}

	return &directories, err
}

func (d *Database) RemoveDirectory(path string) error {
	return d.db.Delete(model.Directory{Path: path}).Error
}

func (d *Database) RemoveTagFromDirectory(direcotryPath string, tagId uint64) error {
	return d.db.Model(model.Directory{Path: direcotryPath}).Association("Tags").Delete(model.Tag{Id: tagId})
}
