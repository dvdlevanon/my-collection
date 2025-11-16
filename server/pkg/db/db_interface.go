package db

import (
	"context"
	"my-collection/server/pkg/model"
)

type Database interface {
	CreateOrUpdateDirectory(ctx context.Context, directory *model.Directory) error
	UpdateDirectory(ctx context.Context, directory *model.Directory) error
	RemoveDirectory(ctx context.Context, path string) error
	RemoveTagFromDirectory(ctx context.Context, direcotryPath string, tagId uint64) error
	GetDirectory(ctx context.Context, conds ...interface{}) (*model.Directory, error)
	GetDirectories(ctx context.Context, conds ...interface{}) (*[]model.Directory, error)
	GetAllDirectories(ctx context.Context) (*[]model.Directory, error)
	CreateOrUpdateItem(ctx context.Context, item *model.Item) error
	UpdateItem(ctx context.Context, item *model.Item) error
	RemoveItem(ctx context.Context, itemId uint64) error
	RemoveTagFromItem(ctx context.Context, itemId uint64, tagId uint64) error
	GetItem(ctx context.Context, conds ...interface{}) (*model.Item, error)
	GetItems(ctx context.Context, conds ...interface{}) (*[]model.Item, error)
	GetAllItems(ctx context.Context) (*[]model.Item, error)
	GetItemsCount(ctx context.Context) (int64, error)
	GetTotalDurationSeconds(ctx context.Context) (float64, error)
	CreateTagAnnotation(ctx context.Context, tagAnnotation *model.TagAnnotation) error
	RemoveTag(ctx context.Context, tagId uint64) error
	RemoveTagAnnotationFromTag(ctx context.Context, tagId uint64, annotationId uint64) error
	GetTagAnnotation(ctx context.Context, conds ...interface{}) (*model.TagAnnotation, error)
	GetTagAnnotations(ctx context.Context, tagId uint64) ([]model.TagAnnotation, error)
	CreateOrUpdateTagImageType(ctx context.Context, tit *model.TagImageType) error
	GetTagImageType(ctx context.Context, conds ...interface{}) (*model.TagImageType, error)
	GetTagImageTypes(ctx context.Context, conds ...interface{}) (*[]model.TagImageType, error)
	GetAllTagImageTypes(ctx context.Context) (*[]model.TagImageType, error)
	CreateOrUpdateTagCustomCommand(ctx context.Context, command *model.TagCustomCommand) error
	GetTagCustomCommand(ctx context.Context, conds ...interface{}) (*[]model.TagCustomCommand, error)
	GetAllTagCustomCommands(ctx context.Context) (*[]model.TagCustomCommand, error)
	CreateOrUpdateTag(ctx context.Context, tag *model.Tag) error
	UpdateTag(ctx context.Context, tag *model.Tag) error
	GetTag(ctx context.Context, conds ...interface{}) (*model.Tag, error)
	GetTagsWithoutChildren(ctx context.Context, conds ...interface{}) (*[]model.Tag, error)
	GetTags(ctx context.Context, conds ...interface{}) (*[]model.Tag, error)
	GetAllTags(ctx context.Context) (*[]model.Tag, error)
	RemoveTagImageFromTag(ctx context.Context, tagId uint64, imageId uint64) error
	UpdateTagImage(ctx context.Context, image *model.TagImage) error
	GetTagsCount(ctx context.Context) (int64, error)
	CreateTask(ctx context.Context, task *model.Task) error
	UpdateTask(ctx context.Context, task *model.Task) error
	RemoveTasks(ctx context.Context, conds ...interface{}) error
	TasksCount(ctx context.Context, query interface{}, conds ...interface{}) (int64, error)
	GetTasks(ctx context.Context, offset int, limit int) (*[]model.Task, error)
}
