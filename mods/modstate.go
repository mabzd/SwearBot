package mods

import (
	"../settings"
	"github.com/nlopes/slack"
	"path"
)

type ModState struct {
	settings      *settings.AllSettings
	SlackClient   *slack.Client
	AsyncResponse chan Response
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
	filePath := getSettingsFilePath()
	return settings.SaveSettings(filePath, s.settings)
}

func newModState(slackClient *slack.Client, asyncResponse chan Response) *ModState {
	return &ModState{
		settings:      settings.NewSettings(),
		SlackClient:   slackClient,
		AsyncResponse: asyncResponse,
	}
}

func (s *ModState) init() bool {
	filePath := getSettingsFilePath()
	settings, err := settings.LoadSettings(filePath)
	if err != Success {
		return false
	}
	s.settings = settings
	return true
}

func getSettingsFilePath() string {
	return path.Join(ModsDirName, SettingsFileName)
}
