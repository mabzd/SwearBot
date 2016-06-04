package modmention

type Reaction struct {
	Weight    int
	Responses []string
}

type ModMentionConfig struct {
	Reactions []*Reaction
}

func NewModMentionConfig() *ModMentionConfig {
	return &ModMentionConfig{
		Reactions: []*Reaction{
			&Reaction{
				Weight: 10,
				Responses: []string{
					"Sorry, I can't understand you ;_;",
					"Don't know what that means ;_;",
				},
			},
			&Reaction{
				Weight: 1,
				Responses: []string{
					"Hmmm...",
					"What?",
				},
			},
		},
	}
}
