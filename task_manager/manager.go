package task_manager

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5"
	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/config"
	models "github.com/xaosBotTeam/go-shared-models/status"
	"go-bot/collector"
	"go-bot/navigation"
	"go-bot/task"
	"log"
	"time"
)

func New(accounts AbstractAccountStorage, configs AbstractConfigStorage,
	statuses AbstractStatusStorage) *TaskManager {
	if accounts == nil || configs == nil || statuses == nil {
		return nil
	}
	return &TaskManager{
		accountStorage: accounts,
		configStorage:  configs,
		statusStorage:  statuses,
		update:         NewUpdateDetector(),
		collectors:     make(map[int][]AbstractCollector),
		accountQueue:   NewAccountQueue(),
		stoppers:       make(map[int]chan bool),
	}
}

type TaskManager struct {
	accountStorage AbstractAccountStorage
	configStorage  AbstractConfigStorage
	statusStorage  AbstractStatusStorage
	update         *UpdateDetector
	collectors     map[int][]AbstractCollector
	accountQueue   AbstractAccountQueue
	stoppers       map[int]chan bool
}

func (t *TaskManager) UpdateConfig(accountId int, configuration config.Config) error {
	t.update.Store(accountId, true)
	_, err := t.configStorage.GetByAccId(accountId)
	if err != nil && err == pgx.ErrNoRows {
		return t.configStorage.Add(accountId, configuration)
	} else {
		return t.configStorage.Update(accountId, configuration)
	}
}

func (t *TaskManager) initAccount(id int) error {
	_, err := t.configStorage.GetByAccId(id)
	if err == pgx.ErrNoRows {
		err = t.configStorage.Add(id, config.Config{})
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	_, err = t.statusStorage.GetById(id)
	if err == pgx.ErrNoRows {
		err = t.statusStorage.Add(id, models.Status{})
		if err != nil {
            return err
        }
	} else if err != nil {
		return err
	}

	t.update.Store(id, false)
	t.collectors[id] = []AbstractCollector {
		collector.NewEnergyLimit(),
        collector.NewGameId(),
        collector.NewNickname(),
    }
	return nil
}

func (t *TaskManager) Init() error {
	accounts, err := t.accountStorage.GetAll()
	if err != nil {
		return err
	}
	for id := range accounts {
		err = t.initAccount(id)
		if err != nil {
			return err
		}
		t.accountQueue.Add(id)
		t.stoppers[id] = make(chan bool)
	}
	return nil
}

func (t *TaskManager) Start() error {
	for {
		id := t.accountQueue.Get()
		go t.servingLoop(id, t.stoppers[id])
	}

}

func (t *TaskManager) ConfigById(id int) (config.Config, error) {
	return t.configStorage.GetByAccId(id)
}

func (t *TaskManager) AllConfigs() (map[int]config.Config, error) {
	ids, configs, err := t.configStorage.GetAll()
	if err != nil {
		return nil, err
	}
	toReturn := make(map[int]config.Config)
	for i, id := range ids {
		toReturn[id] = configs[i]
	}
	return toReturn, nil
}

func (t *TaskManager) AddAccount(acc account.Account) (int, error) {
	if !navigation.ValidateUrl(acc.URL) {
		return 0, errors.New("url is not valid")
	}

	newId, err := t.accountStorage.Add(acc)
	if err != nil {
		return 0, err
	}
	err = t.initAccount(newId)
	if err != nil {
		return 0, err
	}
	t.accountQueue.Add(newId)
	t.stoppers[newId] = make(chan bool)
	return newId, err
}

func (t *TaskManager) GetAllAccounts() (map[int]account.Account, error) {
	return t.accountStorage.GetAll()
}

func (t *TaskManager) GetAccountById(id int) (account.Account, error) {
	return t.accountStorage.GetById(id)
}

func (t *TaskManager) iterateCollectors(id int, stat models.Status, acc account.Account) (models.Status, error) {
	var finalErr error
	for _, c := range t.collectors[id] {
		var err error
		if c.CheckCondition() {
			stat, err = c.Collect(stat, acc.URL)
			if err != nil {
				finalErr = multierror.Append(finalErr, err)
			}
		}
	}

	return stat, finalErr
}

func (t *TaskManager) iterateTasks(id int, acc account.Account, tasks map[task.Type]AbstractTask, currentStatus models.Status, configuration config.Config) (config.Config, error) {
	var finalErr error
	for _, currentTask := range tasks {
		if currentTask.CheckCondition() {
			log.Printf("Task %T started on account id %d, nickname %s", currentTask, id, currentStatus.FriendlyName)

			err := currentTask.Do(acc, currentStatus)
			if err != nil {
				finalErr = multierror.Append(finalErr, err)
			}
			if !currentTask.IsPersistent() {
				configuration = currentTask.RemoveFromStatus(configuration)
			}

			log.Printf("Task %T ended on account id %d, nickname %s", currentTask, id, currentStatus.FriendlyName)
		}
	}
	return configuration, finalErr
}

func (t *TaskManager) servingLoop(id int, stop chan bool) {
	configuration, err := t.configStorage.GetByAccId(id)
	if err != nil {
		log.Printf("ERR: Task Manager: %s", err.Error())
	}
	stat, err := t.statusStorage.GetById(id)
	if err != nil {
		log.Printf("ERR: Task Manager: %s", err.Error())
	}
	acc, err := t.accountStorage.GetById(id)
	if err != nil {
		log.Printf("ERR: Task Manager: %s", err.Error())
	}
	tasks := StatusToTasks(configuration)
	for {
		select {
		case isStopped := <- t.stoppers[id]:
			if isStopped {
				return
			}
		default:
			oldStatus := stat
			stat, err = t.iterateCollectors(id, stat, acc)
			if err != nil {
				log.Printf("ERR: Task Manager: Info Collecting: id %d, name %s: %s", id, stat.FriendlyName, err.Error())
			}
			if oldStatus != stat {
				err = t.statusStorage.Update(id, stat)
				if err != nil {
					log.Printf("ERR: Task Manager: updating main")
				}
			}
			oldConfig := configuration
			if err != nil {
				log.Printf("ERR: Task Manager: try to get current status: %s", err.Error())
			}
			configuration, err = t.iterateTasks(id, acc, tasks, stat, configuration)
			if err != nil {
				log.Printf("ERR: Task Manager %s", err.Error())
			}
			if configuration != oldConfig && !t.update.Load(id) {
				err = t.configStorage.Update(id, configuration)
				if err != nil && err == pgx.ErrNoRows {
					err = t.configStorage.Add(id, configuration)
				}
				if err != nil {
					log.Printf("ERR: Task Manager: %s", err.Error())
				}
				if err != nil {
					log.Printf("ERR: Task Manager: %s", err.Error())
				}
			}

			configuration, err = t.configStorage.GetByAccId(id)
			if err != nil {
				log.Printf("ERR: Task Manager: %s", err.Error())
			}
			tasks = UpdateTasksWithStatus(tasks, configuration)
			
			if t.update.Load(id) {
				t.update.Store(id, false)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (t *TaskManager) SetConfigForAllAccount(configuration config.Config) error {
	return t.configStorage.UpdateRange(configuration)
}

func (t *TaskManager) GetStatus(id int) (models.Status, error) {
	return t.statusStorage.GetById(id)
}

func (t *TaskManager) GetAllStatuses() (map[int]models.Status, error) {
	return t.statusStorage.GetAll()
}

func (t *TaskManager) RemoveAccount(id int) error {
	var finalErr error
	t.stoppers[id] <- true
	finalErr = t.configStorage.Delete(id)
	err := t.accountStorage.Delete(id)
	if err != nil {
		finalErr = multierror.Append(finalErr, err)
	}
	err = t.statusStorage.Delete(id)
	if err != nil {
		finalErr = multierror.Append(finalErr, err)
	}
	t.update.Delete(id)
	delete(t.collectors, id)
	delete(t.stoppers, id)
	return finalErr
}
