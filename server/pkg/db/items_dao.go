package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (d *databaseImpl) CreateOrUpdateItem(item *model.Item) error {
	if item.Id == 0 && (item.Title == "" || item.Origin == "") {
		return errors.Errorf("invalid item, missing ('id') or ('title' and 'origin') %v", item)
	}

	err := d.create(item)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		if item.Id != 0 {
			return d.update(item)
		}

		existing, err := d.GetItem("title = ? and origin = ?", item.Title, item.Origin)

		if err != nil {
			return err
		}

		item.Id = existing.Id
		return d.update(item)
	}

	return err
}

func (d *databaseImpl) UpdateItem(item *model.Item) error {
	return d.update(item)
}

func (d *databaseImpl) RemoveItem(itemId uint64) error {
	return d.deleteWithAssociations(model.Item{Id: itemId})
}

func (d *databaseImpl) RemoveTagFromItem(itemId uint64, tagId uint64) error {
	return d.deleteAssociation(model.Item{Id: itemId}, model.Tag{Id: tagId}, "Tags")
}

func (d *databaseImpl) getItemModel(includeTagIdsOnly bool) *gorm.DB {
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

	return d.db.Model(&model.Item{}).
		Preload("Tags", tagsPreloading).
		Preload("Covers").
		Preload("SubItems", subItemsPreloading).
		Preload("Highlights", highlightsPreloading)
}

func (d *databaseImpl) GetItem(conds ...interface{}) (*model.Item, error) {
	item := &model.Item{}
	err := d.handleError(d.getItemModel(false).First(item, conds...).Error)
	return item, err
}

func (d *databaseImpl) GetItems(conds ...interface{}) (*[]model.Item, error) {
	var items []model.Item
	err := d.handleError(d.getItemModel(false).Find(&items, conds...).Error)
	return &items, err
}

func (d *databaseImpl) GetAllItems() (*[]model.Item, error) {
	return d.GetItems()
}

func (d *databaseImpl) GetItemsCount() (int64, error) {
	var count int64
	err := d.handleError(d.db.Model(&model.Item{}).Count(&count).Error)
	return count, err
}

func (d *databaseImpl) GetTotalDurationSeconds() (float64, error) {
	var result struct {
		Total float64
	}

	err := d.handleError(d.db.Model(&model.Item{}).Select("sum(duration_seconds) as total").Scan(&result).Error)
	return result.Total, err
}
