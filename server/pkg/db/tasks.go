package db

import "my-collection/server/pkg/model"

func (d *Database) CreateTask(task *model.Task) error {
	return d.create(task)
}

func (d *Database) UpdateTask(task *model.Task) error {
	return d.update(task)
}

func (d *Database) RemoveTasks(conds ...interface{}) error {
	return d.delete(model.Task{}, conds...)
}

func (d *Database) TasksCount(query interface{}, conds ...interface{}) (int64, error) {
	var count int64
	err := d.handleError(d.db.Model(model.Task{}).Where(query, conds...).Count(&count).Error)
	return count, err
}

func (d *Database) GetTasks(offset int, limit int) (*[]model.Task, error) {
	var tasks []model.Task
	err := d.handleError(d.db.Model(model.Task{}).Offset(offset).Limit(limit).Find(&tasks).Error)
	return &tasks, err
}
