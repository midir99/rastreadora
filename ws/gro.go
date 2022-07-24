package ws

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/midir99/rastreadora/mpp"
	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ParseGroDate(value string) (time.Time, error) {
	content := strings.Split(value, "T")
	if len(content) != 2 {
		return time.Time{}, fmt.Errorf("unable to parse date %s", value)
	}
	date, err := time.Parse("2006-01-02", content[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date %s", value)
	}
	return date, nil
}

func ParseGroFound(value string) bool {
	switch strings.ToLower(value) {
	case "localizada":
		return true
	case "localizado":
		return true
	default:
		return false
	}
}

func ParseGroSex(value string) mpp.Sex {
	switch strings.ToLower(value) {
	case "localizada":
		return mpp.SexFemale
	case "desaparecida":
		return mpp.SexFemale
	case "localizado":
		return mpp.SexMale
	case "desaparecido":
		return mpp.SexMale
	default:
		return mpp.Sex("")
	}
}

func ParseNameSexFound(value string) (string, mpp.Sex, bool) {
	var segments []string
	for _, sep := range []string{";", ":", ",", "."} {
		segments = strings.Split(value, sep)
		if len(segments) == 2 {
			break
		}
	}
	nonEmptySegments := []string{}
	for _, seg := range segments {
		if seg != "" {
			nonEmptySegments = append(nonEmptySegments, seg)
		}
	}
	name := cases.Title(language.LatinAmericanSpanish).String(value)
	found := false
	sex := mpp.Sex("")
	if len(nonEmptySegments) == 2 {
		name = cases.Title(language.LatinAmericanSpanish).String(strings.TrimSpace(nonEmptySegments[1]))
		foundLegend := strings.TrimSpace(nonEmptySegments[0])
		sex = ParseGroSex(foundLegend)
		found = ParseGroFound(foundLegend)
	}
	return name, sex, found
}

func MakeGroAlbaUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliaguerrero.gob.mx/category/alba/page/%d/", pageNum)
}

func ScrapeGroAlbaAlerts(doc *html.Node) []mpp.MissingPersonPoster {
	mpps := []mpp.MissingPersonPoster{}
	for _, article := range QueryAll(doc, ".article_content") {
		foundAndName := strings.TrimSpace(Query(article, "h2 a").FirstChild.Data)
		mpName, _, found := ParseNameSexFound(foundAndName)
		if mpName == "" {
			continue
		}
		poPostUrl, err := url.Parse(strings.TrimSpace(AttrOr(Query(article, "h2 a"), "href", "")))
		if err != nil {
			continue
		}
		pubDate := AttrOr(Query(article, ".entry-date.published"), "datetime", "")
		if pubDate == "" {
			pubDate = AttrOr(Query(article, ".entry-date"), "datetime", "")
		}
		poPostPublicationDate, _ := ParseGroDate(pubDate)
		posterUrl := strings.TrimSpace(AttrOr(Query(article, "a"), "data-src", ""))
		posterUrl = strings.Replace(posterUrl, "-480x320", "", 1)
		poPosterUrl, _ := url.Parse(posterUrl)
		mpps = append(mpps, mpp.MissingPersonPoster{
			AlertType:             mpp.AlertTypeAlba,
			Found:                 found,
			MpName:                mpName,
			MpSex:                 mpp.SexFemale,
			PoPosterUrl:           poPosterUrl,
			PoPostPublicationDate: poPostPublicationDate,
			PoPostUrl:             poPostUrl,
			PoState:               mpp.StateGuerrero,
		})
	}
	return mpps
}

func MakeGroAmberUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliaguerrero.gob.mx/category/amber/page/%d/", pageNum)
}

func ScrapeGroAmberAlerts(doc *html.Node) []mpp.MissingPersonPoster {
	mpps := []mpp.MissingPersonPoster{}
	for _, article := range QueryAll(doc, ".article_content") {
		foundAndName := strings.TrimSpace(Query(article, "h2 a").FirstChild.Data)
		mpName, mpSex, found := ParseNameSexFound(foundAndName)
		if mpName == "" {
			continue
		}
		poPostUrl, err := url.Parse(strings.TrimSpace(AttrOr(Query(article, "h2 a"), "href", "")))
		if err != nil {
			continue
		}
		pubDate := AttrOr(Query(article, ".entry-date.published"), "datetime", "")
		if pubDate == "" {
			pubDate = AttrOr(Query(article, ".entry-date"), "datetime", "")
		}
		poPostPublicationDate, _ := ParseGroDate(pubDate)
		posterUrl := strings.TrimSpace(AttrOr(Query(article, "a"), "data-src", ""))
		posterUrl = strings.Replace(posterUrl, "-480x320", "", 1)
		poPosterUrl, _ := url.Parse(posterUrl)
		mpps = append(mpps, mpp.MissingPersonPoster{
			AlertType:             mpp.AlertTypeAmber,
			Found:                 found,
			MpName:                mpName,
			MpSex:                 mpSex,
			PoPosterUrl:           poPosterUrl,
			PoPostPublicationDate: poPostPublicationDate,
			PoPostUrl:             poPostUrl,
			PoState:               mpp.StateGuerrero,
		})
	}
	return mpps
}

func MakeGroHasVistoAAlertsUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliaguerrero.gob.mx/hasvistoa/?pagina=%d", pageNum)
}

func ScrapeGroHasVistoAAlerts(doc *html.Node) []mpp.MissingPersonPoster {
	mpps := []mpp.MissingPersonPoster{}
	for _, figure := range QueryAll(doc, "figure") {
		h4 := Query(figure, "h4")
		mpName := h4.FirstChild.Data
		missingDate, _ := time.Parse("2006-01-02", h4.LastChild.Data)
		postUrl := AttrOr(Query(figure, "a"), "href", "")
		if postUrl == "" {
			continue
		}
		poPostUrl, err := url.Parse("https://fiscaliaguerrero.gob.mx" + postUrl)
		if err != nil {
			continue
		}
		var poPosterUrl *url.URL
		posterUrl := AttrOr(Query(figure, "img"), "src", "")
		if posterUrl != "" {
			posterUrl = "https://fiscaliaguerrero.gob.mx" + posterUrl
			poPosterUrl, _ = url.Parse(posterUrl)
		}
		mpps = append(mpps, mpp.MissingPersonPoster{
			AlertType:   mpp.AlertTypeHasVistoA,
			MpName:      mpName,
			MissingDate: missingDate,
			PoPosterUrl: poPosterUrl,
			PoPostUrl:   poPostUrl,
			PoState:     mpp.StateGuerrero,
		})
	}
	return mpps
}
