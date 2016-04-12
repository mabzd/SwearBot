package swears

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	SettingsFileCreateErr = 31
	SettingsFileReadErr   = 32
	SettingsUnmarshalErr  = 33
	SettingsMarshalErr    = 34
	SettingsSaveErr       = 35
)

type AllSettings struct {
	UserSettings map[string][]*UserSettings
}

type UserSettings struct {
	UserId   string
	Channel  string
	Settings map[string]string
}

func (settings *AllSettings) GetSetting(
	userId string,
	channel string,
	key string) (string, bool) {

	userSettings, ok := settings.UserSettings[userId]
	if ok {
		for _, userSetting := range userSettings {
			if userSetting.Channel == channel {
				value, ok := userSetting.Settings[key]
				return value, ok
			}
		}
	}

	return "", false
}

func (settings *AllSettings) SetSetting(
	userId string,
	channel string,
	key string,
	value string) {

	userSettings, ok := settings.UserSettings[userId]
	if !ok {
		userSettings = []*UserSettings{}
	}

	for _, userSetting := range userSettings {
		if userSetting.Channel == channel {
			userSetting.Settings[key] = value
			return
		}
	}

	userSetting := &UserSettings{
		UserId:   userId,
		Channel:  channel,
		Settings: map[string]string{key: value},
	}

	userSettings = append(userSettings, userSetting)
	settings.UserSettings[userId] = userSettings
}

func (sw *Swears) ReadSettings() (*AllSettings, int) {
	return readSettings(sw.config.SettingsFileName)
}

func (sw *Swears) WriteSettings(settings *AllSettings) int {
	return writeSettings(sw.config.SettingsFileName, settings)
}

func createSettingsFileIfNotExist(fileName string) int {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		settings := &AllSettings{
			UserSettings: make(map[string][]*UserSettings),
		}
		return writeSettings(fileName, settings)
	}

	return Success
}

func readSettings(fileName string) (*AllSettings, int) {
	createErr := createSettingsFileIfNotExist(fileName)
	if createErr != Success {
		log.Println("Settings: Settings file creation failed.")
		return nil, SettingsFileCreateErr
	}

	bytes, fileReadErr := ioutil.ReadFile(fileName)
	if fileReadErr != nil {
		log.Printf("Settings: Cannot read settings from file '%s': %s\n", fileName, fileReadErr)
		return nil, SettingsFileReadErr
	}

	var settings AllSettings
	unmarshalErr := json.Unmarshal(bytes, &settings)
	if unmarshalErr != nil {
		log.Printf("Settings: Error when unmarshaling settings from JSON: %s\n", unmarshalErr)
		return nil, SettingsUnmarshalErr
	}

	return &settings, Success
}

func writeSettings(fileName string, settings *AllSettings) int {
	bytes, marshalErr := json.Marshal(settings)
	if marshalErr != nil {
		log.Printf("Settings: Error when marshaling settings to JSON: %s\n", marshalErr)
		return SettingsMarshalErr
	}

	saveErr := ioutil.WriteFile(fileName, bytes, 0666)
	if saveErr != nil {
		log.Printf("Settings: Cannot write settings to file '%s': %s\n", fileName, saveErr)
		return SettingsSaveErr
	}

	return Success
}
