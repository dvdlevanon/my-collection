package model

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
	GetItem(conds ...interface{}) (*Item, error)
	UpdateItem(item *Item) error
}

type TagReader interface {
	GetTag(conds ...interface{}) (*Tag, error)
	GetTags(conds ...interface{}) (*[]Tag, error)
	GetAllTags() (*[]Tag, error)
}

type TagWriter interface {
	CreateOrUpdateTag(tag *Tag) error
	UpdateTag(tag *Tag) error
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
	CreateOrUpdateDirectory(directory *Directory) error
	RemoveDirectory(path string) error
	RemoveTagFromDirectory(direcotryPath string, tagId uint64) error
}

type DirectoryWriter interface {
	GetDirectory(conds ...interface{}) (*Directory, error)
	GetDirectories(conds ...interface{}) (*[]Directory, error)
	GetAllDirectories() (*[]Directory, error)
}

type DirectoryReaderWriter interface {
	DirectoryReader
	DirectoryWriter
}

type StorageUploader interface {
	GetStorageUrl(name string) string
	GetFileForWriting(name string) (string, error)
	GetTempFile() string
}
