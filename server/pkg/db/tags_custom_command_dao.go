package db

import (
	"context"
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
)

func (d *databaseImpl) CreateOrUpdateTagCustomCommand(ctx context.Context, command *model.TagCustomCommand) error {
	if command.Id == 0 && command.Title == "" {
		return errors.Errorf("invalid command, missing ('id') or ('title') %v", command)
	}

	err := d.create(ctx, command)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		if command.Id != 0 {
			return d.update(ctx, command)
		}

		existing, err := d.GetItem(ctx, "title = ?", command.Title)

		if err != nil {
			return err
		}

		command.Id = existing.Id
		return d.update(ctx, command)
	}

	return err
}

func (d *databaseImpl) GetTagCustomCommand(ctx context.Context, conds ...interface{}) (*[]model.TagCustomCommand, error) {
	var commands []model.TagCustomCommand
	err := d.handleError(d.db.WithContext(ctx).Model(model.TagCustomCommand{}).Find(&commands, conds...).Error)
	return &commands, err
}

func (d *databaseImpl) GetAllTagCustomCommands(ctx context.Context) (*[]model.TagCustomCommand, error) {
	return d.GetTagCustomCommand(ctx)
}
