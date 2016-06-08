package modicm

type ModIcmConfig struct {
	IcmUrl                   string
	WeatherRegex             string
	PlaceWeatherRegex        string
	AddPlaceRegex            string
	RemovePlaceRegex         string
	ImplictPlaceRegex        string
	NoImplicitPlaceResponse  string
	NoPlaceResponse          string
	PlaceExistsResponse      string
	PlaceNotExistsResponse   string
	InvalidXCoordResponse    string
	InvalidYCoordResponse    string
	PlaceAddedResponse       string
	PlaceRemovedResponse     string
	ImplictPlaceSetResponse  string
	OnSettingsSaveErr        string
	DefaultImplicitPlaceName string
	DefaultPlaces            []IcmPlace
}

type IcmPlace struct {
	Name string
	X    int
	Y    int
}

func NewModIcmConfig() *ModIcmConfig {
	return &ModIcmConfig{
		IcmUrl:                   "http://www.meteo.pl/um/metco/mgram_pict.php?ntype=0u&fdate={date}&row={y}&col={x}&lang=en",
		WeatherRegex:             "(?i)^\\s*icm\\s*$",
		PlaceWeatherRegex:        "(?i)^\\s*icm\\s+([a-zA-Z0-9-_ ]+)\\s*$",
		AddPlaceRegex:            "(?i)^\\s*icm\\s+add\\s+([a-zA-Z0-9-_ ]+)\\s+([0-9]+)\\s+([0-9]+)\\s*$",
		RemovePlaceRegex:         "(?i)^\\s*icm\\s+remove\\s+([a-zA-Z0-9-_ ]+)\\s*$",
		ImplictPlaceRegex:        "(?i)^\\s*icm\\s+set\\s+([a-zA-Z0-9-_ ]+)\\s*$",
		NoImplicitPlaceResponse:  "No place defined, run 'icm set PLACE_NAME'.",
		NoPlaceResponse:          "No such place defined, run 'icm add PLACE_NAME X Y'.",
		PlaceExistsResponse:      "Such place already added, run 'icm remove PLACE_NAME'.",
		PlaceNotExistsResponse:   "Such place does not exist, run 'icm add PLACE_NAME X Y'.",
		InvalidXCoordResponse:    "Invalid X coordinate.",
		InvalidYCoordResponse:    "Invalid Y coordinate.",
		PlaceAddedResponse:       "Place '{place}' added.",
		PlaceRemovedResponse:     "Place '{place}' removed",
		ImplictPlaceSetResponse:  "Place '{place}' set as default.",
		OnSettingsSaveErr:        "Error when saving to settings file!",
		DefaultImplicitPlaceName: "Poznań",
		DefaultPlaces: []IcmPlace{
			IcmPlace{Name: "Białystok", X: 285, Y: 379},
			IcmPlace{Name: "Bydgoszcz", X: 199, Y: 381},
			IcmPlace{Name: "Gdańsk", X: 210, Y: 346},
			IcmPlace{Name: "Gorzów Wlkp", X: 152, Y: 390},
			IcmPlace{Name: "Katowice", X: 215, Y: 461},
			IcmPlace{Name: "Kielce", X: 244, Y: 443},
			IcmPlace{Name: "Kraków", X: 232, Y: 466},
			IcmPlace{Name: "Lublin", X: 277, Y: 432},
			IcmPlace{Name: "Łódź", X: 223, Y: 418},
			IcmPlace{Name: "Olsztyn", X: 240, Y: 363},
			IcmPlace{Name: "Opole", X: 196, Y: 449},
			IcmPlace{Name: "Poznań", X: 180, Y: 400},
			IcmPlace{Name: "Rzeszów", X: 269, Y: 465},
			IcmPlace{Name: "Szczecin", X: 142, Y: 370},
			IcmPlace{Name: "Toruń", X: 209, Y: 383},
			IcmPlace{Name: "Warszawa", X: 250, Y: 406},
			IcmPlace{Name: "Wrocław", X: 181, Y: 436},
			IcmPlace{Name: "Zielona Góra", X: 155, Y: 412},
		},
	}
}
