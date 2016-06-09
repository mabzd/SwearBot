package modicm

import (
	"../../mods"
	"../../utils"
	"os"
	"testing"
	"time"
)

var testAsyncChan chan mods.Response = make(chan mods.Response)

func TestResponses(t *testing.T) {
	settingsFilePath := createTmpSettingsPath(t)
	configFilePath := createTmpConfigPath(t)
	defer os.Remove(settingsFilePath)
	defer os.Remove(configFilePath)

	config := newTestConfig()
	modIcm := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, modIcm, "i", "d=20060102x=3y=4", true)
	assertProcessMention(t, modIcm, "i PlaceA", "d=20060102x=1y=2", true)
	assertProcessMention(t, modIcm, "i PlaceB", "d=20060102x=3y=4", true)
	assertProcessMention(t, modIcm, "i PlaceC", "NoPlace 'PlaceC'", false)
	assertProcessMention(t, modIcm, "ia PlaceC 5 6", "PlaceAdded 'PlaceC'", false)
	assertProcessMention(t, modIcm, "i PlaceC", "d=20060102x=5y=6", true)
	assertProcessMention(t, modIcm, "ia  _a-B cd   7 8", "PlaceAdded '_a-B cd'", false)
	assertProcessMention(t, modIcm, "i _a-B cd", "d=20060102x=7y=8", true)
	assertProcessMention(t, modIcm, "ia PlaceC 9 10", "PlaceExists 'PlaceC'", false)
	assertProcessMention(t, modIcm, "ir PlaceD", "PlaceNotExists 'PlaceD'", false)
	assertProcessMention(t, modIcm, "is PlaceD", "PlaceNotExists 'PlaceD'", false)
	assertProcessMention(t, modIcm, "is PlaceA", "ImplicitPlaceSet 'PlaceA'", false)
	assertProcessMention(t, modIcm, "i", "d=20060102x=1y=2", true)
	assertProcessMention(t, modIcm, "ir  placea ", "PlaceRemoved 'placea'", false)
	assertProcessMention(t, modIcm, "i", "NoPlace 'PlaceA'", false)
	assertProcessMention(t, modIcm, "i  pLACeb ", "d=20060102x=3y=4", true)
	assertProcessMention(t, modIcm, "ia  PlacEb  0 0", "PlaceExists 'PlacEb'", false)
	assertProcessMention(t, modIcm, "is plACEB", "ImplicitPlaceSet 'plACEB'", false)
	assertProcessMention(t, modIcm, "i", "d=20060102x=3y=4", true)
}

func TestNoImplicitPlace(t *testing.T) {
	settingsFilePath := createTmpSettingsPath(t)
	configFilePath := createTmpConfigPath(t)
	defer os.Remove(settingsFilePath)
	defer os.Remove(configFilePath)

	config := newTestConfig()
	config.DefaultImplicitPlaceName = ""
	modIcm := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, modIcm, "i", "NoImplicitPlace", false)
}

func TestGetWeatherError(t *testing.T) {
	settingsFilePath := createTmpSettingsPath(t)
	configFilePath := createTmpConfigPath(t)
	defer os.Remove(settingsFilePath)
	defer os.Remove(configFilePath)

	config := newTestConfig()
	config.LastModelDateRegex = "(not-matching-regex)"
	modIcm := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, modIcm, "i", "GetWeatherErr", true)
}

func TestMultipleModInitializations(t *testing.T) {
	settingsFilePath := createTmpSettingsPath(t)
	configFilePath := createTmpConfigPath(t)
	defer os.Remove(settingsFilePath)
	defer os.Remove(configFilePath)

	config := newTestConfig()
	m1 := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, m1, "i", "d=20060102x=3y=4", true)
	assertProcessMention(t, m1, "i PlaceC", "NoPlace 'PlaceC'", false)
	assertProcessMention(t, m1, "ia PlaceC 5 6", "PlaceAdded 'PlaceC'", false)
	assertProcessMention(t, m1, "is PlaceA", "ImplicitPlaceSet 'PlaceA'", false)
	assertProcessMention(t, m1, "ir PlaceB", "PlaceRemoved 'PlaceB'", false)

	m2 := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, m2, "i", "d=20060102x=1y=2", true)
	assertProcessMention(t, m2, "i PlaceB", "NoPlace 'PlaceB'", false)
	assertProcessMention(t, m2, "i PlaceC", "d=20060102x=5y=6", true)
}

func newTestConfig() *ModIcmConfig {
	return &ModIcmConfig{
		IcmUrl:                   "d={date}x={x}y={y}",
		IcmLastModelDateUrl:      "url",
		LastModelDateRegex:       "<date>([0-9]+)</date>",
		WeatherRegex:             "^i$",
		PlaceWeatherRegex:        "^i ([a-zA-Z0-9-_ ]+)$",
		AddPlaceRegex:            "^ia ([a-zA-Z0-9-_ ]+) ([0-9]+) ([0-9]+)$",
		RemovePlaceRegex:         "^ir ([a-zA-Z0-9-_ ]+)$",
		ImplictPlaceRegex:        "^is ([a-zA-Z0-9-_ ]+)$",
		NoImplicitPlaceResponse:  "NoImplicitPlace",
		NoPlaceResponse:          "NoPlace '{place}'",
		PlaceExistsResponse:      "PlaceExists '{place}'",
		PlaceNotExistsResponse:   "PlaceNotExists '{place}'",
		InvalidXCoordResponse:    "InvalidXCoord",
		InvalidYCoordResponse:    "InvalidYCoord",
		PlaceAddedResponse:       "PlaceAdded '{place}'",
		PlaceRemovedResponse:     "PlaceRemoved '{place}'",
		ImplictPlaceSetResponse:  "ImplicitPlaceSet '{place}'",
		OnGetWeatherErr:          "GetWeatherErr",
		OnSettingsSaveErr:        "SettingsSaveErr",
		DefaultImplicitPlaceName: "PlaceB",
		DefaultPlaces: []IcmPlace{
			IcmPlace{Name: "PlaceA", X: 1, Y: 2},
			IcmPlace{Name: "PlaceB", X: 3, Y: 4},
		},
	}
}

func newTestState(t *testing.T, settingsFilePath string) mods.State {
	state := mods.NewState(nil, testAsyncChan)
	if !state.Init(settingsFilePath) {
		t.Fatal("Cannot init mod state")
	}
	return state
}

func newTestModIcm(
	t *testing.T,
	config *ModIcmConfig,
	settingsFilePath string,
	configFilePath string) *ModIcm {

	state := newTestState(t, settingsFilePath)
	modIcm := NewModIcm()
	modIcm.config = config
	modIcm.configFilePath = configFilePath
	modIcm.getLastModelDateFunc = func(x string) string { return "<date>20060102</date>" }
	if !modIcm.Init(state) {
		t.Fatal("Cannot init ModIcm")
	}
	return modIcm
}

func assertProcessMention(
	t *testing.T,
	m *ModIcm,
	message string,
	expected string,
	async bool) {

	actual := m.ProcessMention(message, "u1", "c1")
	if actual == nil {
		t.Errorf("Message '%s' expected response, got nil", message)
		return
	}
	if async {
		if actual.Message != "" {
			t.Errorf(
				"Message '%s' expected empty response for async call, got '%s'",
				message,
				actual.Message)
		}
		// Pass controll to async go routiness
		time.Sleep(1)
		select {
		case response := <-testAsyncChan:
			actual = &response
		default:
			t.Errorf("Message '%s' expected async response, got nothing", message)
			return
		}
	}
	if actual.Message != expected {
		t.Errorf("Message '%s' expected response '%s', got '%s'", message, expected, actual)
	}
}

func createTmpSettingsPath(t *testing.T) string {
	fileName := utils.CreateTmpFileName("Settings")
	if fileName == "" {
		t.Fatal("Cannot create temp settings file path")
	}
	return fileName
}

func createTmpConfigPath(t *testing.T) string {
	fileName := utils.CreateTmpFileName("Config")
	if fileName == "" {
		t.Fatal("Cannot create temp config file path")
	}
	return fileName
}
