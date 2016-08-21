package modicm

type ModIcmConfig struct {
	IcmUrl                   string
	IcmLastModelDateUrl      string
	LastModelDateRegex       string
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
	OnGetWeatherErr          string
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
		IcmUrl:                   "*{place}*, {year}-{month}-{day} {hour}:00 http://www.meteo.pl/um/metco/mgram_pict.php?ntype=0u&fdate={date}&row={y}&col={x}&lang=en",
		IcmLastModelDateUrl:      "http://meteo.pl/xml_um_date.php",
		LastModelDateRegex:       "<act_model_date>\\s*([0-9]+)\\s*</act_model_date>",
		WeatherRegex:             "(?i)^\\s*icm\\s*$",
		PlaceWeatherRegex:        "(?i)^\\s*icm\\s+([^\\s]+)\\s*$",
		AddPlaceRegex:            "(?i)^\\s*icm\\s+add\\s+([^\\s]+)\\s+([0-9]+)\\s+([0-9]+)\\s*$",
		RemovePlaceRegex:         "(?i)^\\s*icm\\s+remove\\s+([^\\s]+)\\s*$",
		ImplictPlaceRegex:        "(?i)^\\s*icm\\s+set\\s+([^\\s]+)\\s*$",
		NoImplicitPlaceResponse:  "No implicit place defined, run 'icm set PLACE_NAME'.",
		NoPlaceResponse:          "No place '{place}' defined, run 'icm add {place} X Y'.",
		PlaceExistsResponse:      "Place '{place}' already added, run 'icm remove {place}.",
		PlaceNotExistsResponse:   "Place '{place}' does not exist, run 'icm add {place} X Y'.",
		InvalidXCoordResponse:    "Invalid X coordinate.",
		InvalidYCoordResponse:    "Invalid Y coordinate.",
		PlaceAddedResponse:       "Place '{place}' added.",
		PlaceRemovedResponse:     "Place '{place}' removed",
		ImplictPlaceSetResponse:  "Place '{place}' set as default.",
		OnSettingsSaveErr:        "Error when saving to settings file!",
		OnGetWeatherErr:          "Error when fetching weather data.",
		DefaultImplicitPlaceName: "Poznań",
		DefaultPlaces: []IcmPlace{
			IcmPlace{Name: "Białystok", X: 285, Y: 379},
			IcmPlace{Name: "Bydgoszcz", X: 199, Y: 381},
			IcmPlace{Name: "Gdańsk", X: 210, Y: 346},
			IcmPlace{Name: "GorzówWlkp", X: 152, Y: 390},
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
			IcmPlace{Name: "ZielonaGóra", X: 155, Y: 412},
		},
	}
}
