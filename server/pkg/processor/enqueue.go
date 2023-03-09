package processor

import (
	"my-collection/server/pkg/model"
	"time"

	"github.com/google/uuid"
	"k8s.io/utils/pointer"
)

func (p *itemProcessorImpl) enqueue(t *model.Task) {
	t.Id = uuid.New().String()
	t.EnequeueTime = pointer.Int64(time.Now().UnixMilli())
	if err := p.dque.Enqueue(t); err != nil {
		logger.Errorf("Error enqueuing task %s - %v", err, *t)
		return
	}

	if err := p.db.CreateTask(t); err != nil {
		logger.Errorf("Error adding task to db %s - %v", err, *t)
		return
	}

	if p.notifier != nil {
		p.notifier.OnTaskAdded(t)
	}
}

func (p itemProcessorImpl) EnqueueAllItemsVideoMetadata(force bool) error {
	items, err := p.db.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *items {
		if !force && item.DurationSeconds != 0 {
			continue
		}

		p.EnqueueItemVideoMetadata(item.Id)
	}

	return nil
}

func (p itemProcessorImpl) EnqueueItemVideoMetadata(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_METADATA_TASK, IdParam: id})
}

func (p itemProcessorImpl) EnqueueAllItemsPreview(force bool) error {
	items, err := p.db.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *items {
		if !force && item.PreviewUrl != "" {
			continue
		}

		p.EnqueueItemPreview(item.Id)
	}

	return nil
}

func (p itemProcessorImpl) EnqueueItemPreview(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_PREVIEW_TASK, IdParam: id})
}

func (p itemProcessorImpl) EnqueueAllItemsCovers(force bool) error {
	items, err := p.db.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *items {
		if !force && len(item.Covers) >= p.coversCount {
			continue
		}

		p.EnqueueItemCovers(item.Id)
	}

	return nil
}

func (p itemProcessorImpl) EnqueueMainCover(id uint64, second float64) {
	p.enqueue(&model.Task{TaskType: model.SET_MAIN_COVER, IdParam: id, FloatParam: second})
}

func (p itemProcessorImpl) EnqueueItemCovers(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_COVER_TASK, IdParam: id})
}
