package task

import (
	account "github.com/xaosBotTeam/go-shared-models/account"
	models "github.com/xaosBotTeam/go-shared-models/task"
)

type Abstract interface {
	CheckCondition() bool
	Do(acc account.Account) error
	IsPersistent() bool
	GetName() string
	RemoveFromStatus(status models.Status) models.Status
}
