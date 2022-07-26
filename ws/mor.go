package ws

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/midir99/rastreadora/mpp"
	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ParseMorDate(value string) (time.Time, error) {
	date := strings.Split(strings.ToLower(value), " ")
	if len(date) != 3 {
		return time.Time{}, fmt.Errorf("unable to parse date %s", value)
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
		return time.Time{}, fmt.Errorf("unable to parse date %s (unknown month: %s)", value, month)
	}
	day, err := strconv.Atoi(strings.TrimSpace(strings.Replace(date[DAY_INDEX], ",", "", 1)))
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date %s (invalid day number: %s)", value, date[DAY_INDEX])
	}
	year, err := strconv.Atoi(date[YEAR_INDEX])
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date %s (invalid year number: %s)", value, date[YEAR_INDEX])
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
}

func MakeMorAmberUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliamorelos.gob.mx/category/alerta-amber/page/%d/", pageNum)
}

func ScrapeMorAmberPoPosterUrl(pageUrl string) (string, error) {
	doc, err := RetrieveDocument(pageUrl)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve the page %s", pageUrl)
	}
	return AttrOr(Query(doc, "div .post-thumb-img-content img"), "src", ""), nil
}

func ScrapeMorAmberAlerts(doc *html.Node) ([]mpp.MissingPersonPoster, map[int]error) {
	mpps := []mpp.MissingPersonPoster{}
	errs := make(map[int]error)
	for i, article := range QueryAll(doc, "article") {
		mpName := cases.Title(language.LatinAmericanSpanish).String(strings.TrimSpace(Query(article, "h2 a").FirstChild.Data))
		if mpName == "" {
			errs[i+1] = fmt.Errorf("MpName can't be empty")
			continue
		}
		poPostUrl, err := url.Parse(strings.TrimSpace(AttrOr(Query(article, "a"), "href", "")))
		if err != nil {
			errs[i+1] = fmt.Errorf("can't parse PoPostUrl: %s", err)
			continue
		}
		poPostPublicationDate, _ := ParseMorDate(strings.TrimSpace(Query(article, "span .published").FirstChild.Data))
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

func ScrapeMorCustomAlerts(doc *html.Node) ([]mpp.MissingPersonPoster, map[int]error) {
	mpps := []mpp.MissingPersonPoster{}
	errs := make(map[int]error)
	for i, article := range QueryAll(doc, "article") {
		mpName := cases.Title(language.LatinAmericanSpanish).String(strings.TrimSpace(Query(article, "h3 a").FirstChild.Data))
		if mpName == "" {
			errs[i+1] = fmt.Errorf("MpName can't be empty")
			continue
		}
		poPostUrl, err := url.Parse(strings.TrimSpace(AttrOr(Query(article, "h3 a"), "href", "")))
		if err != nil {
			errs[i+1] = fmt.Errorf("can't parse PoPostUrl: %s", err)
			continue
		}
		poPostPublicationDate, _ := ParseMorDate(strings.TrimSpace(Query(article, "span").FirstChild.Data))
		posterUrl := strings.TrimSpace(AttrOr(Query(article, "img"), "src", ""))
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
