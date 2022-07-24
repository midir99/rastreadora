package ws

import (
	"testing"

	"github.com/midir99/rastreadora/mpp"
)

func TestParseNameSexFound(t *testing.T) {
	testCases := []struct {
		legend      string
		wantedName  string
		wantedSex   mpp.Sex
		wantedFound bool
	}{
		{
			"Desaparecida; Martina Bello Morales",
			"Martina Bello Morales",
			mpp.SexFemale,
			false,
		},
		{
			"Localizada: Valeria Benítez Domínguez",
			"Valeria Benítez Domínguez",
			mpp.SexFemale,
			true,
		},
		{
			"Fiscalía General del Estado solicita su colaboración para localizar a Milagros Gabriela Leyva Santiago.",
			"Fiscalía General Del Estado Solicita Su Colaboración Para Localizar A Milagros Gabriela Leyva Santiago.",
			mpp.Sex(""),
			false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.legend, func(t *testing.T) {
			name, sex, found := ParseNameSexFound(tc.legend)
			if name != tc.wantedName {
				t.Errorf("got %s; want %s", name, tc.wantedName)
			}
			if sex != tc.wantedSex {
				t.Errorf("got %s; want %s", sex, tc.wantedSex)
			}
			if found != tc.wantedFound {
				t.Errorf("got %t; want %t", found, tc.wantedFound)
			}
		})
	}
}
