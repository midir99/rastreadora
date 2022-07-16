package ws

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/midir99/rastreadora/mpp"
)

func ParseMorDate(value string) (string, error) {
	date := strings.Split(value, " ")
	if len(date) != 3 {
		return "", fmt.Errorf("unable to parse date %s", value)
	}
	MONTH_INDEX, DAY_INDEX, YEAR_INDEX := 0, 1, 2
	var month time.Month
	switch date[MONTH_INDEX] {
	case "enero":
		month = time.January
	case "febrero":
		month = time.February
	case "marzo":
		month = time.March
	case "abril":
		month = time.April
	case "mayo":
		month = time.May
	case "junio":
		month = time.June
	case "julio":
		month = time.July
	case "agosto":
		month = time.August
	case "septiembre":
		month = time.September
	case "octubre":
		month = time.October
	case "noviembre":
		month = time.November
	case "diciembre":
		month = time.December
	default:
		return "", fmt.Errorf("unable to parse date %s", value)
	}
	day, err := strconv.Atoi(date[DAY_INDEX])
	if err != nil {
		return "", fmt.Errorf("unable to parse date %s", value)
	}
	year, err := strconv.Atoi(date[YEAR_INDEX])
	if err != nil {
		return "", fmt.Errorf("unable to parse date %s", value)
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02"), nil
}

func MakeMorAmberUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliamorelos.gob.mx/category/alerta-amber/page/%d/", pageNum)
}

func ScrapeMorAmberPoPostUrl(pageUrl string, client *http.Client) (string, error) {
	doc, err := RetrieveDocument(pageUrl, client)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve the page %s", pageUrl)
	}
	return doc.Find("div .post-thumb-img-content img").AttrOr("src", ""), nil
}

func ScrapeMorAmberAlerts(doc *goquery.Document, client *http.Client) []mpp.MissingPersonPoster {
	mpps := []mpp.MissingPersonPoster{}
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		mpName := strings.Title(strings.TrimSpace(s.Find("h2").Text()))
		if mpName == "" {
			return
		}
		poPostUrl := strings.TrimSpace(s.Find("a").AttrOr("href", ""))
		if poPostUrl == "" {
			return
		}
		poPostPublicationDate, _ := ParseMorDate(strings.TrimSpace(s.Find("span .published").Text()))
		poPosterUrl, _ := ScrapeMorAmberPoPostUrl(poPostUrl, client)
		mpps = append(mpps, mpp.MissingPersonPoster{
			AlertType:             mpp.AlertTypeAmber,
			MpName:                mpName,
			PoPosterUrl:           poPosterUrl,
			PoPostPublicationDate: poPostPublicationDate,
			PoPostUrl:             poPostUrl,
			PoState:               mpp.StateMorelos,
		})
	})
	return mpps
}

func MakeMorCustomUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliamorelos.gob.mx/cedulas/%d/", pageNum)
}

func ScrapeMorCustomAlerts(doc *goquery.Document, client *http.Client) []mpp.MissingPersonPoster {
	mpps := []mpp.MissingPersonPoster{}
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		mpName := strings.Title(strings.TrimSpace(s.Find("h3 a").Text()))
		if mpName == "" {
			return
		}
		poPostUrl := strings.TrimSpace(s.Find("h3 a").AttrOr("href", ""))
		if poPostUrl == "" {
			return
		}
		poPostPublicationDate, _ := ParseMorDate(strings.TrimSpace(s.Find("span").Text()))
		poPosterUrl := strings.TrimSpace(s.Find("img").AttrOr("src", ""))
		mpps = append(mpps, mpp.MissingPersonPoster{
			MpName:                mpName,
			PoPosterUrl:           poPosterUrl,
			PoPostPublicationDate: poPostPublicationDate,
			PoPostUrl:             poPostUrl,
			PoState:               mpp.StateMorelos,
		})
	})
	return mpps
}
