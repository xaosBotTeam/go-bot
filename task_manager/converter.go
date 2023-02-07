package task_manager

import (
	models "github.com/xaosBotTeam/go-shared-models/config"
	"go-bot/task"
)

func StatusToTasks(status models.Config) map[task.Type]task.Abstract {
	tasks := make(map[task.Type]task.Abstract)
	if status.ArenaFarming {
		tasks[task.ArenaBoostingTask] = task.NewArenaBoosting(status)
	}
	if status.Travelling {
		tasks[task.TravellingTask] = task.NewTravelling()
	}
	return tasks
}

func UpdateTasksWithStatus(tasks map[task.Type]task.Abstract, configuration models.Config) map[task.Type]task.Abstract {
	if configuration.Travelling {
		if _, ok := tasks[task.TravellingTask]; !ok {
			tasks[task.TravellingTask] = task.NewTravelling()
		}
	} else {
		delete(tasks, task.TravellingTask)
	}
	if configuration.ArenaFarming {
		tasks[task.ArenaBoostingTask] = task.NewArenaBoosting(configuration)
	} else {
		delete(tasks, task.ArenaBoostingTask)
	}
	return tasks
}