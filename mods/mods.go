package mods

import (
	"../utils"
	"github.com/nlopes/slack"
	"log"
	"path"
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
	Init(state *ModState) bool
	ProcessMention(message string, userId string, channelId string) string
	ProcessMessage(message string, userId string, channelId string) string
}

type ModInfo struct {
	Name     string
	Enabled  bool
	Priority int
	Instance Mod
}

type ModContainer struct {
	init     bool
	modInfos []*ModInfo
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
		modInfos: []*ModInfo{},
	}
}

func (mc *ModContainer) LoadConfig() bool {
	filePath := getModConfigFilePath()
	err := utils.LoadJson(filePath, &mc.modInfos)
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
	modState := newModState(slackClient)
	if !modState.init() {
		log.Println("ModContainer: mod state failed to initialize")
		return false
	}
	modsRegistered := []string{}
	modsEnabled := []string{}
	modsInitialized := []string{}
	for _, modInfo := range mc.modInfos {
		if modInfo.Instance != nil {
			modsRegistered = append(modsRegistered, modInfo.Name)
			if modInfo.Enabled {
				modsEnabled = append(modsEnabled, modInfo.Name)
				if modInfo.Instance.Init(modState) {
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
		return mod.ProcessMention(message, userId, channelId)
	})
}

func (mc *ModContainer) ProcessMessage(
	message string,
	userId string,
	channelId string) string {

	return mc.executeOnActiveMod(func(mod Mod) string {
		return mod.ProcessMessage(message, userId, channelId)
	})
}

func GetPath(mod Mod, fileName string) string {
	return path.Join(ModsDirName, mod.Name(), fileName)
}

func (mc *ModContainer) executeOnActiveMod(action func(Mod) string) string {
	for _, modInfo := range mc.modInfos {
		if modInfo.Instance != nil && modInfo.Enabled {
			response := action(modInfo.Instance)
			if response != "" {
				return response
			}
		}
	}
	return ""
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
