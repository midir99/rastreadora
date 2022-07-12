package ws

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/midir99/rastreadora/mpp"
)

func RetrieveDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func Scrape(pageUrl string, scraper func(*goquery.Document) []mpp.MissingPersonPoster, ch chan []mpp.MissingPersonPoster) {
	doc, err := RetrieveDocument(pageUrl)
	if err != nil {
		ch <- []mpp.MissingPersonPoster{}
		return
	}
	ch <- scraper(doc)
}
