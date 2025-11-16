package db

import (
	"log"
	"my-collection/server/pkg/model"
	"os"
	"time"

	"github.com/go-errors/errors"
	"github.com/op/go-logging"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormlogger "gorm.io/gorm/logger"
)

var logger = logging.MustGetLogger("server")

type Database interface {
	CreateOrUpdateDirectory(directory *model.Directory) error
	UpdateDirectory(directory *model.Directory) error
	RemoveDirectory(path string) error
	RemoveTagFromDirectory(direcotryPath string, tagId uint64) error
	getDirectoryModel() *gorm.DB
	GetDirectory(conds ...interface{}) (*model.Directory, error)
	GetDirectories(conds ...interface{}) (*[]model.Directory, error)
	GetAllDirectories() (*[]model.Directory, error)
	CreateOrUpdateItem(item *model.Item) error
	UpdateItem(item *model.Item) error
	RemoveItem(itemId uint64) error
	RemoveTagFromItem(itemId uint64, tagId uint64) error
	GetItem(conds ...interface{}) (*model.Item, error)
	GetItems(conds ...interface{}) (*[]model.Item, error)
	GetAllItems() (*[]model.Item, error)
	GetItemsCount() (int64, error)
	GetTotalDurationSeconds() (float64, error)
	CreateTagAnnotation(tagAnnotation *model.TagAnnotation) error
	RemoveTag(tagId uint64) error
	RemoveTagAnnotationFromTag(tagId uint64, annotationId uint64) error
	GetTagAnnotation(conds ...interface{}) (*model.TagAnnotation, error)
	GetTagAnnotations(tagId uint64) ([]model.TagAnnotation, error)
	CreateOrUpdateTagImageType(tit *model.TagImageType) error
	GetTagImageType(conds ...interface{}) (*model.TagImageType, error)
	GetTagImageTypes(conds ...interface{}) (*[]model.TagImageType, error)
	GetAllTagImageTypes() (*[]model.TagImageType, error)
	CreateOrUpdateTagCustomCommand(command *model.TagCustomCommand) error
	GetTagCustomCommand(conds ...interface{}) (*[]model.TagCustomCommand, error)
	GetAllTagCustomCommands() (*[]model.TagCustomCommand, error)
	CreateOrUpdateTag(tag *model.Tag) error
	UpdateTag(tag *model.Tag) error
	getTagModel(withChildren bool) *gorm.DB
	GetTag(conds ...interface{}) (*model.Tag, error)
	GetTagsWithoutChildren(conds ...interface{}) (*[]model.Tag, error)
	GetTags(conds ...interface{}) (*[]model.Tag, error)
	GetAllTags() (*[]model.Tag, error)
	RemoveTagImageFromTag(tagId uint64, imageId uint64) error
	UpdateTagImage(image *model.TagImage) error
	GetTagsCount() (int64, error)
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task) error
	RemoveTasks(conds ...interface{}) error
	TasksCount(query interface{}, conds ...interface{}) (int64, error)
	GetTasks(offset int, limit int) (*[]model.Task, error)
}

type databaseImpl struct {
	db *gorm.DB
}

func New(dbfile string) (*databaseImpl, error) {
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,       // Slow SQL threshold
			LogLevel:                  gormlogger.Silent, // Log level
			IgnoreRecordNotFoundError: true,              // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,             // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{
		Logger: newLogger,
	})

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

	if err = db.AutoMigrate(&model.Task{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.TagImageType{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.TagImage{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if err = db.AutoMigrate(&model.TagCustomCommand{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	logger.Infof("DB initialized with db file: %s", dbfile)

	return &databaseImpl{
		db: db,
	}, nil
}

func (d *databaseImpl) handleError(err error) error {
	if err != nil {
		return errors.Wrap(err, 1)
	}

	return nil
}

func (d *databaseImpl) deleteAssociation(value interface{}, association interface{}, name string) error {
	return d.handleError(d.db.Model(value).Association(name).Delete(association))
}

func (d *databaseImpl) delete(value interface{}, conds ...interface{}) error {
	return d.handleError(d.db.Delete(value, conds...).Error)
}

func (d *databaseImpl) deleteWithAssociations(value interface{}, conds ...interface{}) error {
	return d.handleError(d.db.Select(clause.Associations).Delete(value, conds...).Error)
}

func (d *databaseImpl) create(value interface{}) error {
	return d.handleError(d.db.Create(value).Error)
}

func (d *databaseImpl) update(value interface{}) error {
	return d.handleError(d.db.Updates(value).Error)
}
