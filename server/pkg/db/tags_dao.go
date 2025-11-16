package db

import (
	"context"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (d *databaseImpl) CreateOrUpdateTag(ctx context.Context, tag *model.Tag) error {
	if tag.Id == 0 && tag.Title == "" {
		return errors.Errorf("Invalid tag, missing id or title %v", tag)
	}

	err := d.create(ctx, tag)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		if tag.Id != 0 {
			return d.update(ctx, tag)
		}

		existing, err := d.GetTag(ctx, "title = ? and parent_id = ?", tag.Title, tag.ParentID)

		if err != nil {
			return err
		}

		tag.Id = existing.Id
		return d.update(ctx, tag)
	}

	return err
}

func (d *databaseImpl) UpdateTag(ctx context.Context, tag *model.Tag) error {
	return d.update(ctx, tag)
}

func (d *databaseImpl) getTagModel(ctx context.Context, withChildren bool) *gorm.DB {
	itemsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	annotationsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	model := d.db.WithContext(ctx).Model(model.Tag{}).
		Preload("Images").
		Preload("Items", itemsPreloading).
		Preload("Annotations", annotationsPreloading)

	if withChildren {
		model = model.Preload("Children")
	}

	return model
}

func (d *databaseImpl) GetTag(ctx context.Context, conds ...interface{}) (*model.Tag, error) {
	tag := &model.Tag{}
	err := d.handleError(d.getTagModel(ctx, true).First(tag, conds...).Error)
	return tag, err
}

func (d *databaseImpl) GetTagsWithoutChildren(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	var tags []model.Tag
	err := d.handleError(d.getTagModel(ctx, false).Find(&tags, conds...).Error)
	return &tags, err
}

func (d *databaseImpl) GetTags(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	var tags []model.Tag
	err := d.handleError(d.getTagModel(ctx, true).Find(&tags, conds...).Error)
	return &tags, err
}

func (d *databaseImpl) GetAllTags(ctx context.Context) (*[]model.Tag, error) {
	return d.GetTags(ctx)
}

func (d *databaseImpl) RemoveTagImageFromTag(ctx context.Context, tagId uint64, imageId uint64) error {
	return d.deleteAssociation(ctx, model.Tag{Id: tagId}, model.TagImage{Id: imageId}, "Images")
}

func (d *databaseImpl) UpdateTagImage(ctx context.Context, image *model.TagImage) error {
	return d.update(ctx, image)
}

func (d *databaseImpl) GetTagsCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.handleError(d.db.WithContext(ctx).Model(&model.Tag{}).Count(&count).Error)
	return count, err
}
