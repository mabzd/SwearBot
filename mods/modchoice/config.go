package modchoice

type ModChoiceConfig struct {
	OrKeywords            []string
	ChoiceResponseFormat  []string
	NullChoiceResponses   []string
	NullChoiceProbability float32
}

func NewModChoiceConfig() *ModChoiceConfig {
	return &ModChoiceConfig{
		OrKeywords:            []string{"or"},
		ChoiceResponseFormat:  []string{"I choose *{option}*.", "*{option}*!", "I would recommend *{option}*."},
		NullChoiceResponses:   []string{"Choose them all! :-)", "Neither."},
		NullChoiceProbability: 0.1,
	}
}
