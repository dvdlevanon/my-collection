package db

import (
	"context"
	"my-collection/server/pkg/model"
)

func (d *databaseImpl) CreateTask(ctx context.Context, task *model.Task) error {
	return d.create(ctx, task)
}

func (d *databaseImpl) UpdateTask(ctx context.Context, task *model.Task) error {
	return d.update(ctx, task)
}

func (d *databaseImpl) RemoveTasks(ctx context.Context, conds ...any) error {
	return d.delete(ctx, model.Task{}, conds...)
}

func (d *databaseImpl) TasksCount(ctx context.Context, query any, conds ...any) (int64, error) {
	var count int64
	err := d.handleError(d.db.WithContext(ctx).Model(model.Task{}).Where(query, conds...).Count(&count).Error)
	return count, err
}

func (d *databaseImpl) GetTasks(ctx context.Context, offset int, limit int) (*[]model.Task, error) {
	var tasks []model.Task
	err := d.handleError(d.db.WithContext(ctx).Model(model.Task{}).Offset(offset).Limit(limit).Find(&tasks).Error)
	return &tasks, err
}
