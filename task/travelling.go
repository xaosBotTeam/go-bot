package task

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/config"
	"github.com/xaosBotTeam/go-shared-models/status"
	"go-bot/navigation"
	"go-bot/resources"
	"time"
)

func NewTravelling() *Travelling {
	return &Travelling{
		lastVisit: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

type Travelling struct {
	lastVisit time.Time
}

func (t *Travelling) CheckCondition() bool {
	if time.Since(t.lastVisit) >= 5*time.Minute {
		t.lastVisit = time.Now()
		return true
	}
	return false
}

func (t *Travelling) Do(acc account.Account, _ status.Status) error {
	doc, err := navigation.GetPage(acc.URL)
	if err != nil {
		return err
	}
	var isEnd bool
	for {
		doc, err = navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlMyButtAtt, " Отправиться")
		doc, isEnd, err = t.fixErrors(doc, err)
		if isEnd {
			return err
		}
		doc, err = navigation.SingleFight(doc)
		doc, isEnd, err = t.fixErrors(doc, err)
		if isEnd {
			return err
		}
	}
}

func (t *Travelling) IsPersistent() bool {
	return true
}

func (t *Travelling) GetName() string {
	return "Travelling"
}

func (t *Travelling) RemoveFromStatus(configuration config.Config) config.Config {
	configuration.Travelling = false
	return configuration
}

func (t *Travelling) returnToTravelling(doc *goquery.Document) (*goquery.Document, error) {
	var err error
	if navigation.IsTopTitleContains(doc, resources.HtmlTravellingTitle) {
		return doc, nil
	}
	doc, err = navigation.GoToMainPagePyMenuLink(doc)
	if err != nil {
		return doc, err
	}
	if navigation.IsMainPage(doc) {
		return navigation.GoToFirstMenuLink(doc, resources.HtmlTravellingButton)
	}
	return doc, navigation.ErrNotFound
}

func (t *Travelling) fixErrors(doc *goquery.Document, err error) (*goquery.Document, bool, error) {
	if err == navigation.ErrNotFound && navigation.IsTopTitleContains(doc, resources.HtmlTravellingTitle) {
		return doc, true, nil
	} else if err == navigation.ErrNotFound && navigation.IsMainPage(doc) {
		doc, err = t.returnToTravelling(doc)
		return doc, false, err
	} else if err == navigation.ErrNotFound && navigation.IsVisibleTextContains(doc, " Атаковать ") {
		doc, err = navigation.SingleFight(doc)
		return doc, false, err
	} else if err != nil {
		return doc, true, err
	}
	return doc, false, err
}
