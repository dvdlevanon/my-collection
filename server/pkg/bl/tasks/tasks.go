package tasks

import (
	"context"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"

	"k8s.io/utils/pointer"
)

var logger = utils.MustGetLogger("tasks")

func BuildQueueMetadata(ctx context.Context, tr model.TaskReader, ps model.ProcessorStatus) (model.QueueMetadata, error) {
	size, err := tr.TasksCount(ctx, "")
	if err != nil {
		logger.Errorf(ctx, "Unable to get queue size %s", err)
		return model.QueueMetadata{}, nil
	}

	unfinishedTasks, err := tr.TasksCount(ctx, "processing_end is null")
	if err != nil {
		logger.Errorf(ctx, "Unable to get unfinished tasks count %s", err)
		return model.QueueMetadata{}, nil
	}

	queueMetadata := model.QueueMetadata{
		Size:            pointer.Int64(size),
		Paused:          pointer.Bool(ps.IsPaused()),
		UnfinishedTasks: pointer.Int64(unfinishedTasks),
	}

	return queueMetadata, err
}
