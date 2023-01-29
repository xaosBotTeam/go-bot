package info_collector

import "github.com/xaosBotTeam/go-shared-models/account"

type Abstract interface {
	Collect(acc account.Account) (account.Account, error)
	CheckCondition() bool
}
