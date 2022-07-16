package ws

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/midir99/rastreadora/mpp"
)

func ParseGroDate(value string) (string, error) {
	content := strings.Split(value, "T")
	if len(content) != 2 {
		return "", fmt.Errorf("unable to parse date %s", value)
	}
	date, err := time.Parse("2006-01-02", content[0])
	if err != nil {
		return "", fmt.Errorf("unable to parse date %s", value)
	}
	return date.Format("2006-01-02"), nil
}

func ParseGroFound(value string) bool {
	switch value {
	case "localizada":
		return true
	case "localizado":
		return true
	default:
		return false
	}
}

func ParseGroSex(value string) mpp.Sex {
	switch value {
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

func MakeGroAlbaUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliaguerrero.gob.mx/category/alba/page/%d/", pageNum)
}

func ScrapeGroAlbaAlerts(doc *goquery.Document, client *http.Client) []mpp.MissingPersonPoster {
	mpps := []mpp.MissingPersonPoster{}
	doc.Find(".article_content").Each(func(i int, s *goquery.Selection) {
		foundAndName := strings.TrimSpace(s.Find("h2 a").Text())
		seps := []string{";", ":", ",", "."}
		var content []string
		for _, sep := range seps {
			content = strings.Split(foundAndName, sep)
			if len(content) == 2 {
				break
			}
		}
		mpName := strings.Title(foundAndName)
		found := false
		if len(content) == 2 {
			mpName = strings.Title(strings.TrimSpace(content[1]))
			foundLegend := strings.ToLower(strings.TrimSpace(content[0]))
			found = ParseGroFound(foundLegend)
		}
		if mpName == "" {
			return
		}
		poPostUrl := strings.TrimSpace(s.Find("h2 a").AttrOr("href", ""))
		if poPostUrl == "" {
			return
		}
		poPostPublicationDate := s.Find(".entry-date.published").First().AttrOr("datetime", "")
		if poPostPublicationDate == "" {
			poPostPublicationDate = s.Find(".entry-date").First().AttrOr("datetime", "")
		}
		poPostPublicationDate, _ = ParseGroDate(poPostPublicationDate)
		poPosterUrl := strings.TrimSpace(s.Find("a").AttrOr("data-src", ""))
		poPosterUrl = strings.Replace(poPosterUrl, "-480x320", "", 1)
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
	})
	return mpps
}

func MakeGroAmberUrl(pageNum uint64) string {
	return fmt.Sprintf("https://fiscaliaguerrero.gob.mx/category/amber/page/%d/", pageNum)
}

func ScrapeGroAmberAlerts(doc *goquery.Document, client *http.Client) []mpp.MissingPersonPoster {
	mpps := []mpp.MissingPersonPoster{}
	doc.Find(".article_content").Each(func(i int, s *goquery.Selection) {
		foundAndName := strings.TrimSpace(s.Find("h2 a").Text())
		seps := []string{";", ":", ",", "."}
		var content []string
		for _, sep := range seps {
			content = strings.Split(foundAndName, sep)
			if len(content) == 2 {
				break
			}
		}
		mpName := strings.Title(foundAndName)
		mpSex := mpp.Sex("")
		found := false
		if len(content) == 2 {
			mpName = strings.Title(strings.TrimSpace(content[1]))
			foundLegend := strings.ToLower(strings.TrimSpace(content[0]))
			mpSex = ParseGroSex(foundLegend)
			found = ParseGroFound(foundLegend)
		}
		if mpName == "" {
			return
		}
		poPostUrl := strings.TrimSpace(s.Find("h2 a").AttrOr("href", ""))
		if poPostUrl == "" {
			return
		}
		poPostPublicationDate := s.Find(".entry-date.published").First().AttrOr("datetime", "")
		if poPostPublicationDate == "" {
			poPostPublicationDate = s.Find(".entry-date").First().AttrOr("datetime", "")
		}
		poPostPublicationDate, _ = ParseGroDate(poPostPublicationDate)
		poPosterUrl := strings.TrimSpace(s.Find("a").AttrOr("data-src", ""))
		poPosterUrl = strings.Replace(poPosterUrl, "-480x320", "", 1)
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
	})
	return mpps
}
