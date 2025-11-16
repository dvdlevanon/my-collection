package db

import (
	"my-collection/server/pkg/model"
)

type Database interface {
	CreateOrUpdateDirectory(directory *model.Directory) error
	UpdateDirectory(directory *model.Directory) error
	RemoveDirectory(path string) error
	RemoveTagFromDirectory(direcotryPath string, tagId uint64) error
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
