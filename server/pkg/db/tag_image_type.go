package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
)

func (d *databaseImpl) CreateOrUpdateTagImageType(tit *model.TagImageType) error {
	if tit.Id == 0 && tit.Nickname == "" {
		return errors.Errorf("invalid tag image type, missing ('id') or ('title') %v", tit)
	}

	err := d.create(tit)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		if tit.Id != 0 {
			return d.update(tit)
		}

		existing, err := d.GetItem("nickname = ?", tit.Nickname)

		if err != nil {
			return err
		}

		tit.Id = existing.Id
		return d.update(tit)
	}

	return err
}

func (d *databaseImpl) GetTagImageType(conds ...interface{}) (*model.TagImageType, error) {
	tit := &model.TagImageType{}
	err := d.handleError(d.db.Model(tit).First(tit, conds...).Error)
	return tit, err
}

func (d *databaseImpl) GetTagImageTypes(conds ...interface{}) (*[]model.TagImageType, error) {
	var tits []model.TagImageType
	err := d.handleError(d.db.Model(model.TagImageType{}).Find(&tits, conds...).Error)
	return &tits, err
}

func (d *databaseImpl) GetAllTagImageTypes() (*[]model.TagImageType, error) {
	return d.GetTagImageTypes()
}
