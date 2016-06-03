package settings

import (
	"../utils"
	"log"
	"os"
)

const (
	Success               = 0
	SettingsFileCreateErr = 31
	SettingsFileReadErr   = 32
	SettingsSaveErr       = 33
)

type AllSettings struct {
	UserSettings map[string]*UserSettings
	ChanSettings map[string]*ChanSettings
	Settings     map[string]string
}

type UserSettings struct {
	UserId       string
	ChanSettings map[string]*ChanSettings
	Settings     map[string]string
}

type ChanSettings struct {
	ChannelId string
	Settings  map[string]string
}

func NewSettings() *AllSettings {
	return &AllSettings{
		UserSettings: map[string]*UserSettings{},
		ChanSettings: map[string]*ChanSettings{},
		Settings:     map[string]string{},
	}
}

func (settings *AllSettings) GetUserChanSetting(
	userId string,
	channelId string,
	key string) (string, bool) {

	userSettings, userOk := settings.UserSettings[userId]
	if userOk {
		chanSettings, chanOk := userSettings.ChanSettings[channelId]
		if chanOk {
			value, ok := chanSettings.Settings[key]
			return value, ok
		}
	}

	return "", false
}

func (settings *AllSettings) GetUserSetting(
	userId string,
	key string) (string, bool) {

	userSettings, userOk := settings.UserSettings[userId]
	if userOk {
		value, ok := userSettings.Settings[key]
		return value, ok
	}

	return "", false
}

func (settings *AllSettings) GetChanSetting(
	channelId string,
	key string) (string, bool) {

	chanSettings, chanOk := settings.ChanSettings[channelId]
	if chanOk {
		value, ok := chanSettings.Settings[key]
		return value, ok
	}

	return "", false
}

func (settings *AllSettings) GetSetting(key string) (string, bool) {
	value, ok := settings.Settings[key]
	return value, ok
}

func (settings *AllSettings) SetUserChanSetting(
	userId string,
	channelId string,
	key string,
	value string) {

	userSettings, userOk := settings.UserSettings[userId]
	if !userOk {
		userSettings = createUserSettings(userId)
		settings.UserSettings[userId] = userSettings
	}

	chanSettings, chanOk := userSettings.ChanSettings[channelId]
	if !chanOk {
		chanSettings = createChanSettings(channelId)
		userSettings.ChanSettings[channelId] = chanSettings
	}

	chanSettings.Settings[key] = value
}

func (settings *AllSettings) SetUserSetting(
	userId string,
	key string,
	value string) {

	userSettings, userOk := settings.UserSettings[userId]
	if !userOk {
		userSettings = createUserSettings(userId)
		settings.UserSettings[userId] = userSettings
	}

	userSettings.Settings[key] = value
}

func (settings *AllSettings) SetChanSetting(
	channelId string,
	key string,
	value string) {

	chanSettings, chanOk := settings.ChanSettings[channelId]
	if !chanOk {
		chanSettings = createChanSettings(channelId)
		settings.ChanSettings[channelId] = chanSettings
	}

	chanSettings.Settings[key] = value
}

func (settings *AllSettings) SetSetting(key string, value string) {
	settings.Settings[key] = value
}

func LoadSettings(fileName string) (*AllSettings, int) {
	createErr := createSettingsFileIfNotExist(fileName)
	if createErr != Success {
		log.Println("Settings: Settings file creation failed.")
		return nil, SettingsFileCreateErr
	}
	var settings AllSettings
	err := utils.LoadJson(fileName, &settings)
	if err != nil {
		log.Printf("Settings: Cannot read settings from file '%s'\n", fileName)
		return nil, SettingsFileReadErr
	}
	return &settings, Success
}

func SaveSettings(fileName string, settings *AllSettings) int {
	err := utils.SaveJson(fileName, settings)
	if err != nil {
		log.Printf("Settings: Cannot write settings to file '%s'\n", fileName)
		return SettingsSaveErr
	}
	return Success
}

func createSettingsFileIfNotExist(fileName string) int {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return SaveSettings(fileName, NewSettings())
	}
	return Success
}

func createUserSettings(userId string) *UserSettings {
	return &UserSettings{
		UserId:       userId,
		ChanSettings: map[string]*ChanSettings{},
		Settings:     map[string]string{},
	}
}

func createChanSettings(channelId string) *ChanSettings {
	return &ChanSettings{
		ChannelId: channelId,
		Settings:  map[string]string{},
	}
}
