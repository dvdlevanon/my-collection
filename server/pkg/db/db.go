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

var logger = logging.MustGetLogger("db")

type databaseImpl struct {
	db *gorm.DB
}

func New(dbfile string, shouldLog bool) (Database, error) {
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

	result := &databaseImpl{
		db: db,
	}
	if shouldLog {
		return &dbLogger{
			db: result,
		}, nil
	}

	return result, nil
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
