package ws

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/midir99/rastreadora/mpp"
)

func MakeClient(skipCert bool) *http.Client {
	if skipCert {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr}
	} else {
		return http.DefaultClient
	}
}

func RetrieveDocument(url string, client *http.Client) (*goquery.Document, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func Scrape(pageUrl string, client *http.Client, scraper func(*goquery.Document, *http.Client) []mpp.MissingPersonPoster, ch chan []mpp.MissingPersonPoster) {
	doc, err := RetrieveDocument(pageUrl, client)
	if err != nil {
		log.Printf("Error: %s\n", err)
		ch <- []mpp.MissingPersonPoster{}
		return
	}
	ch <- scraper(doc, client)
}
