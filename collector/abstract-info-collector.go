package collector

import (
	models "github.com/xaosBotTeam/go-shared-models/status"
)

type Abstract interface {
	Collect(config models.Status, url string) (models.Status, error)
	CheckCondition() bool
}

func NewInfoCollectorList() []Abstract {
	return []Abstract{
		NewNickname(),
		NewEnergyLimit(),
		NewGameId(),
	}
}
