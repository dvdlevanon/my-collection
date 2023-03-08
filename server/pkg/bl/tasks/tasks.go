package tasks

import "my-collection/server/pkg/model"

func AddDescriptionToTasks(ir model.ItemReader, tasks *[]model.Task) {
	for i, task := range *tasks {
		item, err := ir.GetItem(task.IdParam)
		if err != nil {
			(*tasks)[i].Description = task.TaskType.ToDescription("Unknown")
		}

		(*tasks)[i].Description = task.TaskType.ToDescription(item.Title)
	}
}
