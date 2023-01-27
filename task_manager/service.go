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
	return &TaskManager{accountStorage: accounts, statusStorage: statuses, status: make(map[int][]task.Abstract), stop: false}
}

type TaskManager struct {
	accountStorage storage.AbstractAccountStorage
	statusStorage  storage.AbstractStatusStorage
	status         map[int][]task.Abstract
	stop           bool
}

func (t *TaskManager) UpdateStatus(accountId int, status models.Status) error {
	t.status[accountId] = StatusToTasks(&status)
	return t.statusStorage.Update(accountId, status)
}

func (t *TaskManager) Init() error {
	ids, statuses, err := t.statusStorage.GetAll()
	if err != nil {
		return err
	}
	for i, id := range ids {
		t.status[id] = StatusToTasks(&statuses[i])
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
						for _, currentTask := range tasks {
							if currentTask.CheckCondition() {
								err = currentTask.Do(acc)
								if err != nil {
									log.Printf("Error while doing task %s, error: %s\n", currentTask.GetName(), err.Error())
								}
								log.Printf("Task %s ended", currentTask.GetName())
							}
						}
						if t.stop {
							wg.Done()
							return
						}

						tasks, ok = t.status[acc.ID]
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
		t.stop = false
	}
	return nil
}

func (t *TaskManager) RefreshAccounts() {
	t.stop = true
}
