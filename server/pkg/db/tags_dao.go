package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (d *databaseImpl) CreateOrUpdateTag(tag *model.Tag) error {
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

func (d *databaseImpl) UpdateTag(tag *model.Tag) error {
	return d.update(tag)
}

func (d *databaseImpl) getTagModel(withChildren bool) *gorm.DB {
	itemsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	annotationsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	model := d.db.Model(model.Tag{}).
		Preload("Images").
		Preload("Items", itemsPreloading).
		Preload("Annotations", annotationsPreloading)

	if withChildren {
		model = model.Preload("Children")
	}

	return model
}

func (d *databaseImpl) GetTag(conds ...interface{}) (*model.Tag, error) {
	tag := &model.Tag{}
	err := d.handleError(d.getTagModel(true).First(tag, conds...).Error)
	return tag, err
}

func (d *databaseImpl) GetTagsWithoutChildren(conds ...interface{}) (*[]model.Tag, error) {
	var tags []model.Tag
	err := d.handleError(d.getTagModel(false).Find(&tags, conds...).Error)
	return &tags, err
}

func (d *databaseImpl) GetTags(conds ...interface{}) (*[]model.Tag, error) {
	var tags []model.Tag
	err := d.handleError(d.getTagModel(true).Find(&tags, conds...).Error)
	return &tags, err
}

func (d *databaseImpl) GetAllTags() (*[]model.Tag, error) {
	return d.GetTags()
}

func (d *databaseImpl) RemoveTagImageFromTag(tagId uint64, imageId uint64) error {
	return d.deleteAssociation(model.Tag{Id: tagId}, model.TagImage{Id: imageId}, "Images")
}

func (d *databaseImpl) UpdateTagImage(image *model.TagImage) error {
	return d.update(image)
}

func (d *databaseImpl) GetTagsCount() (int64, error) {
	var count int64
	err := d.handleError(d.db.Model(&model.Tag{}).Count(&count).Error)
	return count, err
}
