package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (d *Database) CreateOrUpdateTag(tag *model.Tag) error {
	if tag.Id == 0 && tag.Title == "" {
		return errors.Errorf("Invalid tag, missing id or title %v", tag)
	}

	err := d.create(tag)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		if tag.Id != 0 {
			return d.update(tag)
		}

		existing, err := d.GetTag("title = ? and parent_id = ?", tag.Title, tag.ParentID)

		if err != nil {
			return err
		}

		tag.Id = existing.Id
		return d.update(tag)
	}

	return err
}

func (d *Database) UpdateTag(tag *model.Tag) error {
	return d.update(tag)
}

func (d *Database) getTagModel() *gorm.DB {
	itemsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	annotationsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	return d.db.Model(model.Tag{}).
		Preload("Children").
		Preload("Images").
		Preload("Items", itemsPreloading).
		Preload("Annotations", annotationsPreloading)
}

func (d *Database) GetTag(conds ...interface{}) (*model.Tag, error) {
	tag := &model.Tag{}
	err := d.handleError(d.getTagModel().First(tag, conds...).Error)
	return tag, err
}

func (d *Database) GetTags(conds ...interface{}) (*[]model.Tag, error) {
	var tags []model.Tag
	err := d.handleError(d.getTagModel().Find(&tags, conds...).Error)
	return &tags, err
}

func (d *Database) GetAllTags() (*[]model.Tag, error) {
	return d.GetTags()
}
