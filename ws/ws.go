package ws

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/andybalholm/cascadia"
	"github.com/midir99/rastreadora/mpp"
	"golang.org/x/net/html"
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

func RetrieveDocument(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with a %d status code", url, resp.StatusCode)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func Query(n *html.Node, query string) *html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return &html.Node{}
	}
	return cascadia.Query(n, sel)
}

func QueryAll(n *html.Node, query string) []*html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return []*html.Node{}
	}
	return cascadia.QueryAll(n, sel)
}

func AttrOr(n *html.Node, attr, or string) string {
	for _, a := range n.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return or
}

func Scrape(pageUrl string, scraper func(*html.Node) []mpp.MissingPersonPoster, ch chan []mpp.MissingPersonPoster) {
	doc, err := RetrieveDocument(pageUrl)
	if err != nil {
		log.Printf("Error: %s\n", err)
		ch <- []mpp.MissingPersonPoster{}
		return
	}
	ch <- scraper(doc)
}
