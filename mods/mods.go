package mods

import (
	"../utils"
	"github.com/nlopes/slack"
	"log"
	"os"
	"path"
	"runtime/debug"
	"sort"
	"strings"
)

const (
	Success           = 0
	ModsDirName       = "mods"
	SettingsFileName  = "settings.json"
	ModConfigFileName = "config.json"
)

type Mod interface {
	Name() string
	Init(state State) bool
	ProcessMention(message string, userId string, channelId string) string
	ProcessMessage(message string, userId string, channelId string) string
}

type Response struct {
	Message   string
	ChannelId string
}

type ModContainer struct {
	modInfos      []*ModInfo
	AsyncResponse chan Response
}

type ByModPriority []*ModInfo

func (a ByModPriority) Len() int {
	return len(a)
}

func (a ByModPriority) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByModPriority) Less(i, j int) bool {
	return a[i].Priority > a[j].Priority
}

func NewModContainer() *ModContainer {
	return &ModContainer{
		modInfos:      NewModInfos(),
		AsyncResponse: make(chan Response),
	}
}

func (mc *ModContainer) LoadConfig() bool {
	os.MkdirAll(ModsDirName, 0777)
	filePath := getModConfigFilePath()
	err := utils.JsonFromFileCreate(filePath, &mc.modInfos)
	if err != nil {
		log.Println("ModContainer: cannot load mod config file.")
		return false
	}
	return true
}

func (mc *ModContainer) AddMod(mod Mod) bool {
	modName := mod.Name()
	modInfo := getModInfoByName(mc.modInfos, modName)
	if modInfo == nil {
		log.Printf("ModContainer: no mod '%s' in mod config file.\n", modName)
		return false
	}
	if modInfo.Instance != nil {
		log.Printf("ModContainer: mod '%s' already added.\n", modName)
		return false
	}
	modInfo.Instance = mod
	return true
}

func (mc *ModContainer) InitMods(slackClient *slack.Client) bool {
	settingsFilePath := getSettingsFilePath()
	state := NewState(slackClient, mc.AsyncResponse)
	if !state.Init(settingsFilePath) {
		log.Println("ModContainer: mod state failed to initialize")
		return false
	}
	modsRegistered := []string{}
	modsEnabled := []string{}
	modsInitialized := []string{}
	for _, modInfo := range mc.modInfos {
		if modInfo.Instance != nil {
			os.MkdirAll(getModDirPath(modInfo.Instance), 0777)
			modsRegistered = append(modsRegistered, modInfo.Name)
			if modInfo.Enabled {
				modsEnabled = append(modsEnabled, modInfo.Name)
				if modInfo.Instance.Init(state) {
					modInfo.Active = true
					modsInitialized = append(modsInitialized, modInfo.Name)
				} else {
					log.Printf(
						"ModContainer: mod '%s' failed to initialize\n",
						modInfo.Name)
				}
			}
		}
	}
	sort.Sort(ByModPriority(mc.modInfos))
	log.Printf("ModContainer: mod initialization complete "+
		"(mods active: %d, mods enabled: %d, mods registered: %d)\n",
		len(modsInitialized),
		len(modsEnabled),
		len(modsRegistered))
	if len(modsInitialized) == 0 {
		log.Println("ModContainer: no active mods")
		return false
	}
	log.Printf("ModContainer: active mods: %s\n", strings.Join(modsInitialized, ", "))
	return true
}

func (mc *ModContainer) ProcessMention(
	message string,
	userId string,
	channelId string) string {

	return mc.executeOnActiveMod(func(mod Mod) string {
		defer recoverMod("ProcessMention", mod.Name(), message, userId, channelId)
		return mod.ProcessMention(message, userId, channelId)
	})
}

func (mc *ModContainer) ProcessMessage(
	message string,
	userId string,
	channelId string) string {

	return mc.executeOnActiveMod(func(mod Mod) string {
		defer recoverMod("ProcessMessage", mod.Name(), message, userId, channelId)
		return mod.ProcessMessage(message, userId, channelId)
	})
}

func GetPath(mod Mod, fileName string) string {
	return path.Join(getModDirPath(mod), fileName)
}

func getModDirPath(mod Mod) string {
	return path.Join(ModsDirName, mod.Name())
}

func (mc *ModContainer) executeOnActiveMod(action func(Mod) string) string {
	for _, modInfo := range mc.modInfos {
		if modInfo.Active {
			response := action(modInfo.Instance)
			if response != "" {
				return response
			}
		}
	}
	return ""
}

func recoverMod(
	function string,
	modName string,
	message string,
	userId string,
	channelId string) {

	if r := recover(); r != nil {
		log.Printf(
			"ModContainer: recovered panicking mod '%s' on call to %s('%s', '%s', '%s'): %s\n",
			modName,
			function,
			message,
			userId,
			channelId,
			r)
		log.Printf("Stacktrace: %s", string(debug.Stack()))
	}
}

func getModInfoByName(modInfos []*ModInfo, name string) *ModInfo {
	for _, modInfo := range modInfos {
		if modInfo.Name == name {
			return modInfo
		}
	}
	return nil
}

func getModConfigFilePath() string {
	return path.Join(ModsDirName, ModConfigFileName)
}

func getSettingsFilePath() string {
	return path.Join(ModsDirName, SettingsFileName)
}
