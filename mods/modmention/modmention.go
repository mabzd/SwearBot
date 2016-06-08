package modmention

import (
	"../../mods"
	"../../utils"
	"log"
	"math/rand"
)

const (
	ConfigFileName = "config.json"
)

type ModMention struct {
	state     mods.State
	config    *ModMentionConfig
	maxWeight int
}

func NewModMention() *ModMention {
	return &ModMention{
		config: NewModMentionConfig(),
	}
}

func (mod *ModMention) Name() string {
	return "modmention"
}

func (mod *ModMention) Init(state mods.State) bool {
	mod.state = state
	configFilePath := mods.GetPath(mod, ConfigFileName)
	err := utils.JsonFromFileCreate(configFilePath, mod.config)
	if err != nil {
		log.Printf("ModMention: cannot load config.")
		return false
	}
	return mod.validateConfig()
}

func (mod *ModMention) ProcessMention(
	message string,
	userId string,
	channelId string) string {

	distribution := rand.Intn(mod.maxWeight)
	return mod.getReaction(distribution)
}

func (mod *ModMention) ProcessMessage(
	message string,
	userId string,
	channelId string) string {

	return ""
}

func (mod *ModMention) validateConfig() bool {
	validReactions := []*Reaction{}
	for _, reaction := range mod.config.Reactions {
		if reaction.Weight <= 0 {
			log.Printf(
				"ModMention: '%d' is not valid reaction weight. Corrected to '1'\n",
				reaction.Weight)
			reaction.Weight = 1
		} else if reaction.Weight >= 1000 {
			log.Printf(
				"ModMention: '%d' is not valid reaction weight. Corrected to '1000'\n",
				reaction.Weight)
			reaction.Weight = 1000
		}
		if len(reaction.Responses) > 0 {
			mod.maxWeight += reaction.Weight * len(reaction.Responses)
			reaction.Weight = mod.maxWeight
			validReactions = append(validReactions, reaction)
		} else {
			log.Println("ModMention: omitted empty reaction")
		}
	}
	if len(validReactions) == 0 {
		log.Println("ModMention: no valid reactions in config file.")
		return false
	}
	mod.config.Reactions = validReactions
	return true
}

func (mod *ModMention) getReaction(distribution int) string {
	for _, reaction := range mod.config.Reactions {
		if distribution < reaction.Weight {
			return utils.RandSelect(reaction.Responses)
		}
	}
	log.Printf(
		"ModMention: invalid distribution %d. Should be from range [0, %d)\n",
		distribution,
		mod.maxWeight)
	return ""
}
