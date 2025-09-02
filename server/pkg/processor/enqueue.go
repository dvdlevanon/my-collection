package processor

import (
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/model"
	"os"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"k8s.io/utils/pointer"
)

func (p *Processor) enqueue(t *model.Task) {
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

	p.pushQueueMetadata()
}

func (p *Processor) GetFileMetadata(path string) (int64, int64, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, 0, errors.Wrap(err, 1)
	}

	return file.ModTime().UnixMilli(), file.Size(), nil
}

func (p *Processor) EnqueueAllItemsVideoMetadata(force bool) error {
	allItems, err := p.db.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *allItems {
		if !force && item.DurationSeconds != 0 {
			modified, err := items.IsModified(&item, p)
			if !modified || err != nil {
				continue
			}
		}

		p.EnqueueItemVideoMetadata(item.Id)
	}

	return nil
}

func (p *Processor) EnqueueItemVideoMetadata(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_METADATA_TASK, IdParam: id})
}

func (p *Processor) EnqueueAllItemsPreview(force bool) error {
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

func (p *Processor) EnqueueItemPreview(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_PREVIEW_TASK, IdParam: id})
}

func (p *Processor) EnqueueAllItemsCovers(force bool) error {
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

func (p *Processor) EnqueueMainCover(id uint64, second float64) {
	p.enqueue(&model.Task{TaskType: model.SET_MAIN_COVER, IdParam: id, FloatParam: second})
}

func (p *Processor) EnqueueCropFrame(id uint64, second float64, rect model.RectFloat) {
	p.enqueue(&model.Task{TaskType: model.CROP_FRAME, IdParam: id, FloatParam: second, StringParam: rect.Serialize()})
}

func (p *Processor) EnqueueChangeResolution(id uint64, newResolution string) {
	p.enqueue(&model.Task{TaskType: model.CHANGE_RESOLUTION, IdParam: id, StringParam: newResolution})
}

func (p *Processor) EnqueueItemCovers(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_COVER_TASK, IdParam: id})
}

func (p *Processor) EnqueueAllItemsFileMetadata() error {
	allItems, err := p.db.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *allItems {
		p.EnqueueItemFileMetadata(item.Id)
	}

	return nil
}

func (p *Processor) EnqueueItemFileMetadata(id uint64) {
	p.enqueue(&model.Task{TaskType: model.REFRESH_FILE_TASK, IdParam: id})
}
