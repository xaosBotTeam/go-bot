package task

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/xaosBotTeam/go-shared-models/account"
	models "github.com/xaosBotTeam/go-shared-models/task"
	"go-bot/navigation"
	"go-bot/random"
	"go-bot/resources"
	"strconv"
	"time"
)

var _ Abstract = (*ArenaBoosting)(nil)

func NewArenaBoosting(status *models.Status) *ArenaBoosting {
	if status.ArenaFarming {
		return &ArenaBoosting{
			UseEnergyCans: status.ArenaUseEnergyCans,
		}
	}
	return nil
}

type ArenaBoosting struct {
	UseEnergyCans bool
}

func (t *ArenaBoosting) Do(acc account.Account) error {
	doc, err := navigation.GetPage(acc.URL)

	if err != nil {
		return err
	}
	doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
	if err == navigation.ErrNotFound {
		doc, err = navigation.GoToMainPagePyMenuLink(doc)
		if err != nil {
			return err
		}
	} else if navigation.IsVisibleTextContains(doc, resources.HtmlYouAreTooWeak) {
		doc, err = navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlRestoreHealth, "Восстановить жизни")

		if err != nil {
			return err
		}
		if navigation.IsMainPage(doc) {
			doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
			if err != nil {
				return err
			}
		}
	}

	energyEnough := true
	for energyEnough {
		for {
			doc, err = navigation.GoByClass(doc, resources.HtmlArenaGoldButton)

			if navigation.IsVisibleTextContains(doc, resources.HtmlEnergyIsEmpty) {
				break
			} else if navigation.IsVisibleTextContains(doc, resources.HtmlYouAreTooWeak) {
				doc, err = navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlRestoreHealth, "Восстановить жизни")
				if err != nil {
					return err
				}
				if navigation.IsMainPage(doc) {
					doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
					if err != nil {
						return err
					}
				}
			} else if err == navigation.ErrNotFound && !navigation.IsMainPage(doc) {
				break
			} else if err == navigation.ErrNotFound && navigation.IsMainPage(doc) {
				doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		}

		doc, err = navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlArenaSkipButton, resources.HtmlAnotherRival)

		if err != nil && !navigation.IsVisibleTextContains(doc, resources.HtmlEnergyIsEmpty) && !navigation.IsMainPage(doc) {
			return err
		}
		if navigation.IsMainPage(doc) {
			doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
		}
		if err != nil {
			return err
		}
		doc, energyEnough, err = t.restoreEnergy(doc, acc.EnergyLimit)
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

//func (t *ArenaBoosting) tryToReturnToMainProcess(doc *goquery.Document) (*goquery.Document, bool, error) {
//	if navigation.IsVisibleTextContains(doc, resources.HtmlEnergyIsEmpty) {
//		return doc, false, nil
//	} else if navigation.IsVisibleTextContains(doc, resources.HtmlYouAreTooWeak) {
//		doc, err := navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlRestoreHealth, "Восстановить жизни")
//		if err != nil {
//			return doc, false, err
//		}
//		doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
//		if err != nil {
//			return doc, false, err
//		}
//		return doc, true, nil
//	} else if err == navigation.ErrNotFound && !navigation.IsMainPage(doc) {
//		break
//	} else if err == navigation.ErrNotFound && navigation.IsMainPage(doc) {
//		doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
//		if err != nil {
//			return err
//		}
//	} else if err != nil {
//		return err
//	}
//}

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

func (t *ArenaBoosting) RemoveFromStatus(status models.Status) models.Status {
	status.ArenaFarming = false
	return status
}
