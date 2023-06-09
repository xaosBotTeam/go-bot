package task

import (
	"github.com/PuerkitoBio/goquery"
	"go-bot/navigation"
	"go-bot/resources"
	"time"

	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/config"
	"github.com/xaosBotTeam/go-shared-models/status"
)

func NewChests() *Chests {
	return &Chests{lastSync: time.Time{}}
}

type Chests struct {
	lastSync time.Time
}

func (c *Chests) CheckCondition() bool {
	if time.Since(c.lastSync) >= 240*time.Minute {
		c.lastSync = time.Now()
		return true
	}
	return false
}

func (c *Chests) Do(acc account.Account, _ status.Status) error {
	doc, err := navigation.GetPage(acc.URL)

	if err != nil {
		return err
	}
	var isEnd bool
	for {
		doc, err = navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlArenaSkipButton, "открыть сундук")
		doc, isEnd, err = c.fixErrors(doc, err)
		if isEnd || err != nil {
			return err
		}
	}
}

func (c *Chests) fixErrors(doc *goquery.Document, err error) (*goquery.Document, bool, error) {
	if navigation.IsMainPage(doc) {
		doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlBagButton)
		if err != nil {
			return doc, true, err
		}
		return doc, false, nil
	}
	if err == navigation.ErrNotFound {
		return doc, true, nil
	}

	if navigation.IsTopTitleContains(doc, resources.HtmlBagButton) {
		return doc, false, err
	}

	return doc, true, ErrUB
}

func (c *Chests) IsPersistent() bool {
	return true
}

func (c *Chests) RemoveFromStatus(configuration config.Config) config.Config {
	configuration.OpenChests = false
	return configuration
}
