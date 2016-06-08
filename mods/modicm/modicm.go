package modicm

import (
	"../../mods"
	"../../utils"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	Success                    = 0
	ConfigFileName             = "config.json"
	ImplicitPlaceSetting       = "ModIcm.ImplicitPlaceName"
	IcmPlaceSettingPrefix      = "ModIcm.Place"
	DefaultPlacesLoadedSetting = "ModIcm.DefaultPlacesLoaded"
)

type ModIcm struct {
	state             mods.State
	config            *ModIcmConfig
	configFilePath    string
	weatherRegex      *regexp.Regexp
	placeWeatherRegex *regexp.Regexp
	addPlaceRegex     *regexp.Regexp
	removePlaceRegex  *regexp.Regexp
	implictPlaceRegex *regexp.Regexp
}

func NewModIcm() *ModIcm {
	modIcm := &ModIcm{
		config: NewModIcmConfig(),
	}
	modIcm.configFilePath = mods.GetPath(modIcm, ConfigFileName)
	return modIcm
}

func (m *ModIcm) Name() string {
	return "modicm"
}

func (m *ModIcm) Init(state mods.State) bool {
	var err error
	m.state = state
	err = utils.JsonFromFileCreate(m.configFilePath, m.config)
	if err != nil {
		log.Println("ModIcm: cannot load config.")
		return false
	}
	m.weatherRegex = compileRegex(m.config.WeatherRegex, "WeatherRegex", 0)
	if m.weatherRegex == nil {
		return false
	}
	m.placeWeatherRegex = compileRegex(m.config.PlaceWeatherRegex, "PlaceWeatherRegex", 1)
	if m.placeWeatherRegex == nil {
		return false
	}
	m.addPlaceRegex = compileRegex(m.config.AddPlaceRegex, "AddPlaceRegex", 3)
	if m.addPlaceRegex == nil {
		return false
	}
	m.removePlaceRegex = compileRegex(m.config.RemovePlaceRegex, "RemovePlaceRegex", 1)
	if m.removePlaceRegex == nil {
		return false
	}
	m.implictPlaceRegex = compileRegex(m.config.ImplictPlaceRegex, "ImplictPlaceRegex", 1)
	if m.implictPlaceRegex == nil {
		return false
	}
	return m.validateConfig()
}

func (m *ModIcm) ProcessMention(message string, userId string, channelId string) string {
	var groups [][]string
	if m.weatherRegex.MatchString(message) {
		return m.getImplicitPlaceWeather()
	}
	groups = m.placeWeatherRegex.FindAllStringSubmatch(message, 1)
	if groups != nil {
		return m.getPlaceWeather(groups[0][1])
	}
	groups = m.addPlaceRegex.FindAllStringSubmatch(message, 1)
	if groups != nil {
		return m.addNewPlace(groups[0][1], groups[0][2], groups[0][3])
	}
	groups = m.removePlaceRegex.FindAllStringSubmatch(message, 1)
	if groups != nil {
		return m.removePlace(groups[0][1])
	}
	groups = m.implictPlaceRegex.FindAllStringSubmatch(message, 1)
	if groups != nil {
		return m.setImplicitPlace(groups[0][1])
	}
	return ""
}

func (m *ModIcm) ProcessMessage(message string, userId string, channelId string) string {
	return ""
}

func (m *ModIcm) getImplicitPlaceWeather() string {
	implictPlaceName, ok := m.state.Settings().GetSetting(ImplicitPlaceSetting)
	if !ok {
		implictPlaceName = m.config.DefaultImplicitPlaceName
	}
	if implictPlaceName == "" {
		return m.config.NoImplicitPlaceResponse
	}
	return m.getPlaceWeather(implictPlaceName)
}

func (m *ModIcm) getPlaceWeather(placeName string) string {
	icmPlace, ok := m.getIcmPlace(placeName)
	if !ok {
		return m.config.NoPlaceResponse
	}
	return formatUrlResponse(m.config.IcmUrl, icmPlace.X, icmPlace.Y)
}

func (m *ModIcm) addNewPlace(placeName string, xStr string, yStr string) string {
	_, ok := m.getIcmPlace(placeName)
	if ok {
		return m.config.PlaceExistsResponse
	}
	x, errX := strconv.Atoi(xStr)
	if errX != nil {
		return m.config.InvalidXCoordResponse
	}
	y, errY := strconv.Atoi(yStr)
	if errY != nil {
		return m.config.InvalidYCoordResponse
	}
	icmPlace := IcmPlace{
		Name: trimWhitespaces(placeName),
		X:    x,
		Y:    y,
	}
	return m.addIcmPlace(icmPlace)
}

func (m *ModIcm) removePlace(placeName string) string {
	_, ok := m.getIcmPlace(placeName)
	if !ok {
		return m.config.PlaceNotExistsResponse
	}
	// TODO: implement
	return "Not implemented yet"
}

func (m *ModIcm) setImplicitPlace(placeName string) string {
	_, ok := m.getIcmPlace(placeName)
	if !ok {
		return m.config.PlaceNotExistsResponse
	}
	m.state.Settings().SetSetting(ImplicitPlaceSetting, normalizePlaceName(placeName))
	errNum := m.state.SaveSettings()
	if errNum != Success {
		log.Println("ModIcm: cannot save implicit place '%s' setting.", placeName)
		return m.config.OnSettingsSaveErr
	}
	return formatPlaceNameResponse(m.config.ImplictPlaceSetResponse, placeName)
}

func (m *ModIcm) addIcmPlace(icmPlace IcmPlace) string {
	placeName := icmPlace.Name
	err := m.addIcmPlaceToSettings(icmPlace)
	if err != nil {
		log.Printf("ModIcm: error on adding ICM place '%s' to settings: %v\n", placeName, err)
		return m.config.OnSettingsSaveErr
	}
	errNum := m.state.SaveSettings()
	if errNum != Success {
		log.Printf("ModIcm: error on saving settings for ICM place '%s': %v\n", placeName, err)
		return m.config.OnSettingsSaveErr
	}
	return formatIcmPlaceResponse(m.config.PlaceAddedResponse, icmPlace)
}

func (m *ModIcm) getIcmPlace(placeName string) (IcmPlace, bool) {
	var icmPlace IcmPlace
	placeSetting := getPlaceSetting(placeName)
	icmPlaceStr, ok := m.state.Settings().GetSetting(placeSetting)
	if !ok {
		return icmPlace, false
	}
	err := json.Unmarshal([]byte(icmPlaceStr), &icmPlace)
	if err != nil {
		log.Printf("ModIcm: error on unmarshaling JSON for ICM place '%s': %v\n", placeName, err)
		return icmPlace, false
	}
	return icmPlace, true
}

func (m *ModIcm) addIcmPlaceToSettings(icmPlace IcmPlace) error {
	placeSetting := getPlaceSetting(icmPlace.Name)
	icmPlaceBytes, err := json.Marshal(icmPlace)
	if err != nil {
		return err
	}
	m.state.Settings().SetSetting(placeSetting, string(icmPlaceBytes))
	return nil
}

func formatUrlResponse(format string, x int, y int) string {
	params := map[string]string{
		"date": utils.TimeClock.Now().Format("20060102"),
		"x":    strconv.Itoa(x),
		"y":    strconv.Itoa(y),
	}
	return utils.ParamFormat(format, params)
}

func formatIcmPlaceResponse(format string, icmPlace IcmPlace) string {
	params := map[string]string{
		"place": icmPlace.Name,
		"x":     strconv.Itoa(icmPlace.X),
		"y":     strconv.Itoa(icmPlace.Y),
	}
	return utils.ParamFormat(format, params)
}

func formatPlaceNameResponse(format string, placeName string) string {
	params := map[string]string{
		"place": placeName,
	}
	return utils.ParamFormat(format, params)
}

func compileRegex(regex string, name string, groups int) *regexp.Regexp {
	re, err := regexp.Compile(regex)
	if err != nil {
		log.Printf("ModIcm: cannot compile %s: %v\n", name, err)
		return nil
	}
	if groups > 0 && groups != re.NumSubexp() {
		log.Printf(
			"ModIcm: regexp %s must have %d parenthesized groups (%d defined)\n",
			name,
			groups,
			re.NumSubexp())
		return nil
	}
	return re
}

func (m *ModIcm) validateConfig() bool {
	if strings.Trim(m.config.IcmUrl, " \r\n\t") == "" {
		log.Println("ModIcm: icm url is empty.")
		return false
	}
	loaded, _ := m.state.Settings().GetSetting(DefaultPlacesLoadedSetting)
	if loaded != "true" {
		for _, icmPlace := range m.config.DefaultPlaces {
			if m.addIcmPlaceToSettings(icmPlace) != nil {
				log.Printf("ModIcm: cannot add ICM place '%s' to settings.\n", icmPlace.Name)
				return false
			}
		}
		m.state.Settings().SetSetting(DefaultPlacesLoadedSetting, "true")
		errNum := m.state.SaveSettings()
		if errNum != Success {
			log.Println("ModIcm: cannot save settings.")
			return false
		}
	}
	return true
}

func getPlaceSetting(placeName string) string {
	return fmt.Sprintf("%s.%s", IcmPlaceSettingPrefix, normalizePlaceName(placeName))
}

func normalizePlaceName(placeName string) string {
	return strings.ToLower(trimWhitespaces(placeName))
}

func trimWhitespaces(value string) string {
	return strings.Trim(value, " \r\n\t")
}
