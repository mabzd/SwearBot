package modicm

import (
	"../../mods"
	"../../utils"
	"os"
	"testing"
)

func init() {
	utils.TimeClock = utils.MockClock{
		CurrentTime: utils.NewLocalDate(2006, 1, 2),
	}
}

func TestResponses(t *testing.T) {
	settingsFilePath := createTmpSettingsPath(t)
	configFilePath := createTmpConfigPath(t)
	defer os.Remove(settingsFilePath)
	defer os.Remove(configFilePath)

	config := newTestConfig()
	modIcm := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, modIcm, "i", "d=20060102x=3y=4")
	assertProcessMention(t, modIcm, "i PlaceA", "d=20060102x=1y=2")
	assertProcessMention(t, modIcm, "i PlaceB", "d=20060102x=3y=4")
	assertProcessMention(t, modIcm, "i PlaceC", "NoPlace 'PlaceC'")
	assertProcessMention(t, modIcm, "ia PlaceC 5 6", "PlaceAdded 'PlaceC'")
	assertProcessMention(t, modIcm, "i PlaceC", "d=20060102x=5y=6")
	assertProcessMention(t, modIcm, "ia  _a-B cd   7 8", "PlaceAdded '_a-B cd'")
	assertProcessMention(t, modIcm, "i _a-B cd", "d=20060102x=7y=8")
	assertProcessMention(t, modIcm, "ia PlaceC 9 10", "PlaceExists 'PlaceC'")
	assertProcessMention(t, modIcm, "ir PlaceD", "PlaceNotExists 'PlaceD'")
	assertProcessMention(t, modIcm, "is PlaceD", "PlaceNotExists 'PlaceD'")
	assertProcessMention(t, modIcm, "is PlaceA", "ImplicitPlaceSet 'PlaceA'")
	assertProcessMention(t, modIcm, "i", "d=20060102x=1y=2")
	assertProcessMention(t, modIcm, "ir  placea ", "PlaceRemoved 'placea'")
	assertProcessMention(t, modIcm, "i", "NoPlace 'PlaceA'")
	assertProcessMention(t, modIcm, "i  pLACeb ", "d=20060102x=3y=4")
	assertProcessMention(t, modIcm, "ia  PlacEb  0 0", "PlaceExists 'PlacEb'")
	assertProcessMention(t, modIcm, "is plACEB", "ImplicitPlaceSet 'plACEB'")
	assertProcessMention(t, modIcm, "i", "d=20060102x=3y=4")
}

func TestNoImplicitPlace(t *testing.T) {
	settingsFilePath := createTmpSettingsPath(t)
	configFilePath := createTmpConfigPath(t)
	defer os.Remove(settingsFilePath)
	defer os.Remove(configFilePath)

	config := newTestConfig()
	config.DefaultImplicitPlaceName = ""
	modIcm := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, modIcm, "i", "NoImplicitPlace")
}

func TestMultipleModInitializations(t *testing.T) {
	settingsFilePath := createTmpSettingsPath(t)
	configFilePath := createTmpConfigPath(t)
	defer os.Remove(settingsFilePath)
	defer os.Remove(configFilePath)

	config := newTestConfig()
	m1 := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, m1, "i", "d=20060102x=3y=4")
	assertProcessMention(t, m1, "i PlaceC", "NoPlace 'PlaceC'")
	assertProcessMention(t, m1, "ia PlaceC 5 6", "PlaceAdded 'PlaceC'")
	assertProcessMention(t, m1, "is PlaceA", "ImplicitPlaceSet 'PlaceA'")
	assertProcessMention(t, m1, "ir PlaceB", "PlaceRemoved 'PlaceB'")

	m2 := newTestModIcm(t, config, settingsFilePath, configFilePath)
	assertProcessMention(t, m2, "i", "d=20060102x=1y=2")
	assertProcessMention(t, m2, "i PlaceB", "NoPlace 'PlaceB'")
	assertProcessMention(t, m2, "i PlaceC", "d=20060102x=5y=6")
}

func newTestConfig() *ModIcmConfig {
	return &ModIcmConfig{
		IcmUrl:                   "d={date}x={x}y={y}",
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
		OnSettingsSaveErr:        "SettingsSaveErr",
		DefaultImplicitPlaceName: "PlaceB",
		DefaultPlaces: []IcmPlace{
			IcmPlace{Name: "PlaceA", X: 1, Y: 2},
			IcmPlace{Name: "PlaceB", X: 3, Y: 4},
		},
	}
}

func newTestState(t *testing.T, settingsFilePath string) mods.State {
	state := mods.NewState(nil, nil)
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
	if !modIcm.Init(state) {
		t.Fatal("Cannot init ModIcm")
	}
	return modIcm
}

func assertProcessMention(t *testing.T, m *ModIcm, message string, expected string) {
	actual := m.ProcessMention(message, "u1", "c1")
	if actual != expected {
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
