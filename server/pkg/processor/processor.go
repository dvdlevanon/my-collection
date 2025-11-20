package processor

import (
	"context"
	"fmt"
	"my-collection/server/pkg/bl/tasks"
	"my-collection/server/pkg/db"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/storage"
	general_tasks "my-collection/server/pkg/tasks/general"
	video_tasks "my-collection/server/pkg/tasks/videos"
	"my-collection/server/pkg/utils"
	"os"
	"time"

	"github.com/go-errors/errors"
	"github.com/joncrlsn/dque"
	"github.com/op/go-logging"
	"k8s.io/utils/ptr"
)

var logger = logging.MustGetLogger("item-processor")

func taskBuilder() interface{} {
	return &model.Task{}
}

type Processor struct {
	utils.PushSender
	db                   db.Database
	storage              *storage.Storage
	dque                 *dque.DQue
	pauseChannel         chan bool
	paused               bool
	coversCount          int
	previewSceneCount    int
	previewSceneDuration int
	automaticProcessing  bool
}

func New(db db.Database, storage *storage.Storage, paused bool, coversCount int, previewSceneCount int, previewSceneDuration int) (*Processor, error) {
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

	return &Processor{
		db:                   db,
		storage:              storage,
		dque:                 dque,
		coversCount:          coversCount,
		previewSceneCount:    previewSceneCount,
		previewSceneDuration: previewSceneDuration,
		paused:               paused,
		automaticProcessing:  false,
		pauseChannel:         make(chan bool, 10),
	}, nil
}

func (p *Processor) pushQueueMetadata(ctx context.Context) error {
	queueMetadata, err := tasks.BuildQueueMetadata(ctx, p.db, p)
	if err != nil {
		return err
	}

	p.Push(model.PushMessage{MessageType: model.PUSH_QUEUE_METADATA, Payload: queueMetadata})
	return nil
}

func (p *Processor) ClearFinishedTasks(ctx context.Context) error {
	if err := p.db.RemoveTasks(ctx, "processing_end is not null"); err != nil {
		logger.Errorf("Unable to clear finished tasks %s", err)
		return err
	}

	p.pushQueueMetadata(ctx)
	return nil
}

func (p *Processor) IsPaused() bool {
	return p.paused
}
func (p *Processor) IsAutomaticProcessing() bool {
	return p.automaticProcessing
}

func (p *Processor) Pause() {
	p.pauseChannel <- true
}

func (p *Processor) Continue() {
	p.pauseChannel <- false
}

func (p *Processor) GetFileMetadata(path string) (int64, int64, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, 0, errors.Wrap(err, 1)
	}

	return file.ModTime().UnixMilli(), file.Size(), nil
}

func (p *Processor) Run(ctx context.Context) error {
	ctx = utils.ContextWithSubject(ctx, "processor")
	for {
		select {
		case paused := <-p.pauseChannel:
			logger.Infof("Queue paused changed from %t to %t", p.paused, paused)
			p.paused = paused
			p.pushQueueMetadata(ctx)
		case <-ctx.Done():
			return nil
		default:
			if !p.paused {
				p.process(ctx)
			} else {
				time.Sleep(time.Second)
			}
		}
	}
}

func (p *Processor) process(ctx context.Context) {
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
	task.ProcessingStart = ptr.To(time.Now().UnixMilli())
	if err := p.db.UpdateTask(ctx, task); err != nil {
		logger.Warningf("Unable to update task processing start time %s %s", task.Id, err)
	}

	logger.Infof("Start processing task %+v", task)
	if err := p.processTask(ctx, task); err != nil {
		logger.Errorf("Error processing task %+v for params: %s - %s", task.TaskType.String(), task.Params, err)
	}

	task.ProcessingEnd = ptr.To(time.Now().UnixMilli())
	if err := p.db.UpdateTask(ctx, task); err != nil {
		logger.Warningf("Unable to update task processing end time %s %s", task.Id, err)
	}

	if _, err = p.dque.DequeueBlock(); err != nil {
		logger.Errorf("Error dequeuing task %s - %+v", err, task)
	}

	p.pushQueueMetadata(ctx)

	processingMillis := time.Now().UnixMilli() - startMillis
	logger.Infof("Done processing task in %dms %+v", processingMillis, task)
}

func (p *Processor) processTask(ctx context.Context, t *model.Task) error {
	switch t.TaskType {
	case model.REFRESH_COVER_TASK:
		return video_tasks.RefreshVideoCovers(ctx, p.db, p.storage, t.Params)
	case model.SET_MAIN_COVER:
		return video_tasks.UpdateVideoMainCover(ctx, p.db, p.storage, t.Params)
	case model.CROP_FRAME:
		return video_tasks.CropVideoFrame(ctx, p.db, p.storage, t.Params)
	case model.REFRESH_PREVIEW_TASK:
		return video_tasks.RefreshVideoPreview(ctx, p.db, p.storage, t.Params)
	case model.REFRESH_METADATA_TASK:
		return video_tasks.UpdateVideoMetadata(ctx, p.db, t.Params)
	case model.REFRESH_FILE_TASK:
		return general_tasks.UpdateFileMetadata(ctx, p.db, t.Params)
	case model.CHANGE_RESOLUTION:
		return video_tasks.ChangeVideoResolution(ctx, p.db, p.storage, t.Params)
	default:
		return fmt.Errorf("unknown task %+v", t)
	}
}
