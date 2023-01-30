package task_manager

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5"
	"github.com/xaosBotTeam/go-shared-models/account"
	models "github.com/xaosBotTeam/go-shared-models/task"
	"go-bot/collector"
	"go-bot/navigation"
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
	return &TaskManager{
		accountStorage: accounts,
		statusStorage:  statuses,
		status:         make(map[int][]task.Abstract),
		stop:           false, update: make(map[int]bool),
		collectors: make(map[int][]collector.Abstract),
	}
}

type TaskManager struct {
	accountStorage storage.AbstractAccountStorage
	statusStorage  storage.AbstractStatusStorage
	status         map[int][]task.Abstract
	stop           bool
	update         map[int]bool
	collectors     map[int][]collector.Abstract
}

func (t *TaskManager) UpdateStatus(accountId int, status models.Status) error {
	t.status[accountId] = StatusToTasks(&status)
	t.update[accountId] = true
	_, err := t.statusStorage.GetByAccId(accountId)
	if err != nil && err.Error() == "no rows in result set" {
		return t.statusStorage.Add(accountId, status)
	} else {
		return t.statusStorage.Update(accountId, status)
	}
}

func (t *TaskManager) initAccount(acc account.Account) error {
	status, err := t.statusStorage.GetByAccId(acc.ID)
	if err == pgx.ErrNoRows {
		status = models.Status{}
		err = t.statusStorage.Add(acc.ID, status)
	} else if err != nil {
		return err
	}

	t.status[acc.ID] = StatusToTasks(&status)
	t.update[acc.ID] = false
	t.collectors[acc.ID] = collector.NewInfoCollectorList()
	return nil
}

func (t *TaskManager) Init() error {
	accs, err := t.accountStorage.GetAll()
	if err != nil {
		return err
	}
	for _, acc := range accs {
		err = t.initAccount(acc)
	}
	return nil
}

func (t *TaskManager) Start() error {
	for {
		accounts, err := t.accountStorage.GetAll()
		if err != nil {
			return err
		}

		for _, acc := range accounts {
			t.update[acc.ID] = false
		}

		var wg sync.WaitGroup
		for _, acc := range accounts {
			tasks, ok := t.status[acc.ID]
			currentAccount := acc
			wg.Add(1)
			go func() {
				for {
					oldAccount := currentAccount
					currentAccount, err := t.IterateCollectors(currentAccount)
					if err != nil {
						log.Printf("ERR: Task Manager: Info Collecting: id %d, name %s: %s", acc.ID, acc.FriendlyName, err.Error())
					}
					if oldAccount != currentAccount {
						err = t.accountStorage.Update(currentAccount)
					}

					currentStatus, err := t.GetStatusById(currentAccount.ID)
					oldStatus := currentStatus
					if err != nil {
						log.Printf("ERR: Task Manager: try to get current status: %s", err.Error())
					}
					currentStatus, err = t.IterateTasks(currentAccount, tasks, currentStatus)
					if err != nil {
						log.Printf("ERR: Task Manager %s", err.Error())
					}

					if currentStatus != oldStatus && !t.update[currentAccount.ID] {
						err = t.statusStorage.Update(currentAccount.ID, currentStatus)
						if err != nil && err == pgx.ErrNoRows {
							err = t.statusStorage.Add(currentAccount.ID, currentStatus)
						}
						if err != nil {
							log.Printf("ERR: Task Manager: %s", err.Error())
						}
						t.status[currentAccount.ID] = StatusToTasks(&currentStatus)
						if err != nil {
							log.Printf("ERR: Task Manager: %s", err.Error())
						}
					}

					if t.stop {
						wg.Done()
						return
					}

					tasks, ok = t.status[currentAccount.ID]
					if t.update[currentAccount.ID] {
						t.update[currentAccount.ID] = false
					}
					if !ok || len(tasks) == 0 {
						wg.Done()
						return
					}

					time.Sleep(1 * time.Second)
				}
			}()
		}
		wg.Wait()
		log.Printf("Restart tasks in 30 seconds...")
		time.Sleep(30 * time.Second)
		t.stop = false
	}
}

func (t *TaskManager) RefreshAccounts() {
	t.stop = true
}

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

func (t *TaskManager) AddAccount(url string, ownerId int) (account.Account, error) {
	if !navigation.ValidateUrl(url) {
		return account.Account{}, errors.New("url is not valid")
	}

	acc, err := t.accountStorage.Add(url, ownerId)
	if err != nil {
		return account.Account{}, err
	}
	err = t.initAccount(acc)
	return acc, nil
}

func (t *TaskManager) GetAllAccounts() ([]account.Account, error) {
	return t.accountStorage.GetAll()
}

func (t *TaskManager) GetAccountById(id int) (account.Account, error) {
	return t.accountStorage.GetById(id)
}

func (t *TaskManager) IterateCollectors(acc account.Account) (account.Account, error) {
	var finalErr error
	for _, c := range t.collectors[acc.ID] {
		var err error
		if c.CheckCondition() {
			acc, err = c.Collect(acc)
			if err != nil {
				finalErr = multierror.Append(finalErr, err)
			}
		}
	}

	return acc, finalErr
}

func (t *TaskManager) IterateTasks(acc account.Account, tasks []task.Abstract, currentStatus models.Status) (models.Status, error) {
	var finalErr error
	for _, currentTask := range tasks {
		if currentTask.CheckCondition() {
			log.Printf("Task %s started on account id %d, nickname %s", currentTask.GetName(), acc.ID, acc.FriendlyName)

			err := currentTask.Do(acc)
			if err != nil {
				finalErr = multierror.Append(finalErr, err)
			}
			if !currentTask.IsPersistent() {
				currentStatus = currentTask.RemoveFromStatus(currentStatus)
			}

			log.Printf("Task %s ended on account id %d, nickname %s", currentTask.GetName(), acc.ID, acc.FriendlyName)
		}
	}
	return currentStatus, finalErr
}
