package ws

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/midir99/rastreadora/doc"
	"github.com/midir99/rastreadora/mpp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ParseMorDate(value string) (time.Time, error) {
	date := strings.Split(strings.ToLower(value), " ")
	if len(date) != 3 {
		return time.Time{}, fmt.Errorf("unable to parse date %s", value)
	}
	MONTH, DAY, YEAR := 0, 1, 2
	var month time.Month
	switch date[MONTH] {
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
		return time.Time{}, fmt.Errorf("unable to parse date %s (invalid month: %s)", value, month)
	}
	day, err := strconv.ParseUint(strings.TrimSpace(strings.Replace(date[DAY], ",", "", 1)), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date %s (invalid day number: %s)", value, date[DAY])
	}
	year, err := strconv.ParseUint(date[YEAR], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date %s (invalid year number: %s)", value, date[YEAR])
	}
	return time.Date(int(year), month, int(day), 0, 0, 0, 0, time.UTC), nil
}

func MakeMorAmberUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliamorelos.gob.mx/category/alerta-amber/page/%d/", pageNum)
}

func ScrapeMorAmberPoPosterUrl(pageUrl string) (string, error) {
	doc, err := RetrieveDocument(pageUrl, false)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve the page %s", pageUrl)
	}
	return doc.Query("div .post-thumb-img-content img").AttrOr("src", ""), nil
}

func ScrapeMorAmberAlerts(d *doc.Doc) ([]mpp.MissingPersonPoster, map[int]error) {
	mpps := []mpp.MissingPersonPoster{}
	errs := make(map[int]error)
	for i, article := range d.QueryAll("article") {
		mpName := cases.Title(language.LatinAmericanSpanish).String(strings.TrimSpace(article.Query("h2 a").Text()))
		if mpName == "" {
			errs[i+1] = fmt.Errorf("MpName can't be empty")
			continue
		}
		poPostUrl, err := url.Parse(strings.TrimSpace(article.Query("a").AttrOr("href", "")))
		if err != nil {
			errs[i+1] = fmt.Errorf("can't parse PoPostUrl: %s", err)
			continue
		}
		poPostPublicationDate, _ := ParseMorDate(strings.TrimSpace(article.Query("span .published").Text()))
		posterUrl, _ := ScrapeMorAmberPoPosterUrl(poPostUrl.String())
		poPosterUrl, _ := url.Parse(posterUrl)
		mpps = append(mpps, mpp.MissingPersonPoster{
			AlertType:             mpp.AlertTypeAmber,
			MpName:                mpName,
			PoPosterUrl:           poPosterUrl,
			PoPostPublicationDate: poPostPublicationDate,
			PoPostUrl:             poPostUrl,
			PoState:               mpp.StateMorelos,
		})
	}
	return mpps, errs
}

func MakeMorCustomUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliamorelos.gob.mx/cedulas/%d/", pageNum)
}

func ScrapeMorCustomAlerts(d *doc.Doc) ([]mpp.MissingPersonPoster, map[int]error) {
	mpps := []mpp.MissingPersonPoster{}
	errs := make(map[int]error)
	for i, article := range d.QueryAll("article") {
		mpName := cases.Title(language.LatinAmericanSpanish).String(strings.TrimSpace(article.Query("h3 a").Text()))
		if mpName == "" {
			errs[i+1] = fmt.Errorf("MpName can't be empty")
			continue
		}
		poPostUrl, err := url.Parse(strings.TrimSpace(article.Query("h3 a").AttrOr("href", "")))
		if err != nil {
			errs[i+1] = fmt.Errorf("can't parse PoPostUrl: %s", err)
			continue
		}
		poPostPublicationDate, _ := ParseMorDate(strings.TrimSpace(article.Query("span").Text()))
		posterUrl := strings.TrimSpace(article.Query("img").AttrOr("src", ""))
		posterUrl = strings.Replace(posterUrl, "-300x225", "", 1)
		posterUrl = strings.Replace(posterUrl, "-300x240", "", 1)
		poPosterUrl, _ := url.Parse(posterUrl)
		mpps = append(mpps, mpp.MissingPersonPoster{
			MpName:                mpName,
			PoPosterUrl:           poPosterUrl,
			PoPostPublicationDate: poPostPublicationDate,
			PoPostUrl:             poPostUrl,
			PoState:               mpp.StateMorelos,
		})
	}
	return mpps, errs
}
