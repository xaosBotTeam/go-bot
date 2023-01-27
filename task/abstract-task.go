package task

import account "github.com/xaosBotTeam/go-shared-models/dbAccountInformation"

type Abstract interface {
	CheckCondition() bool
	Do(acc account.DbAccountInformation) error
	GetName() string
}
