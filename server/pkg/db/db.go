package db

import (
	"my-collection/server/pkg/model"
	"path/filepath"

	"github.com/go-errors/errors"
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

	if err = db.AutoMigrate(&model.Task{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	logger.Infof("DB initialized with db file: %s", actualpath)

	return &Database{
		db: db,
	}, nil
}

func (d *Database) handleError(err error) error {
	if err != nil {
		return errors.Wrap(err, 1)
	}

	return nil
}

func (d *Database) deleteAssociation(value interface{}, association interface{}, name string) error {
	return d.handleError(d.db.Model(value).Association(name).Delete(association))
}

func (d *Database) delete(value interface{}, conds ...interface{}) error {
	return d.handleError(d.db.Delete(value, conds...).Error)
}

func (d *Database) create(value interface{}) error {
	return d.handleError(d.db.Create(value).Error)
}

func (d *Database) update(value interface{}) error {
	return d.handleError(d.db.Updates(value).Error)
}
