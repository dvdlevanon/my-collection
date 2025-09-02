package tasks

import (
	"my-collection/server/pkg/model"

	"github.com/op/go-logging"
	"k8s.io/utils/pointer"
)

var logger = logging.MustGetLogger("tasks")

func AddDescriptionToTasks(ir model.ItemReader, tasks *[]model.Task) {
	for i, task := range *tasks {
		item, err := ir.GetItem(task.IdParam)
		if err != nil {
			(*tasks)[i].Description = task.TaskType.ToDescription("Unknown")
		}

		(*tasks)[i].Description = task.TaskType.ToDescription(item.Title)
	}
}

func BuildQueueMetadata(tr model.TaskReader, ps model.ProcessorStatus) (model.QueueMetadata, error) {
	size, err := tr.TasksCount("")
	if err != nil {
		logger.Errorf("Unable to get queue size %s", err)
		return model.QueueMetadata{}, nil
	}

	unfinishedTasks, err := tr.TasksCount("processing_end is null")
	if err != nil {
		logger.Errorf("Unable to get unfinished tasks count %s", err)
		return model.QueueMetadata{}, nil
	}

	queueMetadata := model.QueueMetadata{
		Size:            pointer.Int64(size),
		Paused:          pointer.Bool(ps.IsPaused()),
		UnfinishedTasks: pointer.Int64(unfinishedTasks),
	}

	return queueMetadata, err
}
