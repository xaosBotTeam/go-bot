package task_manager

import (
	models "github.com/xaosBotTeam/go-shared-models/task"
	"go-bot/task"
)

func StatusToTasks(status *models.Status) []task.Abstract {
	tasks := make([]task.Abstract, 0)
	if status.ArenaFarming {
		tasks = append(tasks, task.NewArenaBoosting(status))
	}
	return tasks
}
