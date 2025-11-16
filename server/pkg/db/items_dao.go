package db

import (
	"context"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (d *databaseImpl) CreateOrUpdateItem(ctx context.Context, item *model.Item) error {
	if item.Id == 0 && (item.Title == "" || item.Origin == "") {
		return errors.Errorf("invalid item, missing ('id') or ('title' and 'origin') %v", item)
	}

	err := d.create(ctx, item)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		if item.Id != 0 {
			return d.update(ctx, item)
		}

		existing, err := d.GetItem(ctx, "title = ? and origin = ?", item.Title, item.Origin)

		if err != nil {
			return err
		}

		item.Id = existing.Id
		return d.update(ctx, item)
	}

	return err
}

func (d *databaseImpl) UpdateItem(ctx context.Context, item *model.Item) error {
	return d.update(ctx, item)
}

func (d *databaseImpl) RemoveItem(ctx context.Context, itemId uint64) error {
	return d.deleteWithAssociations(ctx, model.Item{Id: itemId})
}

func (d *databaseImpl) RemoveTagFromItem(ctx context.Context, itemId uint64, tagId uint64) error {
	return d.deleteAssociation(ctx, model.Item{Id: itemId}, model.Tag{Id: tagId}, "Tags")
}

func (d *databaseImpl) getItemModel(ctx context.Context, includeTagIdsOnly bool) *gorm.DB {
	tagsPreloading := func(db *gorm.DB) *gorm.DB {
		if includeTagIdsOnly {
			return db.Select("ID", "ParentID")
		} else {
			return db.Preload("Images")
		}
	}

	subItemsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Preload("Covers").Preload("Tags", tagsPreloading)
	}

	highlightsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Preload("Covers").Preload("Tags", tagsPreloading)
	}

	return d.db.WithContext(ctx).Model(&model.Item{}).
		Preload("Tags", tagsPreloading).
		Preload("Covers").
		Preload("SubItems", subItemsPreloading).
		Preload("Highlights", highlightsPreloading)
}

func (d *databaseImpl) GetItem(ctx context.Context, conds ...interface{}) (*model.Item, error) {
	item := &model.Item{}
	err := d.handleError(d.getItemModel(ctx, false).First(item, conds...).Error)
	return item, err
}

func (d *databaseImpl) GetItems(ctx context.Context, conds ...interface{}) (*[]model.Item, error) {
	var items []model.Item
	err := d.handleError(d.getItemModel(ctx, false).Find(&items, conds...).Error)
	return &items, err
}

func (d *databaseImpl) GetAllItems(ctx context.Context) (*[]model.Item, error) {
	return d.GetItems(ctx)
}

func (d *databaseImpl) GetItemsCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.handleError(d.db.Model(&model.Item{}).WithContext(ctx).Count(&count).Error)
	return count, err
}

func (d *databaseImpl) GetTotalDurationSeconds(ctx context.Context) (float64, error) {
	var result struct {
		Total float64
	}

	err := d.handleError(d.db.Model(&model.Item{}).WithContext(ctx).Select("sum(duration_seconds) as total").Scan(&result).Error)
	return result.Total, err
}
