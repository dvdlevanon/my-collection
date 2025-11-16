package model

import (
	"context"
	"time"
)

//go:generate mockgen -package model -source interfaces.go -destination interfaces_mock.go

type ItemReader interface {
	GetItem(ctx context.Context, conds ...interface{}) (*Item, error)
	GetItems(ctx context.Context, conds ...interface{}) (*[]Item, error)
	GetAllItems(ctx context.Context) (*[]Item, error)
}

type ItemWriter interface {
	CreateOrUpdateItem(ctx context.Context, item *Item) error
	UpdateItem(ctx context.Context, item *Item) error
	RemoveItem(ctx context.Context, itemId uint64) error
	RemoveTagFromItem(ctx context.Context, itemId uint64, tagId uint64) error
}

type ItemReaderWriter interface {
	ItemReader
	ItemWriter
}

type TagReader interface {
	GetTag(ctx context.Context, conds ...interface{}) (*Tag, error)
	GetTags(ctx context.Context, conds ...interface{}) (*[]Tag, error)
	GetAllTags(ctx context.Context) (*[]Tag, error)
	GetTagsWithoutChildren(ctx context.Context, conds ...interface{}) (*[]Tag, error)
}

type TagWriter interface {
	CreateOrUpdateTag(ctx context.Context, tag *Tag) error
	UpdateTag(ctx context.Context, tag *Tag) error
	RemoveTag(ctx context.Context, tagId uint64) error
	RemoveTagImageFromTag(ctx context.Context, tagId uint64, imageId uint64) error
}

type TagReaderWriter interface {
	TagReader
	TagWriter
}

type TagAnnotationReader interface {
	GetTagAnnotation(ctx context.Context, conds ...interface{}) (*TagAnnotation, error)
	GetTagAnnotations(ctx context.Context, tagId uint64) ([]TagAnnotation, error)
}

type TagAnnotationWriter interface {
	CreateTagAnnotation(ctx context.Context, tagAnnotation *TagAnnotation) error
	RemoveTag(ctx context.Context, tagId uint64) error
	RemoveTagAnnotationFromTag(ctx context.Context, tagId uint64, annotationId uint64) error
}

type TagAnnotationReaderWriter interface {
	TagAnnotationReader
	TagAnnotationWriter
}

type DirectoryReader interface {
	GetDirectory(ctx context.Context, conds ...interface{}) (*Directory, error)
	GetDirectories(ctx context.Context, conds ...interface{}) (*[]Directory, error)
	GetAllDirectories(ctx context.Context) (*[]Directory, error)
}

type DirectoryWriter interface {
	CreateOrUpdateDirectory(ctx context.Context, directory *Directory) error
	UpdateDirectory(ctx context.Context, directory *Directory) error
	RemoveDirectory(ctx context.Context, path string) error
	RemoveTagFromDirectory(ctx context.Context, direcotryPath string, tagId uint64) error
}

type DirectoryReaderWriter interface {
	DirectoryReader
	DirectoryWriter
}

type TagImageTypeReader interface {
	GetTagImageType(ctx context.Context, conds ...interface{}) (*TagImageType, error)
	GetAllTagImageTypes(ctx context.Context) (*[]TagImageType, error)
}

type TagImageTypeWriter interface {
	CreateOrUpdateTagImageType(ctx context.Context, tit *TagImageType) error
}

type TagImageTypeReaderWriter interface {
	TagImageTypeReader
	TagImageTypeWriter
}

type TagImageWriter interface {
	UpdateTagImage(ctx context.Context, image *TagImage) error
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
	GetBelongingItems(ctx context.Context, path string) (*[]Item, error)
	GetBelongingItem(ctx context.Context, path string, filename string) (*Item, error)
}

type DirectoryItemsSetter interface {
	AddBelongingItem(ctx context.Context, item *Item) error
}

type DirectoryItemsGetterSetter interface {
	DirectoryItemsGetter
	DirectoryItemsSetter
}

type TaskReader interface {
	GetTasks(ctx context.Context, offset int, limit int) (*[]Task, error)
	TasksCount(ctx context.Context, query interface{}, conds ...interface{}) (int64, error)
}

type ProcessorStatus interface {
	IsPaused() bool
}

type DirectoryAutoTagsGetter interface {
	GetAutoTags(ctx context.Context, path string) ([]*Tag, error)
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
	GetTagCustomCommand(ctx context.Context, conds ...interface{}) (*[]TagCustomCommand, error)
	GetAllTagCustomCommands(ctx context.Context) (*[]TagCustomCommand, error)
}

type DbMetadataReader interface {
	GetItemsCount(ctx context.Context) (int64, error)
	GetTagsCount(ctx context.Context) (int64, error)
	GetTotalDurationSeconds(ctx context.Context) (float64, error)
}

type PushListener interface {
	Push(m PushMessage)
}
