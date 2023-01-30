package collector

import (
	"github.com/xaosBotTeam/go-shared-models/account"
	"go-bot/navigation"
	"go-bot/resources"
	"strings"
	"time"
)

func NewNickname() *Nickname {
	return &Nickname{
		lastSync: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

type Nickname struct {
	lastSync time.Time
}

func (n *Nickname) Collect(acc account.Account) (account.Account, error) {
	doc, err := navigation.GetPage(acc.URL)
	if err != nil {
		return account.Account{}, err
	}
	doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlCharacter)
	if err != nil {
		return account.Account{}, err
	}
	acc.FriendlyName = strings.TrimSpace(strings.Split(navigation.GetTopTitle(doc), ",")[0])
	return acc, nil
}
func (n *Nickname) CheckCondition() bool {
	if time.Now().Sub(n.lastSync) > 24*time.Hour {
		n.lastSync = time.Now()
		return true
	}
	return false
}
