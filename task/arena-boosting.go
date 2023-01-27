package task

import (
	"github.com/PuerkitoBio/goquery"
	account "github.com/xaosBotTeam/go-shared-models/dbAccountInformation"
	"go-bot/navigation"
	"go-bot/random"
	"go-bot/resources"
	"net/http"
	"strconv"
	"time"
)

var _ Abstract = (*ArenaBoosting)(nil)

type ArenaBoosting struct {
	UseEnergyCans bool
}

func (t *ArenaBoosting) Do(acc account.DbAccountInformation) error {
	response, err := http.Get(acc.URL)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
		if err == navigation.ErrNotFound {
			time.Sleep(random.RandomWaitTime())
			doc, err = navigation.GoToMainPagePyMenuLink(doc)
			if err != nil {
				return err
			}
			time.Sleep(random.RandomWaitTime())
		} else if navigation.IsVisibleTextContains(doc, resources.HtmlYouAreTooWeak) {
			doc, err = navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlRestoreHealth, "Восстановить жизни")
			if err != nil {
				return err
			}
			doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
			if err != nil {
				return err
			}
			break
		} else {
			break
		}
	}
	time.Sleep(random.RandomWaitTime())
	energyEnough := true
	for energyEnough {
		for {
			doc, err = navigation.GoByClass(doc, resources.HtmlArenaGoldButton)

			time.Sleep(random.RandomWaitTime())
			if navigation.IsVisibleTextContains(doc, resources.HtmlEnergyIsEmpty) {
				break
			} else if navigation.IsVisibleTextContains(doc, resources.HtmlYouAreTooWeak) {
				doc, err = navigation.GoByClassAndVisibleTextContains(doc, resources.HtmlRestoreHealth, "Восстановить жизни")
				if err != nil {
					return err
				}
				doc, err = navigation.GoToFirstMenuLink(doc, resources.HtmlDeathArena)
				if err != nil {
					return err
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

		doc, err = navigation.GoByClass(doc, resources.HtmlArenaSkipButton)
		time.Sleep(random.RandomWaitTime())

		if err != nil && !navigation.IsVisibleTextContains(doc, resources.HtmlEnergyIsEmpty) {
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
