package db

import (
	"my-collection/server/pkg/model"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var logger = logging.MustGetLogger("server")

type Database struct {
	db *gorm.DB
}

func New(rootDirectory string, filename string) (*Database, error) {
	actualpath := filepath.Join(rootDirectory, filename)
	db, err := gorm.Open(sqlite.Open(actualpath), &gorm.Config{})

	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.Item{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.Tag{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.Cover{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.TagAnnotation{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.Directory{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	logger.Infof("DB initialized with db file: %s", actualpath)

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

func (d *Database) CreateTagAnnotation(tagAnnotation *model.TagAnnotation) error {
	return d.create(tagAnnotation)
}

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

func (d *Database) CreateOrUpdateItem(item *model.Item) error {
	if item.Id == 0 && (item.Title == "" || item.Origin == "") {
		return errors.Errorf("Invalid item, missing id or title and origin %v", item)
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

	annotationsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	err := d.db.Model(tag).
		Preload("Children").
		Preload("Items", itemsPreloading).
		Preload("Annotations", annotationsPreloading).
		First(tag, conds...).Error

	if err != nil {
		return nil, err
	}

	return tag, err
}

func (d *Database) GetTags(conds ...interface{}) (*[]model.Tag, error) {
	var tags []model.Tag

	itemsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	annotationsPreloading := func(db *gorm.DB) *gorm.DB {
		return db.Select("ID")
	}

	err := d.db.Model(model.Tag{}).
		Preload("Children").
		Preload("Items", itemsPreloading).
		Preload("Annotations", annotationsPreloading).
		Find(&tags, conds...).Error
	return &tags, err
}

func (d *Database) GetAllTags() (*[]model.Tag, error) {
	return d.GetTags()
}

func (d *Database) GetTagAnnotation(conds ...interface{}) (*model.TagAnnotation, error) {
	tagAnnotation := &model.TagAnnotation{}
	err := d.db.Model(tagAnnotation).First(tagAnnotation, conds...).Error
	return tagAnnotation, err
}

func (d *Database) getItemModel(includeTagIdsOnly bool) *gorm.DB {
	tagsPreloading := func(db *gorm.DB) *gorm.DB {
		if includeTagIdsOnly {
			return db.Select("ID")
		} else {
			return db
		}
	}

	return d.db.Model(&model.Item{}).Preload("Tags", tagsPreloading).Preload("Covers")
}

func (d *Database) GetItem(conds ...interface{}) (*model.Item, error) {
	item := &model.Item{}
	err := d.getItemModel(false).First(item, conds...).Error

	if err != nil {
		err = errors.Wrap(err, 0)
	}

	return item, err
}

func (d *Database) GetItems(conds ...interface{}) (*[]model.Item, error) {
	var items []model.Item
	err := d.getItemModel(true).Find(&items, conds...).Error

	if err != nil {
		err = errors.Wrap(err, 0)
	}

	return &items, err
}

func (d *Database) GetAllItems() (*[]model.Item, error) {
	return d.GetItems()
}

func (d *Database) RemoveTagFromItem(itemId uint64, tagId uint64) error {
	return d.db.Model(&model.Item{Id: itemId}).Association("Tags").Delete(model.Tag{Id: tagId})
}

func (d *Database) RemoveTagAnnotationFromTag(tagId uint64, annotationId uint64) error {
	return d.db.Model(&model.Tag{Id: tagId}).Association("Annotations").Delete(model.TagAnnotation{Id: annotationId})
}

func (d *Database) GetTagAnnotations(tagId uint64) ([]model.TagAnnotation, error) {
	var annotations []model.TagAnnotation
	err := d.db.Model(&model.Tag{Id: tagId}).Association("Annotations").Find(&annotations)

	if err != nil {
		err = errors.Wrap(err, 0)
	}

	return annotations, err
}

func (d *Database) RemoveTag(tagId uint64) error {
	return d.db.Delete(model.Tag{Id: tagId}).Error
}

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

func (d *Database) RemoveItem(itemId uint64) error {
	return d.db.Delete(model.Item{Id: itemId}).Error
}
