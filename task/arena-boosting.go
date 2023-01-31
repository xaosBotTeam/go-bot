package task

import (
	"go-bot/navigation"
	"go-bot/random"
	"go-bot/resources"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/config"
	"github.com/xaosBotTeam/go-shared-models/status"
)

var _ Abstract = (*ArenaBoosting)(nil)

func NewArenaBoosting(configuration config.Config) *ArenaBoosting {
	if configuration.ArenaFarming {
		return &ArenaBoosting{
			UseEnergyCans: configuration.ArenaUseEnergyCans,
		}
	}
	return nil
}

type ArenaBoosting struct {
	UseEnergyCans bool
}

func (t *ArenaBoosting) Do(acc account.Account, s status.Status) error {
	doc, err := navigation.GetPage(acc.URL)
	if err != nil {
		return err
	}

	doc, err = t.returnToArena(doc)

	energyEnough := true
	for energyEnough {
		doc, err = navigation.GoByClass(doc, resources.HtmlArenaGoldButton)
		if navigation.IsVisibleTextContains(doc, resources.HtmlEnergyIsEmpty) {
			doc, energyEnough, err = t.restoreEnergy(doc, s.EnergyLimit)
		} else if err == navigation.ErrNotFound && navigation.IsTopTitleContains(doc, "Арена Смерти") {
			doc, err = navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlArenaSkipButton, resources.HtmlAnotherRival)
			if err != nil {
				return err
			}
			doc, energyEnough, err = t.restoreEnergy(doc, s.EnergyLimit)
			if err != nil {
				return err
			}
		}
		doc, err = t.returnToArena(doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ArenaBoosting) CheckCondition() bool {
	return true
}

func (t *ArenaBoosting) GetName() string {
	return "ArenaBoosting"
}

func (t *ArenaBoosting) restoreEnergy(doc *goquery.Document, limit int) (*goquery.Document, bool, error) {
	err := (error)(nil)
	if t.UseEnergyCans {
		for {
			doc, err = navigation.GoToFirstMenuLink(doc, "Энергия")
			// energy is over
			if err == navigation.ErrNotFound && !navigation.IsMainPage(doc) {
				return doc, false, nil
			} else if err == navigation.ErrNotFound && navigation.IsMainPage(doc) {
				doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
				if err != nil {
					return doc, false, err
				}
			} else if err != nil {
				return doc, false, err
			}
			energyStr, err := navigation.GetTopBarValue(doc, 1)
			if err != nil {
				return doc, false, err
			}
			energy, err := strconv.Atoi(energyStr)
			if err != nil {
				return doc, false, err
			}
			if energy >= limit {
				return doc, true, nil
			}
			time.Sleep(random.RandomWaitTime())
		}
	} else {
		return doc, false, nil
	}
}

func (t *ArenaBoosting) IsPersistent() bool {
	return false
}

func (t *ArenaBoosting) RemoveFromStatus(configuration config.Config) config.Config {
	configuration.ArenaFarming = false
	return configuration
}

func (t *ArenaBoosting) restoreCharacterHealthReturnToArena(doc *goquery.Document) (*goquery.Document, error) {
	doc, err := navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlRestoreHealth, "Восстановить жизни")
	if err != nil {
		return doc, err
	}
	if navigation.IsMainPage(doc) {
		doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
		if err != nil {
			return doc, err
		}
	}
	return doc, err
}

func (t *ArenaBoosting) returnToArena(doc *goquery.Document) (*goquery.Document, error) {
	var err error
	if navigation.IsVisibleTextContains(doc, resources.HtmlYouAreTooWeak) {
		doc, err = t.restoreCharacterHealthReturnToArena(doc)
	}

	if err != nil {
		return doc, err
	}

	if navigation.IsTopTitleContains(doc, "Арена Смерти") {
		return doc, nil
	}
	if navigation.IsMainPage(doc) {
		return navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
	}
	return doc, navigation.ErrNotFound
}
