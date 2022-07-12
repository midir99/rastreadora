package mpp

type State string

const (
	StateCiudadDeMexico             State = "MX-CMX"
	StateAguascalientes             State = "MX-AGU"
	StateBajaCalifornia             State = "MX-BCN"
	StateBajaCaliforniaSur          State = "MX-BCS"
	StateCampeche                   State = "MX-CAM"
	StateCoahuilaDeZaragoza         State = "MX-COA"
	StateColima                     State = "MX-COL"
	StateChiapas                    State = "MX-CHP"
	StateChihuahua                  State = "MX-CHH"
	StateDurango                    State = "MX-DUR"
	StateGuanajuato                 State = "MX-GUA"
	StateGuerrero                   State = "MX-GRO"
	StateHidalgo                    State = "MX-HID"
	StateJalisco                    State = "MX-JAL"
	StateMexico                     State = "MX-MEX"
	StateMichoacanDeOcampo          State = "MX-MIC"
	StateMorelos                    State = "MX-MOR"
	StateNayarit                    State = "MX-NAY"
	StateNuevoLeon                  State = "MX-NLE"
	StateOaxaca                     State = "MX-OAX"
	StatePuebla                     State = "MX-PUE"
	StateQueretaro                  State = "MX-QUE"
	StateQuintanaRoo                State = "MX-ROO"
	StateSanLuisPotosi              State = "MX-SLP"
	StateSinaloa                    State = "MX-SIN"
	StateSonora                     State = "MX-SON"
	StateTabasco                    State = "MX-TAB"
	StateTamaulipas                 State = "MX-TAM"
	StateTlaxcala                   State = "MX-TLA"
	StateVeracruzDeIgnacioDeLaLlave State = "MX-VER"
	StateYucatan                    State = "MX-YUC"
	StateZacatecas                  State = "MX-ZAC"
)

type PhysicalBuild string

const (
	PhysicalBuildSlim    PhysicalBuild = "S"
	PhysicalBuildRegular PhysicalBuild = "R"
	PhysicalBuildHeavy   PhysicalBuild = "H"
)

type Complexion string

const (
	ComplexionVeryLight         Complexion = "VL"
	ComplexionLight             Complexion = "L"
	ComplexionLightIntermediate Complexion = "LI"
	ComplexionDarkIntermediate  Complexion = "DI"
	ComplexionDark              Complexion = "D"
	ComplexionVeryDark          Complexion = "VD"
)

type Sex string

const (
	SexFemale Sex = "F"
	SexMale   Sex = "M"
)

type AlertType string

const (
	AlertTypeAlba      AlertType = "AL"
	AlertTypeAmber     AlertType = "AM"
	AlertTypeHasVistoA AlertType = "HV"
	AlertTypeOdisea    AlertType = "OD"
)

type MissingPersonPoster struct {
	MpName                           string        `json:"mp_name"`
	MpHeight                         uint          `json:"mp_height,omitempty"`
	MpWeight                         uint          `json:"mp_weight,omitempty"`
	MpPhysicalBuild                  PhysicalBuild `json:"mp_physical_build,omitempty"`
	MpComplexion                     Complexion    `json:"mp_complexion,omitempty"`
	MpSex                            Sex           `json:"mp_sex,omitempty"`
	MpDob                            string        `json:"mp_dob,omitempty"`
	MpAgeWhenDisappeared             uint          `json:"mp_age_when_disappeared,omitempty"`
	MpEyesDescription                string        `json:"mp_eyes_description,omitempty"`
	MpHairDescription                string        `json:"mp_hair_description,omitempty"`
	MpOutfitDescription              string        `json:"mp_outfit_description,omitempty"`
	MpIdentifyingCharacteristics     string        `json:"mp_identifying_characteristics,omitempty"`
	CircumstancesBehindDissapearance string        `json:"circumstances_behind_dissapearance,omitempty"`
	MissingFrom                      string        `json:"missing_from,omitempty"`
	MissingDate                      string        `json:"missing_date,omitempty"`
	Found                            bool          `json:"found,omitempty"`
	AlertType                        AlertType     `json:"alert_type,omitempty"`
	PoState                          State         `json:"po_state"`
	PoPostUrl                        string        `json:"po_post_url,omitempty"`
	PoPostPublicationDate            string        `json:"po_post_publication_date,omitempty"`
	PoPosterUrl                      string        `json:"po_poster_url,omitempty"`
	IsMultiple                       bool          `json:"is_multiple,omitempty"`
}
