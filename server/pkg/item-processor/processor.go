package itemprocessor

import (
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/storage"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("item-processor")

type TaskType int

const (
	REFRESH_COVER_TASK = iota
	REFRESH_PREVIEW_TASK
	REFRESH_METADATA_TASK
)

type ItemProcessor interface {
	Run()
	EnqueueAllItemsPreview() error
	EnqueueAllItemsCovers() error
	EnqueueAllItemsVideoMetadata() error
	EnqueueItemPreview(id uint64)
	EnqueueItemCovers(id uint64)
	EnqueueItemVideoMetadata(id uint64)
}

type task struct {
	taskType TaskType
	id       uint64
}

type itemProcessorImpl struct {
	gallery *gallery.Gallery
	storage *storage.Storage
	queue   chan task
}

func New(gallery *gallery.Gallery, storage *storage.Storage) ItemProcessor {
	logger.Infof("Item processor initialized")

	return &itemProcessorImpl{
		gallery: gallery,
		storage: storage,
		queue:   make(chan task, 100000),
	}
}

func (p itemProcessorImpl) Run() {
	for t := range p.queue {
		p.processTask(&t)
	}
}

func (p itemProcessorImpl) processTask(t *task) {
	switch t.taskType {
	case REFRESH_COVER_TASK:
		p.handleError(t, p.refreshItemCovers(t.id))
	case REFRESH_PREVIEW_TASK:
		p.handleError(t, p.refreshItemPreview(t.id))
	case REFRESH_METADATA_TASK:
		p.handleError(t, p.refreshItemMetadata(t.id))
	}
}

func (t TaskType) String() string {
	switch t {
	case REFRESH_COVER_TASK:
		return "cover"
	case REFRESH_PREVIEW_TASK:
		return "preview"
	case REFRESH_METADATA_TASK:
		return "metadata"
	default:
		return "unknown"
	}
}

func (p itemProcessorImpl) handleError(t *task, err error) {
	if err == nil {
		return
	}

	logger.Errorf("Error processing task %v for id: %d - %t", t.taskType.String(), t.id, err)
}
