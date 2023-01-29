package info_collector

import (
	"github.com/xaosBotTeam/go-shared-models/account"
	"go-bot/navigation"
	"strconv"
	"time"
)

func NewEnergyLimit() *EnergyLimit {
	return &EnergyLimit{
		lastSync: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

type EnergyLimit struct {
	lastSync time.Time
}

func (e *EnergyLimit) Collect(acc account.Account) (account.Account, error) {
	doc, err := navigation.GetPage(acc.URL)
	if err != nil {
		return account.Account{}, err
	}
	energyLimitStr, err := navigation.GetTopBarValue(doc, 1)
	if err != nil {
		return account.Account{}, err
	}
	energyLimit, err := strconv.Atoi(energyLimitStr)
	if err != nil {
		return account.Account{}, err
	}
	if acc.EnergyLimit < energyLimit {
		acc.EnergyLimit = energyLimit
	}
	return acc, nil
}

func (e *EnergyLimit) CheckCondition() bool {
	if time.Now().Sub(e.lastSync) > 1*time.Hour {
		e.lastSync = time.Now()
		return true
	}
	return false
}
