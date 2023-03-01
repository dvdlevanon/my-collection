package itemprocessor

import (
	"fmt"
	"my-collection/server/pkg/gallery"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joncrlsn/dque"
	"github.com/op/go-logging"
	"k8s.io/utils/pointer"
)

var logger = logging.MustGetLogger("item-processor")

type ProcessorNotifier interface {
	OnTaskAdded(task *model.Task)
	OnTaskComplete(task *model.Task)
}

type ItemProcessor interface {
	Run()
	EnqueueAllItemsPreview(force bool) error
	EnqueueAllItemsCovers(force bool) error
	EnqueueAllItemsVideoMetadata(force bool) error
	EnqueueItemPreview(id uint64)
	EnqueueItemCovers(id uint64)
	EnqueueMainCover(id uint64, second float64)
	EnqueueItemVideoMetadata(id uint64)
	IsPaused() bool
	Pause()
	Continue()
	SetProcessorNotifier(notifier ProcessorNotifier)
}

func taskBuilder() interface{} {
	return &model.Task{}
}

type itemProcessorImpl struct {
	gallery      *gallery.Gallery
	storage      *storage.Storage
	dque         *dque.DQue
	pauseChannel chan bool
	paused       bool
	notifier     ProcessorNotifier
}

func New(gallery *gallery.Gallery, storage *storage.Storage) (ItemProcessor, error) {
	logger.Infof("Item processor initialized")

	tasksDirectory := storage.GetStorageDirectory("tasks")
	if err := os.MkdirAll(tasksDirectory, 0750); err != nil {
		logger.Errorf("Error creating tasks directory %s", err)
		return nil, err
	}

	dque, err := dque.NewOrOpen("tasks", tasksDirectory, 100, taskBuilder)
	if err != nil {
		logger.Errorf("Error creating tasks queue %s", err)
		return nil, err
	}

	return &itemProcessorImpl{
		gallery:      gallery,
		storage:      storage,
		dque:         dque,
		pauseChannel: make(chan bool, 10),
	}, nil
}

func (p *itemProcessorImpl) SetProcessorNotifier(notifier ProcessorNotifier) {
	p.notifier = notifier
}

func (p *itemProcessorImpl) IsPaused() bool {
	return p.paused
}

func (p *itemProcessorImpl) Pause() {
	p.pauseChannel <- true
}

func (p *itemProcessorImpl) Continue() {
	p.pauseChannel <- false
}

func (p *itemProcessorImpl) Run() {
	for {
		select {
		case paused := <-p.pauseChannel:
			logger.Infof("Queue paused changed from %t to %t", p.paused, paused)
			p.paused = paused
			if p.notifier != nil {
				p.notifier.OnTaskAdded(nil)
			}
		default:
			if !p.paused {
				p.process()
			} else {
				time.Sleep(time.Second)
			}
		}
	}
}

func (p *itemProcessorImpl) process() {
	taskIfc, err := p.dque.Peek()
	if err != nil {
		if err != dque.ErrEmpty {
			logger.Errorf("Error peeking tasks queue %s", err)
		}

		time.Sleep(time.Second)
		return
	}

	task, ok := taskIfc.(*model.Task)
	if !ok {
		logger.Errorf("Unable to convert interface to task %s", task)
		return
	}

	startMillis := time.Now().UnixMilli()
	task.ProcessingStart = pointer.Int64(time.Now().UnixMilli())
	if err := p.gallery.UpdateTask(task); err != nil {
		logger.Warningf("Unable to update task processing start time %s %s", task.Id, err)
	}

	logger.Infof("Start processing task %+v", task)
	if err := p.processTask(task); err != nil {
		logger.Errorf("Error processing task %+v for id: %d - %t", task.TaskType.String(), task.IdParam, err)
	}

	p.gallery.RemoveTask(task.Id)
	if _, err = p.dque.DequeueBlock(); err != nil {
		logger.Errorf("Error dequeuing task %s - %+v", err, task)
	}

	if p.notifier != nil {
		p.notifier.OnTaskComplete(task)
	}

	processingMillis := time.Now().UnixMilli() - startMillis
	logger.Infof("Done processing task in %dms %+v", processingMillis, task)
}

func (p *itemProcessorImpl) enqueue(t *model.Task) {
	t.Id = uuid.New().String()
	t.EnequeueTime = pointer.Int64(time.Now().UnixMilli())
	if err := p.dque.Enqueue(t); err != nil {
		logger.Errorf("Error enqueuing task %s - %v", err, *t)
		return
	}

	if err := p.gallery.CreateTask(t); err != nil {
		logger.Errorf("Error adding task to db %s - %v", err, *t)
		return
	}

	if p.notifier != nil {
		p.notifier.OnTaskAdded(t)
	}
}

func (p *itemProcessorImpl) processTask(t *model.Task) error {
	switch t.TaskType {
	case model.REFRESH_COVER_TASK:
		return p.refreshItemCovers(t.IdParam)
	case model.SET_MAIN_COVER:
		return p.setMainCover(t.IdParam, t.FloatParam)
	case model.REFRESH_PREVIEW_TASK:
		return p.refreshItemPreview(t.IdParam)
	case model.REFRESH_METADATA_TASK:
		return p.refreshItemMetadata(t.IdParam)
	default:
		return fmt.Errorf("unknown task %+v", t)
	}
}
