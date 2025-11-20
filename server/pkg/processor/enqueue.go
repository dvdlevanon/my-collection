package processor

import (
	"context"
	"my-collection/server/pkg/model"
	general_tasks "my-collection/server/pkg/tasks/general"
	video_tasks "my-collection/server/pkg/tasks/videos"
	"time"

	"github.com/google/uuid"
	"k8s.io/utils/ptr"
)

func (p *Processor) enqueue(ctx context.Context, t *model.Task) error {
	t.Id = uuid.New().String()
	t.EnequeueTime = ptr.To(time.Now().UnixMilli())
	if err := p.dque.Enqueue(t); err != nil {
		return err
	}

	if err := p.db.CreateTask(ctx, t); err != nil {
		return err
	}

	return p.pushQueueMetadata(ctx)
}

func createTask(taskType model.TaskType, params string, desc string) *model.Task {
	return &model.Task{TaskType: taskType, Description: desc, Params: params}
}

func (p *Processor) EnqueueItemVideoMetadata(ctx context.Context, id uint64, title string) error {
	params, err := video_tasks.MarshalVideoMetadataParams(id)
	if err != nil {
		return err
	}

	desc := video_tasks.MetadataDesc(id, title)
	return p.enqueue(ctx, createTask(model.REFRESH_METADATA_TASK, params, desc))
}

func (p *Processor) EnqueueItemPreview(ctx context.Context, id uint64, title string) error {
	params, err := video_tasks.MarshalVideoPreviewParams(id, p.previewSceneCount, p.previewSceneDuration)
	if err != nil {
		return err
	}

	desc := video_tasks.PreviewDesc(id, title, p.previewSceneCount, p.previewSceneDuration)
	return p.enqueue(ctx, createTask(model.REFRESH_PREVIEW_TASK, params, desc))
}

func (p *Processor) EnqueueMainCover(ctx context.Context, id uint64, second float64, title string) error {
	params, err := video_tasks.MarshalMainCoverParams(id, second)
	if err != nil {
		return err
	}

	desc := video_tasks.MainCoverDesc(id, title, second)
	return p.enqueue(ctx, createTask(model.SET_MAIN_COVER, params, desc))
}

func (p *Processor) EnqueueCropFrame(ctx context.Context, id uint64, second float64, rect model.RectFloat, title string) error {
	params, err := video_tasks.MarshalVideoCropParams(id, second, rect)
	if err != nil {
		return err
	}

	desc := video_tasks.CropDesc(id, title, second, rect)
	return p.enqueue(ctx, createTask(model.CROP_FRAME, params, desc))
}

func (p *Processor) EnqueueChangeResolution(ctx context.Context, id uint64, w int, h int, title string) error {
	params, err := video_tasks.MarshalVideoResolutionParams(id, w, h)
	if err != nil {
		return err
	}

	desc := video_tasks.ResolutionDesc(id, title, w, h)
	return p.enqueue(ctx, createTask(model.CHANGE_RESOLUTION, params, desc))
}

func (p *Processor) EnqueueItemCovers(ctx context.Context, id uint64, title string) error {
	params, err := video_tasks.MarshalVideoCoversParams(id, p.coversCount)
	if err != nil {
		return err
	}

	desc := video_tasks.CoversDesc(id, title, p.coversCount)
	return p.enqueue(ctx, createTask(model.REFRESH_COVER_TASK, params, desc))
}

func (p *Processor) EnqueueItemFileMetadata(ctx context.Context, id uint64, title string) error {
	params, err := general_tasks.MarshalFileMetadataParams(id)
	if err != nil {
		return err
	}

	desc := general_tasks.MetadataDesc(id, title)
	return p.enqueue(ctx, createTask(model.REFRESH_FILE_TASK, params, desc))
}
