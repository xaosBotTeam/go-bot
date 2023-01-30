package collector

import "github.com/xaosBotTeam/go-shared-models/account"

type Abstract interface {
	Collect(acc account.Account) (account.Account, error)
	CheckCondition() bool
}

func NewInfoCollectorList() []Abstract {
	return []Abstract{
		NewNickname(),
		NewEnergyLimit(),
		NewGameId(),
	}
}
