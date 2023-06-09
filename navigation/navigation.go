package navigation

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	httpbridg "go-bot/http-bridge"
	"go-bot/resources"
	"strings"
)

var ErrNotFound = errors.New("item not found")
var ErrEnergyEmpty = errors.New("energy is over")

func GetPage(url string) (*goquery.Document, error) {
	rsp, err := httpbridg.GetBodyBytes(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(rsp)
}

func IsMainPage(doc *goquery.Document) bool {
	selection := doc.Find("." + resources.HtmlTopTitle)
	return selection.Text() == "Наследие Хаоса"
}

func GoToMainPagePyMenuLink(doc *goquery.Document) (*goquery.Document, error) {
	return GoToFirstMenuLink(doc, resources.HtmlMainPage)
}

func GoToFirstMenuLink(doc *goquery.Document, label string) (*goquery.Document, error) {
	newDoc := doc
	finalError := ErrNotFound
	doc.Find(".menu_link").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		debugStr := s.Text()
		if strings.Contains(debugStr, label) {
			url, ok := s.Attr("href")
			if ok {
				rsp, err := httpbridg.GetBodyBytes(resources.UrlPrefix + url)

				if err != nil {
					finalError = err
					return false
				}
				doc, err := goquery.NewDocumentFromReader(rsp)
				if err != nil {
					finalError = err
					return false
				}
				newDoc = doc
				finalError = nil
			}
			return false
		}
		return true

	})
	return newDoc, finalError
}

func GoByClass(doc *goquery.Document, class string) (*goquery.Document, error) {
	url, ok := doc.Find("." + class).First().Attr("href")
	if ok {
		rsp, err := httpbridg.GetBodyBytes(resources.UrlPrefix + url)

		if err != nil {
			return nil, err
		}
		doc, err := goquery.NewDocumentFromReader(rsp)
		if err != nil {
			return nil, err
		}
		return doc, nil
	}
	return doc, ErrNotFound
}

func GetTopBarValue(doc *goquery.Document, index int) (string, error) {
	res := strings.Fields(doc.Find("." + resources.HtmlStatusBar).First().Text())
	if len(res) < index {
		return "", errors.New("index is out of result array length")
	}
	return res[index], nil
}

func IsVisibleTextContains(doc *goquery.Document, text string) bool {
	return strings.Contains(doc.Text(), text)
}

func GoByClassAndVisibleTextContains(doc *goquery.Document, class string, text string) (*goquery.Document, error) {
	newDoc := doc
	finalError := ErrNotFound
	doc.Find("." + class).EachWithBreak(func(_ int, s *goquery.Selection) bool {
		debugStr := s.Text()
		if strings.Contains(debugStr, text) {
			url, ok := s.Attr("href")
			if ok {
				rsp, err := httpbridg.GetBodyBytes(resources.UrlPrefix + url)
				if err != nil {
					finalError = err
					return false
				}
				doc, err := goquery.NewDocumentFromReader(rsp)
				if err != nil {
					finalError = err
					return false
				}
				newDoc = doc
				finalError = nil
			}
			return false
		}
		return true

	})
	return newDoc, finalError
}

func IsTopTitleContains(doc *goquery.Document, title string) bool {
	selection := doc.Find("." + resources.HtmlTopTitle)
	return strings.Contains(selection.Text(), title)
}

func GetTopTitle(document *goquery.Document) string {
	return document.Find("." + resources.HtmlTopTitle).Text()
}

func ValidateUrl(url string) bool {
	body, err := httpbridg.GetBodyBytes(url)
	if err != nil {
		return false
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return false
	}
	doc, err = GoToMainPagePyMenuLink(doc)
	if err != nil {
		return false
	}
	return IsMainPage(doc)
}

func SingleFight(doc *goquery.Document) (*goquery.Document, error) {
	var err error
	for {
		doc, err = GoByClassAndVisibleTextContains(doc, resources.HtmlMyButtAtt, "Атаковать")
		if err == ErrNotFound {
			return doc, nil
		} else if err != nil {
			return doc, err
		} else if IsVisibleTextContains(doc, resources.HtmlEnergyIsEmpty) {
			return doc, ErrEnergyEmpty
		}
	}
}
