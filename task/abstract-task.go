package task

import (
	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/config"
	"github.com/xaosBotTeam/go-shared-models/status"
)

type Abstract interface {
	CheckCondition() bool
	Do(acc account.Account, status status.Status) error
	IsPersistent() bool
	GetName() string
	RemoveFromStatus(configuration config.Config) config.Config
}
