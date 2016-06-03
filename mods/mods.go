package mods

import (
	"../settings"
	"encoding/json"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"path"
)

const (
	Success          = 0
	ModsDirName      = "mods"
	SettingsFileName = "settings.json"
)

type Mod interface {
	Name() string
	Init(state *ModState) bool
	ProcessMention(message string, userId string, channelId string) string
	ProcessMessage(message string, userId string, channelId string) string
}

type ModState struct {
	settings    *settings.AllSettings
	SlackClient *slack.Client
}

func NewModState(slackClient *slack.Client) *ModState {
	settings, err := settings.LoadSettings(SettingsFileName)
	if err != Success {
		return nil
	}
	return &ModState{
		settings:    settings,
		SlackClient: slackClient,
	}
}

func (s *ModState) GetUserChanSetting(
	userId string,
	channelId string,
	key string) (string, bool) {

	return s.settings.GetUserChanSetting(userId, channelId, key)
}

func (s *ModState) GetUserSetting(userId string, key string) (string, bool) {
	return s.settings.GetUserSetting(userId, key)
}

func (s *ModState) GetChanSetting(channelId string, key string) (string, bool) {
	return s.settings.GetChanSetting(channelId, key)
}

func (s *ModState) GetSetting(key string) (string, bool) {
	return s.settings.GetSetting(key)
}

func (s *ModState) SetUserChanSetting(
	userId string,
	channelId string,
	key string,
	value string) {

	s.settings.SetUserChanSetting(userId, channelId, key, value)
}

func (s *ModState) SetUserSetting(userId string, key string, value string) {
	s.settings.SetUserSetting(userId, key, value)
}

func (s *ModState) SetChanSetting(channelId string, key string, value string) {
	s.settings.SetChanSetting(channelId, key, value)
}

func (s *ModState) SetSetting(key string, value string) {
	s.settings.SetSetting(key, value)
}

func (s *ModState) Save() int {
	return settings.SaveSettings(SettingsFileName, s.settings)
}

func LoadConfig(fileName string, config interface{}) error {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Cannot read config from file '%s': %s", fileName, err)
		return err
	}

	err = json.Unmarshal(bytes, config)
	if err != nil {
		log.Printf("Error when parsing config file '%s' JSON: %s", fileName, err)
		return err
	}

	return nil
}

func GetPath(mod Mod, fileName string) string {
	return path.Join(ModsDirName, mod.Name(), fileName)
}
