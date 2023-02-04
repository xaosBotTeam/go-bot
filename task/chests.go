package task

import (
	"time"

	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/config"
	"github.com/xaosBotTeam/go-shared-models/status"
)

type Chests struct {
	lastSync time.Time
}

func (c *Chests) CheckCondition() bool {
	if time.Since(c.lastSync) >= 240 * time.Minute {
		c.lastSync = time.Now()
		return true
	}
	return false
}

func (c *Chests) Do(acc account.Account, status status.Status) error {
	
}

func (c *Chests) IsPersistent() bool {
	return true
}

func (c *Chests) RemoveFromStatus(configuration config.Config) config.Config {
	panic("not implemented") // TODO: Implement
}

