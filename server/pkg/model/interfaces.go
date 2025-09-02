package model

import (
	"time"
)

//go:generate mockgen -package model -source interfaces.go -destination interfaces_mock.go

type ItemReader interface {
	GetItem(conds ...interface{}) (*Item, error)
	GetItems(conds ...interface{}) (*[]Item, error)
	GetAllItems() (*[]Item, error)
}

type ItemWriter interface {
	CreateOrUpdateItem(item *Item) error
	UpdateItem(item *Item) error
	RemoveItem(itemId uint64) error
	RemoveTagFromItem(itemId uint64, tagId uint64) error
}

type ItemReaderWriter interface {
	ItemReader
	ItemWriter
}

type TagReader interface {
	GetTag(conds ...interface{}) (*Tag, error)
	GetTags(conds ...interface{}) (*[]Tag, error)
	GetAllTags() (*[]Tag, error)
	GetTagsWithoutChildren(conds ...interface{}) (*[]Tag, error)
}

type TagWriter interface {
	CreateOrUpdateTag(tag *Tag) error
	UpdateTag(tag *Tag) error
	RemoveTag(tagId uint64) error
	RemoveTagImageFromTag(tagId uint64, imageId uint64) error
}

type TagReaderWriter interface {
	TagReader
	TagWriter
}

type TagAnnotationReader interface {
	GetTagAnnotation(conds ...interface{}) (*TagAnnotation, error)
	GetTagAnnotations(tagId uint64) ([]TagAnnotation, error)
}

type TagAnnotationWriter interface {
	CreateTagAnnotation(tagAnnotation *TagAnnotation) error
	RemoveTag(tagId uint64) error
	RemoveTagAnnotationFromTag(tagId uint64, annotationId uint64) error
}

type TagAnnotationReaderWriter interface {
	TagAnnotationReader
	TagAnnotationWriter
}

type DirectoryReader interface {
	GetDirectory(conds ...interface{}) (*Directory, error)
	GetDirectories(conds ...interface{}) (*[]Directory, error)
	GetAllDirectories() (*[]Directory, error)
}

type DirectoryWriter interface {
	CreateOrUpdateDirectory(directory *Directory) error
	UpdateDirectory(directory *Directory) error
	RemoveDirectory(path string) error
	RemoveTagFromDirectory(direcotryPath string, tagId uint64) error
}

type DirectoryReaderWriter interface {
	DirectoryReader
	DirectoryWriter
}

type TagImageTypeReader interface {
	GetTagImageType(conds ...interface{}) (*TagImageType, error)
	GetAllTagImageTypes() (*[]TagImageType, error)
}

type TagImageTypeWriter interface {
	CreateOrUpdateTagImageType(tit *TagImageType) error
}

type TagImageTypeReaderWriter interface {
	TagImageTypeReader
	TagImageTypeWriter
}

type TagImageWriter interface {
	UpdateTagImage(image *TagImage) error
}

type StorageDownloader interface {
	IsStorageUrl(name string) bool
	GetFile(name string) string
}

type StorageUploader interface {
	GetStorageUrl(name string) string
	GetFileForWriting(name string) (string, error)
	GetTempFile() string
}

type TempFileProvider interface {
	GetTempFile() string
}

type DirectoryItemsGetter interface {
	GetBelongingItems(path string) (*[]Item, error)
	GetBelongingItem(path string, filename string) (*Item, error)
}

type DirectoryItemsSetter interface {
	AddBelongingItem(item *Item) error
}

type DirectoryItemsGetterSetter interface {
	DirectoryItemsGetter
	DirectoryItemsSetter
}

type TaskReader interface {
	GetTasks(offset int, limit int) (*[]Task, error)
	TasksCount(query interface{}, conds ...interface{}) (int64, error)
}

type ProcessorStatus interface {
	IsPaused() bool
}

type DirectoryAutoTagsGetter interface {
	GetAutoTags(path string) ([]*Tag, error)
}

type FileMetadataGetter interface {
	GetFileMetadata(f string) (int64, int64, error)
}

type CurrentTimeGetter interface {
	GetCurrentTime() time.Time
}

type Database interface {
	DirectoryReaderWriter
	TagReaderWriter
	ItemReaderWriter
}

type TagCustomCommandsReader interface {
	GetTagCustomCommand(conds ...interface{}) (*[]TagCustomCommand, error)
	GetAllTagCustomCommands() (*[]TagCustomCommand, error)
}

type DbMetadataReader interface {
	GetItemsCount() (int64, error)
	GetTagsCount() (int64, error)
	GetTotalDurationSeconds() (float64, error)
}

type PushListener interface {
	Push(m PushMessage)
}
