package db

import (
	"my-collection/server/pkg/model"

	"github.com/go-errors/errors"
	"github.com/mattn/go-sqlite3"
)

func (d *databaseImpl) CreateOrUpdateTagCustomCommand(command *model.TagCustomCommand) error {
	if command.Id == 0 && command.Title == "" {
		return errors.Errorf("invalid command, missing ('id') or ('title') %v", command)
	}

	err := d.create(command)

	if err != nil && err.(*errors.Error).Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
		if command.Id != 0 {
			return d.update(command)
		}

		existing, err := d.GetItem("title = ?", command.Title)

		if err != nil {
			return err
		}

		command.Id = existing.Id
		return d.update(command)
	}

	return err
}

func (d *databaseImpl) GetTagCustomCommand(conds ...interface{}) (*[]model.TagCustomCommand, error) {
	var commands []model.TagCustomCommand
	err := d.handleError(d.db.Model(model.TagCustomCommand{}).Find(&commands, conds...).Error)
	return &commands, err
}

func (d *databaseImpl) GetAllTagCustomCommands() (*[]model.TagCustomCommand, error) {
	return d.GetTagCustomCommand()
}
