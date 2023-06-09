package task_manager

import (
	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/config"
	"github.com/xaosBotTeam/go-shared-models/status"
)

type AbstractTask interface {
	Do(account.Account, status.Status) error
	CheckCondition() bool
	IsPersistent() bool
	RemoveFromStatus(config.Config) config.Config
}

type AbstractCollector interface {
	Collect(status.Status, string) (status.Status, error)
	CheckCondition() bool
}

type AbstractAccountStorage interface {
	GetAll() (map[int]account.Account, error)
	GetById(id int) (account.Account, error)
	Close()
	Add(acc account.Account) (int, error)
	Update(id int, acc account.Account) error
	Delete(id int) error
}

type AbstractConfigStorage interface {
	GetAll() ([]int, []config.Config, error)
	GetByAccId(id int) (config.Config, error)
	Update(id int, configuration config.Config) error
	UpdateRange(configuration config.Config) error
	Delete(id int) error
	Add(id int, configuration config.Config) error
	Close()
}

type AbstractStatusStorage interface {
	GetById(id int) (status.Status, error)
	GetAll() (map[int]status.Status, error)
	Update(id int, stat status.Status) error
	Add(id int, stat status.Status) error
	Close()
	Delete(id int) error
}
