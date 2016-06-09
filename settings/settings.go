package settings

import (
	"../utils"
	"log"
)

const (
	Success             = 0
	SettingsFileReadErr = 31
	SettingsSaveErr     = 32
)

type Settings interface {
	GetUserChanSetting(userId string, channelId string, key string) (string, bool)
	GetUserSetting(userId string, key string) (string, bool)
	GetChanSetting(channelId string, key string) (string, bool)
	GetSetting(key string) (string, bool)
	SetUserChanSetting(userId string, channelId string, key string, value string)
	SetUserSetting(userId string, key string, value string)
	SetChanSetting(channelId string, key string, value string)
	SetSetting(key string, value string)
	RemoveUserChanSetting(userId string, chanelId string, key string) bool
	RemoveUserSetting(userId string, key string) bool
	RemoveChanSetting(channelId string, key string) bool
	RemoveSetting(key string) bool
}

type SettingsManager interface {
	Settings
	Load(fileName string) int
	Save(fileName string) int
}

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

func (settings *AllSettings) RemoveUserChanSetting(
	userId string,
	chanelId string,
	key string) bool {

	userSettings, userOk := settings.UserSettings[userId]
	if userOk {
		chanSettings, chanOk := userSettings.ChanSettings[chanelId]
		if chanOk {
			_, ok := chanSettings.Settings[key]
			delete(chanSettings.Settings, key)
			return ok
		}
	}
	return false
}

func (settings *AllSettings) RemoveUserSetting(userId string, key string) bool {
	userSettings, userOk := settings.UserSettings[userId]
	if userOk {
		_, ok := userSettings.Settings[key]
		delete(userSettings.Settings, key)
		return ok
	}
	return false
}

func (settings *AllSettings) RemoveChanSetting(channelId string, key string) bool {
	chanSettings, chanOk := settings.ChanSettings[channelId]
	if chanOk {
		_, ok := chanSettings.Settings[key]
		delete(chanSettings.Settings, key)
		return ok
	}
	return false
}

func (settings *AllSettings) RemoveSetting(key string) bool {
	_, ok := settings.Settings[key]
	delete(settings.Settings, key)
	return ok
}

func (settings *AllSettings) Load(fileName string) int {
	err := utils.JsonFromFileCreate(fileName, settings)
	if err != nil {
		log.Printf("Settings: Cannot read settings from file '%s'\n", fileName)
		return SettingsFileReadErr
	}
	return Success
}

func (settings *AllSettings) Save(fileName string) int {
	err := utils.JsonToFile(fileName, settings)
	if err != nil {
		log.Printf("Settings: Cannot write settings to file '%s'\n", fileName)
		return SettingsSaveErr
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
