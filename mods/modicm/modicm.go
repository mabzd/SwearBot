package modicm

import (
	"../../mods"
	"../../utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	state                mods.State
	config               *ModIcmConfig
	configFilePath       string
	lastModelDateRegex   *regexp.Regexp
	weatherRegex         *regexp.Regexp
	placeWeatherRegex    *regexp.Regexp
	addPlaceRegex        *regexp.Regexp
	removePlaceRegex     *regexp.Regexp
	implictPlaceRegex    *regexp.Regexp
	getLastModelDateFunc func(string) string
}

func NewModIcm() *ModIcm {
	modIcm := &ModIcm{
		config:               NewModIcmConfig(),
		getLastModelDateFunc: getLastModelDate,
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
	m.lastModelDateRegex = compileRegex(m.config.LastModelDateRegex, "LastModelDateRegex", 1)
	if m.lastModelDateRegex == nil {
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

func (m *ModIcm) ProcessMention(
	message string,
	userId string,
	channelId string) *mods.Response {

	var groups [][]string
	if m.weatherRegex.MatchString(message) {
		return m.asyncImplicitPlaceWeather(channelId)
	}
	groups = m.placeWeatherRegex.FindAllStringSubmatch(message, 1)
	if groups != nil {
		return m.asyncPlaceWeather(groups[0][1], channelId)
	}
	groups = m.addPlaceRegex.FindAllStringSubmatch(message, 1)
	if groups != nil {
		return response(m.addNewPlace(groups[0][1], groups[0][2], groups[0][3]), channelId)
	}
	groups = m.removePlaceRegex.FindAllStringSubmatch(message, 1)
	if groups != nil {
		return response(m.removePlace(groups[0][1]), channelId)
	}
	groups = m.implictPlaceRegex.FindAllStringSubmatch(message, 1)
	if groups != nil {
		return response(m.setImplicitPlace(groups[0][1]), channelId)
	}
	return nil
}

func (m *ModIcm) ProcessMessage(
	message string,
	userId string,
	channelId string) *mods.Response {

	return nil
}

func response(message string, channelId string) *mods.Response {
	if message == "" {
		return nil
	}
	return &mods.Response{
		Message:   message,
		ChannelId: channelId,
	}
}

func (m *ModIcm) asyncImplicitPlaceWeather(channelId string) *mods.Response {
	implictPlaceName, ok := m.state.Settings().GetSetting(ImplicitPlaceSetting)
	if !ok {
		implictPlaceName = m.config.DefaultImplicitPlaceName
	}
	if implictPlaceName == "" {
		return &mods.Response{
			Message:   m.config.NoImplicitPlaceResponse,
			ChannelId: channelId,
		}
	}
	return m.asyncPlaceWeather(implictPlaceName, channelId)
}

func (m *ModIcm) asyncPlaceWeather(placeName string, channelId string) *mods.Response {
	placeName = trimWhitespaces(placeName)
	icmPlace, ok := m.getIcmPlace(placeName)
	if !ok {
		return &mods.Response{
			Message:   formatPlaceNameResponse(m.config.NoPlaceResponse, placeName),
			ChannelId: channelId,
		}
	}
	go func() {
		dateResponse := m.getLastModelDateFunc(m.config.IcmLastModelDateUrl)
		groups := m.lastModelDateRegex.FindAllStringSubmatch(dateResponse, 1)
		if groups == nil {
			log.Printf(
				"ModIcm: cannot find last ICM model date in response: %v\n",
				dateResponse)
			m.state.AsyncResponse(mods.Response{
				Message:   m.config.OnGetWeatherErr,
				ChannelId: channelId,
			})
			return
		}
		date := groups[0][1]
		m.state.AsyncResponse(mods.Response{
			Message:   formatUrlResponse(m.config.IcmUrl, date, icmPlace),
			ChannelId: channelId,
		})
	}()
	return &mods.Response{}
}

func (m *ModIcm) addNewPlace(placeName string, xStr string, yStr string) string {
	placeName = trimWhitespaces(placeName)
	_, ok := m.getIcmPlace(placeName)
	if ok {
		return formatPlaceNameResponse(m.config.PlaceExistsResponse, placeName)
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
		Name: placeName,
		X:    x,
		Y:    y,
	}
	return m.addIcmPlace(icmPlace)
}

func (m *ModIcm) removePlace(placeName string) string {
	placeName = trimWhitespaces(placeName)
	placeSetting := getPlaceSetting(placeName)
	if !m.state.Settings().RemoveSetting(placeSetting) {
		return formatPlaceNameResponse(m.config.PlaceNotExistsResponse, placeName)
	}
	errNum := m.state.SaveSettings()
	if errNum != Success {
		log.Printf("ModIcm: cannot save settings for removed ICM place '%s'\n", placeName)
		return m.config.OnSettingsSaveErr
	}
	return formatPlaceNameResponse(m.config.PlaceRemovedResponse, placeName)
}

func (m *ModIcm) setImplicitPlace(placeName string) string {
	placeName = trimWhitespaces(placeName)
	_, ok := m.getIcmPlace(placeName)
	if !ok {
		return formatPlaceNameResponse(m.config.PlaceNotExistsResponse, placeName)
	}
	m.state.Settings().SetSetting(ImplicitPlaceSetting, placeName)
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

func getLastModelDate(url string) string {
	resp, httpErr := http.Get(url)
	if httpErr != nil {
		log.Printf("ModIcm: cannot get last ICM model date from '%s': %v\n", url, httpErr)
		return ""
	}
	defer resp.Body.Close()
	body, ioErr := ioutil.ReadAll(resp.Body)
	if ioErr != nil {
		fmt.Printf("ModIcm: cannot read HTTP last ICM model date response: %v\n", ioErr)
		return ""
	}
	return string(body)
}

func formatUrlResponse(format string, date string, icmPlace IcmPlace) string {
	params := map[string]string{
		"date":  date,
		"place": icmPlace.Name,
		"x":     strconv.Itoa(icmPlace.X),
		"y":     strconv.Itoa(icmPlace.Y),
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
