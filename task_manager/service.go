package task_manager

import (
	models "github.com/xaosBotTeam/go-shared-models/task"
	"go-bot/storage"
	"go-bot/task"
	"log"
	"sync"
	"time"
)

func New(accounts storage.AbstractAccountStorage, statuses storage.AbstractStatusStorage) *TaskManager {
	if accounts == nil || statuses == nil {
		return nil
	}
	return &TaskManager{accountStorage: accounts, statusStorage: statuses, status: make(map[int][]task.Abstract), stop: false, update: make(map[int]bool)}
}

type TaskManager struct {
	accountStorage storage.AbstractAccountStorage
	statusStorage  storage.AbstractStatusStorage
	status         map[int][]task.Abstract
	stop           bool
	update         map[int]bool
}

func (t *TaskManager) UpdateStatus(accountId int, status models.Status) error {
	t.status[accountId] = StatusToTasks(&status)
	t.update[accountId] = true
	return t.statusStorage.Update(accountId, status)
}

func (t *TaskManager) Init() error {
	ids, statuses, err := t.statusStorage.GetAll()
	if err != nil {
		return err
	}
	for i, id := range ids {
		t.status[id] = StatusToTasks(&statuses[i])
		t.update[id] = false
	}
	return nil
}

func (t *TaskManager) Start() error {
	for {
		accounts, err := t.accountStorage.GetAll()
		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		for _, acc := range accounts {
			tasks, ok := t.status[acc.ID]
			if ok {
				if len(tasks) == 0 {
					continue
				}
				wg.Add(1)
				go func() {
					for {
						currentStatus, err := t.GetStatusById(acc.ID)
						oldStatus := currentStatus
						if err != nil {
							log.Printf("ERR: Task Manager: %s", err.Error())
						}
						for _, currentTask := range tasks {
							if currentTask.CheckCondition() {
								err = currentTask.Do(acc)
								if err != nil {
									log.Printf("ERR: Task Manager, task %s: %s\n", currentTask.GetName(), err.Error())
								}
								if !currentTask.IsPersistent() {
									currentStatus = currentTask.RemoveFromStatus(currentStatus)
								}

								log.Printf("Task %s ended", currentTask.GetName())
							}
						}

						if currentStatus != oldStatus && !t.update[acc.ID] {
							err = t.statusStorage.Update(acc.ID, currentStatus)
							t.status[acc.ID] = StatusToTasks(&currentStatus)
							if err != nil {
								log.Printf("ERR: Task Manager: %s", err.Error())
							}
						}

						if t.stop {
							wg.Done()
							return
						}

						tasks, ok = t.status[acc.ID]
						if t.update[acc.ID] {
							t.update[acc.ID] = false
						}
						if !ok || len(tasks) == 0 {
							wg.Done()
							return
						}

						time.Sleep(1 * time.Second)
					}
				}()
			}
		}
		wg.Wait()
		log.Printf("Restart tasks in 30 seconds...")
		time.Sleep(30 * time.Second)
		t.stop = false
	}
	return nil
}

func (t *TaskManager) RefreshAccounts() {
	t.stop = true
}

// может стоит возращать прям текущий-текущий из status
func (t *TaskManager) GetStatusById(id int) (models.Status, error) {
	return t.statusStorage.GetByAccId(id)
}

func (t *TaskManager) GetAllStatuses() (map[int]models.Status, error) {
	ids, statuses, err := t.statusStorage.GetAll()
	if err != nil {
		return nil, err
	}
	toReturn := make(map[int]models.Status)
	for i, id := range ids {
		toReturn[id] = statuses[i]
	}
	return toReturn, nil

}