package task_manager

import (
	models "github.com/xaosBotTeam/go-shared-models/config"
	"go-bot/task"
)

func StatusToTasks(status models.Config) map[task.Type]AbstractTask {
	tasks := make(map[task.Type]AbstractTask)
	if status.ArenaFarming {
		tasks[task.ArenaBoostingTask] = task.NewArenaBoosting(status)
	}
	if status.Travelling {
		tasks[task.TravellingTask] = task.NewTravelling()
	}
	if status.OpenChests {
		tasks[task.OpenChests] = task.NewChests()
	}
	return tasks
}

func UpdateTasksWithStatus(tasks map[task.Type]AbstractTask, configuration models.Config) map[task.Type]AbstractTask {
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
	
	if configuration.OpenChests {
		tasks[task.OpenChests] = task.NewChests()
	} else {
		delete(tasks, task.OpenChests)
	}
	return tasks
}
