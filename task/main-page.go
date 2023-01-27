package task

import (
	"github.com/PuerkitoBio/goquery"
	account "github.com/xaosBotTeam/go-shared-models/dbAccountInformation"
	"go-bot/navigation"
	"net/http"
)

var _ Abstract = (*GoMainPage)(nil)

type GoMainPage struct {
}

func (t *GoMainPage) Do(acc account.DbAccountInformation) error {
	response, err := http.Get(acc.URL)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}

	if navigation.IsMainPage(doc) {
		return nil
	}

	_, err = navigation.GoToMainPagePyMenuLink(doc)

	return err
}

func (t *GoMainPage) CheckCondition() bool {
	return true
}

func (t *GoMainPage) GetName() string {
	return "GoToMainPage"
}
