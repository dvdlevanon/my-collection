package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func New(filename string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})

	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.Item{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.Tag{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.Preview{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) create(value interface{}) error {
	if err := d.db.Create(value).Error; err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) update(value interface{}) error {
	if err := d.db.Updates(value).Error; err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) CreateTag(tag *model.Tag) error {
	return d.create(tag)
}

func (d *Database) CreateOrUpdateTag(tag *model.Tag) error {
	err := d.create(tag)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		existing, err := d.GetTag("title = ?", tag.Title)

		if err != nil {
			return err
		}

		tag.Id = existing.Id
		return d.update(tag)
	}

	return err
}

func (d *Database) CreateItem(item *model.Item) error {
	return d.create(item)
}

func (d *Database) CreateOrUpdateItem(item *model.Item) error {
	err := d.create(item)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		existing, err := d.GetItem("title = ?", item.Title)

		if err != nil {
			return err
		}

		item.Id = existing.Id
		return d.update(item)
	}

	return err
}

func (d *Database) UpdateItem(item *model.Item) error {
	return d.update(item)
}

func (d *Database) UpdateTag(tag *model.Tag) error {
	return d.update(tag)
}

func (d *Database) GetTag(conds ...interface{}) (*model.Tag, error) {
	tag := &model.Tag{}

	itemsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	err := d.db.Model(tag).Preload("Children").Preload("Items", itemsPreloading).Find(tag, conds...).Error

	return tag, err
}

func (d *Database) getItemModel() *gorm.DB {
	tagsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	return d.db.Model(&model.Item{}).Preload("Tags", tagsPreloading).Preload("Previews")
}

func (d *Database) GetItem(conds ...interface{}) (*model.Item, error) {
	item := &model.Item{}
	err := d.getItemModel().Find(item, conds...).Error

	if err != nil {
		err = errors.Wrap(err, 0)
	}

	return item, err
}

func (d *Database) GetAllTags() (*[]model.Tag, error) {
	var tags []model.Tag

	itemsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	err := d.db.Model(model.Tag{}).Preload("Children").Preload("Items", itemsPreloading).Find(&tags).Error

	if err != nil {
		err = errors.Wrap(err, 0)
	}

	return &tags, err
}

func (d *Database) GetAllItems() (*[]model.Item, error) {
	var items []model.Item
	err := d.getItemModel().Find(&items).Error

	if err != nil {
		err = errors.Wrap(err, 0)
	}

	return &items, err
}
