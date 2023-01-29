package info_collector

import (
	"github.com/xaosBotTeam/go-shared-models/account"
	"go-bot/navigation"
	"go-bot/resources"
	"strconv"
	"strings"
	"time"
)

func NewGameId() *GameId {
	return &GameId{
		lastSync: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

type GameId struct {
	lastSync time.Time
}

func (o *GameId) Collect(acc account.Account) (account.Account, error) {
	doc, err := navigation.GetPage(acc.URL)
	if err != nil {
		return account.Account{}, err
	}
	doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlCharacter)
	if err != nil {
		return account.Account{}, err
	}
	id := 0
	words := strings.Fields(doc.Text())
	for i, word := range words {
		if word == "ID" {
			if len(words) > i+2 && words[i+1] == "игрока:" {
				id, err = strconv.Atoi(words[i+2])
				if err != nil {
					return account.Account{}, err
				}
			}
		}
	}

	acc.GameID = id
	return acc, nil
}
func (g *GameId) CheckCondition() bool {
	if time.Now().Sub(g.lastSync) > 24*time.Hour {
		g.lastSync = time.Now()
		return true
	}
	return false
}
