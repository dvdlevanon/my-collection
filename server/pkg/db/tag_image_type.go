package db

import (
	"context"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
)

func (d *databaseImpl) CreateOrUpdateTagImageType(ctx context.Context, tit *model.TagImageType) error {
	if tit.Id == 0 && tit.Nickname == "" {
		return errors.Errorf("invalid tag image type, missing ('id') or ('title') %v", tit)
	}

	err := d.create(ctx, tit)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		if tit.Id != 0 {
			return d.update(ctx, tit)
		}

		existing, err := d.GetItem(ctx, "nickname = ?", tit.Nickname)

		if err != nil {
			return err
		}

		tit.Id = existing.Id
		return d.update(ctx, tit)
	}

	return err
}

func (d *databaseImpl) GetTagImageType(ctx context.Context, conds ...interface{}) (*model.TagImageType, error) {
	tit := &model.TagImageType{}
	err := d.handleError(d.db.WithContext(ctx).Model(tit).First(tit, conds...).Error)
	return tit, err
}

func (d *databaseImpl) GetTagImageTypes(ctx context.Context, conds ...interface{}) (*[]model.TagImageType, error) {
	var tits []model.TagImageType
	err := d.handleError(d.db.WithContext(ctx).Model(model.TagImageType{}).Find(&tits, conds...).Error)
	return &tits, err
}

func (d *databaseImpl) GetAllTagImageTypes(ctx context.Context) (*[]model.TagImageType, error) {
	return d.GetTagImageTypes(ctx)
}
