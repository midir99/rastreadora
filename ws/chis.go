package ws

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/midir99/rastreadora/doc"
	"github.com/midir99/rastreadora/mpp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ParseChisBuild(value string) mpp.PhysicalBuild {
	switch strings.ToLower(value) {
	case "atletica":
		return mpp.PhysicalBuildSlim
	case "delgada":
		return mpp.PhysicalBuildSlim
	case "mediana":
		return mpp.PhysicalBuildRegular
	case "no especificado":
		return mpp.PhysicalBuild("")
	case "obesa":
		return mpp.PhysicalBuildHeavy
	case "regular":
		return mpp.PhysicalBuildRegular
	case "robusta":
		return mpp.PhysicalBuildHeavy
	case "sin dato":
		return mpp.PhysicalBuild("")
	default:
		return mpp.PhysicalBuild("")
	}
}

func ParseChisComplexion(value string) mpp.Complexion {
	switch strings.ToLower(value) {
	case "albino":
		return mpp.ComplexionVeryLight
	case "api\u00F1onado":
		return mpp.ComplexionLightIntermediate
	case "blanca":
		return mpp.ComplexionLight
	case "morena":
		return mpp.ComplexionDark
	case "morena clara":
		return mpp.ComplexionDarkIntermediate
	case "morena obscura":
		return mpp.ComplexionVeryDark
	default:
		return mpp.Complexion("")
	}
}

func ParseChisDate(value string) (time.Time, error) {
	date, err := time.Parse("02/01/2006", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date %s", value)
	}
	return date, nil
}

func ParseChisFound(value string) bool {
	switch strings.ToLower(value) {
	case "ausente":
		return false
	case "desaparecida":
		return false
	case "extraviada":
		return false
	case "no localizada":
		return false
	default:
		return false
	}
}

func ParseChisHeigth(value string) (int, error) {
	value = strings.ReplaceAll(value, "C", "")
	value = strings.TrimSpace(value)
	meters, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse heigth %s", value)
	}
	centimeters := meters * 100
	ratio := math.Pow(10, float64(2)) // round to 2 decimal places
	rounded := math.Round(centimeters*ratio) / ratio
	return int(rounded), nil
}

func ParseChisSex(value string) mpp.Sex {
	switch strings.ToLower(value) {
	case "hombre":
		return mpp.SexMale
	case "mujer":
		return mpp.SexFemale
	default:
		return mpp.Sex("")
	}
}

func ParseChisWeight(value string) (int, error) {
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, "kg", "")
	value = strings.ReplaceAll(value, ".", "")
	value = strings.TrimSpace(value)
	kgs, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse %s", value)
	}
	return int(kgs), nil
}

func MakeChisHasVistoAUrl(pageNum uint64) string {
	if pageNum > 0 {
		pageNum--
	}
	return fmt.Sprintf("https://www.fge.chiapas.gob.mx/Servicios/Hasvistoa/Page/%d", pageNum)
}

func ScrapeChisHasVistoAExtraData(pageUrl string) (*mpp.MissingPersonPoster, error) {
	doc, err := RetrieveDocument(pageUrl, false)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the page %s", pageUrl)
	}
	var (
		Sex            = 0
		Heigth         = 1
		Complexion     = 2
		Eyes           = 3
		Hair           = 4
		Weight         = 5
		MissingDate    = 6
		Build          = 7
		Mouth          = 8
		NoseSize       = 9
		NoseType       = 10
		SchoolingLevel = 11
		From           = 12
	)
	var (
		Dob                              = 0
		IdentifyingCharacteristics       = 1
		CircumstancesBehindDissapearance = 2
	)
	missing := mpp.MissingPersonPoster{}
	identifyingCharacteristics := []string{
		"Registro: " + strings.TrimSpace(doc.Query(".proile-rating span").Text()),
	}
	if data := doc.QueryAll("p.color-subtitulo-theme1"); len(data) == 13 {
		missing.MpSex = ParseChisSex(strings.TrimSpace(data[Sex].Text()))
		if height, err := ParseChisHeigth(strings.TrimSpace(data[Heigth].Text())); err == nil {
			missing.MpHeight = height
		}
		if weigth, err := ParseChisWeight(strings.TrimSpace(data[Weight].Text())); err == nil {
			missing.MpWeight = weigth
		}
		missing.MpEyesDescription = strings.TrimSpace(data[Eyes].Text())
		missing.MpHairDescription = strings.TrimSpace(data[Hair].Text())
		if md, err := ParseChisDate(strings.TrimSpace(data[MissingDate].Text())); err == nil {
			missing.MissingDate = md
		}
		missing.MpPhysicalBuild = ParseChisBuild(strings.TrimSpace(data[Build].Text()))
		missing.MpComplexion = ParseChisComplexion(strings.TrimSpace(data[Complexion].Text()))
		identifyingCharacteristics = append(identifyingCharacteristics, []string{
			"Boca: " + strings.TrimSpace(data[Mouth].Text()),
			"Tama\u00F1o de nariz: " + strings.TrimSpace(data[NoseSize].Text()),
			"Tipo de nariz: " + strings.TrimSpace(data[NoseType].Text()),
			"Escolaridad: " + strings.TrimSpace(data[SchoolingLevel].Text()),
			"Originario de: " + strings.TrimSpace(data[From].Text()),
		}...)
	}
	if moreData := doc.QueryAll(".profile-work p"); len(moreData) == 3 {
		missing.MpDob, _ = ParseChisDate(strings.TrimSpace(moreData[Dob].Text()))
		identifyingCharacteristics = append(identifyingCharacteristics, "Se\u00F1as particulares: "+strings.TrimSpace(moreData[IdentifyingCharacteristics].Text()))
		missing.CircumstancesBehindDissapearance = strings.TrimSpace(moreData[CircumstancesBehindDissapearance].Text())
	}
	missing.MpIdentifyingCharacteristics = strings.Join(identifyingCharacteristics, ", ")
	return &missing, nil
}

func ScrapeChisHasVistoAAlerts(d *doc.Doc) ([]mpp.MissingPersonPoster, map[int]error) {
	mpps := []mpp.MissingPersonPoster{}
	errs := make(map[int]error)
	for i, div := range d.QueryAll(".column_hasvistoa") {
		a := div.Query(".nombre")
		mpName := cases.Title(language.LatinAmericanSpanish).String(strings.TrimSpace(a.Text()))
		if mpName == "" {
			errs[i+1] = fmt.Errorf("MpName can't be empty")
			continue
		}
		postUrl := a.AttrOr("href", "")
		if postUrl == "" {
			errs[i+1] = fmt.Errorf("PoPostUrl can't be empty")
			continue
		}
		poPostUrl, err := url.Parse("https://www.fge.chiapas.gob.mx" + postUrl)
		if err != nil {
			errs[i+1] = fmt.Errorf("can't parse PoPostUrl: %s", err)
			continue
		}
		poPosterUrl, _ := url.Parse(div.Query(".contenido-img img").AttrOr("src", ""))
		found := ParseChisFound(strings.TrimSpace(div.Query("span").Text()))
		missing := mpp.MissingPersonPoster{
			AlertType:   mpp.AlertTypeHasVistoA,
			Found:       found,
			MpName:      mpName,
			PoPosterUrl: poPosterUrl,
			PoPostUrl:   poPostUrl,
			PoState:     mpp.StateChiapas,
		}
		if mppData, err := ScrapeChisHasVistoAExtraData(poPostUrl.String()); err == nil {
			missing.CircumstancesBehindDissapearance = mppData.CircumstancesBehindDissapearance
			missing.MissingDate = mppData.MissingDate
			missing.MpComplexion = mppData.MpComplexion
			missing.MpDob = mppData.MpDob
			missing.MpEyesDescription = mppData.MpEyesDescription
			missing.MpHairDescription = mppData.MpHairDescription
			missing.MpHeight = mppData.MpHeight
			missing.MpIdentifyingCharacteristics = mppData.MpIdentifyingCharacteristics
			missing.MpPhysicalBuild = mppData.MpPhysicalBuild
			missing.MpSex = mppData.MpSex
			missing.MpWeight = mppData.MpWeight
		}
		mpps = append(mpps, missing)
	}
	return mpps, errs
}
