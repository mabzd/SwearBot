package mods

import (
	"../settings"
	"github.com/nlopes/slack"
	"log"
)

type State interface {
	Settings() settings.Settings
	SaveSettings() int
	SlackClient() *slack.Client
	AsyncResponse(response Response)
}

type state struct {
	settings         settings.SettingsManager
	slackClient      *slack.Client
	asyncResponse    chan Response
	settingsFilePath string
}

func (s *state) Settings() settings.Settings {
	return s.settings
}

func (s *state) SaveSettings() int {
	return s.settings.Save(s.settingsFilePath)
}

func (s *state) SlackClient() *slack.Client {
	return s.slackClient
}

func (s *state) AsyncResponse(response Response) {
	s.asyncResponse <- response
}

func NewState(slackClient *slack.Client, asyncResponse chan Response) *state {
	return &state{
		settings:      settings.NewSettings(),
		slackClient:   slackClient,
		asyncResponse: asyncResponse,
	}
}

func (s *state) Init(settingsFilePath string) bool {
	err := s.settings.Load(settingsFilePath)
	if err != Success {
		log.Printf("Mods: cannot load state settings from file '%s'\n", settingsFilePath)
		return false
	}
	s.settingsFilePath = settingsFilePath
	return true
}
