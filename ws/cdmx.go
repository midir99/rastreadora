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

func ParseCdmxDate(value string) (time.Time, error) {
	date := strings.Split(strings.ToLower(value), " ")
	if len(date) != 5 {
		return time.Time{}, fmt.Errorf("unable to parse date %s", value)
	}
	DAY, MONTH, YEAR := 0, 2, 4
	day, err := strconv.ParseUint(date[DAY], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date %s (invalid day number: %s)", value, date[DAY])
	}
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
	year, err := strconv.ParseUint(date[YEAR], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date %s (invalid year number: %s)", value, date[YEAR])
	}
	return time.Date(int(year), month, int(day), 0, 0, 0, 0, time.UTC), nil
}

func ParseCdmxFound(value string) bool {
	switch strings.ToLower(value) {
	case "localizado":
		return true
	case "no localizado":
		return false
	case "ausente":
		return false
	default:
		return false
	}
}

func ParseCdmxAge(value string) (int, error) {
	age := strings.Split(value, " ")
	YEARS := 0
	years, err := strconv.ParseUint(age[YEARS], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse age %s", value)
	}
	return int(years), nil
}

func MakeCdmxCustomUrl(pageNum uint64) string {
	return fmt.Sprintf("https://personasdesaparecidas.fgjcdmx.gob.mx/listado.php?pa=%d&re=100", pageNum)
}

func ScrapeCdmxCustomAlerts(d *doc.Doc) ([]mpp.MissingPersonPoster, map[int]error) {
	mpps := []mpp.MissingPersonPoster{}
	errs := make(map[int]error)
	for i, tr := range d.QueryAll("tbody tr") {
		tds := tr.QueryAll("td")
		if len(tds) != 2 {
			errs[i+1] = fmt.Errorf("entry only has not 2 td elements")
			continue
		}
		posterTd := tds[0]
		dataTd := tds[1]
		mpName := cases.Title(language.LatinAmericanSpanish).String(strings.ReplaceAll(strings.TrimSpace(dataTd.NthChild(0).Text()), "\u00A0", " "))
		if mpName == "" {
			errs[i+1] = fmt.Errorf("MpName can't be empty")
			continue
		}
		postUrl := strings.TrimSpace(dataTd.NthChild(10).AttrOr("href", ""))
		if postUrl == "" {
			errs[i+1] = fmt.Errorf("PoPostUrl can't be empty")
			continue
		}
		poPostUrl, err := url.Parse("https://personasdesaparecidas.fgjcdmx.gob.mx/" + postUrl)
		if err != nil {
			errs[i+1] = fmt.Errorf("can't parse PoPostUrl: %s", err)
			continue
		}
		var poPosterUrl *url.URL
		posterUrl := strings.TrimSpace(posterTd.Query("img").AttrOr("src", ""))
		if posterUrl != "" {
			posterUrl = "https://personasdesaparecidas.fgjcdmx.gob.mx/" + posterUrl
			poPosterUrl, _ = url.Parse(posterUrl)
		}
		var missingDate time.Time
		missingDateLegend := strings.Split(dataTd.NthChild(4).Text(), ":\u00A0")
		if len(missingDateLegend) == 2 {
			missingDate, _ = ParseCdmxDate(missingDateLegend[1])
		}
		var found bool
		foundLegend := strings.Split(dataTd.NthChild(8).Text(), ":\u00A0")
		if len(foundLegend) == 2 {
			found = ParseCdmxFound(foundLegend[1])
		}
		var age int
		ageLegend := strings.Split(dataTd.NthChild(2).Text(), ":\u00A0")
		if len(ageLegend) == 2 {
			age, _ = ParseCdmxAge(ageLegend[1])
		}
		mpps = append(mpps, mpp.MissingPersonPoster{
			Found:                found,
			MissingDate:          missingDate,
			MpAgeWhenDisappeared: age,
			MpName:               mpName,
			PoPosterUrl:          poPosterUrl,
			PoPostUrl:            poPostUrl,
			PoState:              mpp.StateCiudadDeMexico,
		})
	}
	return mpps, errs
}
