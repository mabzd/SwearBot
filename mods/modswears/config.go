package modswears

type ModSwearsConfig struct {
	AddRuleRegex        string
	CurrMonthRankRegex  string
	PrevMonthRankRegex  string
	TotalRankRegex      string
	SwearNotifyOnRegex  string
	SwearNotifyOffRegex string

	SwearFormat              string
	OnSwearsFoundResponse    string
	OnAddRuleResponse        string
	OnEmptyRankResponse      string
	OnSwearNotifyOnResponse  string
	OnSwearNotifyOffResponse string
	MonthlyRankHeaderFormat  string
	TotalRankHeaderFormat    string
	RankLineFormat           string
	MonthNames               []string

	OnUserFetchErr       string
	OnDictFileReadErr    string
	OnAddRuleConflictErr string
	OnAddRuleSaveErr     string
	OnInvalidWildcardErr string

	OnStatsFileReadErr string
	OnStatsSaveErr     string

	OnSettingsFileReadErr string
	OnSettingsSaveErr     string
}

func NewModSwearsConfig() *ModSwearsConfig {
	return &ModSwearsConfig{
		AddRuleRegex:        "(?i)^\\s*add rule:\\s*([a-z0-9*]+)\\s*$",
		CurrMonthRankRegex:  "(?i)^\\s*curr\\s+rank\\s*$",
		PrevMonthRankRegex:  "(?i)^\\s*prev\\s+rank\\s*$",
		TotalRankRegex:      "(?i)^\\s*total\\s+rank\\s*$",
		SwearNotifyOnRegex:  "(?i)^\\s*notify\\s+on\\s*$",
		SwearNotifyOffRegex: "(?i)^\\s*notify\\s+off\\s*$",

		SwearFormat:              "{index}. *{swear}*",
		OnAddRuleResponse:        "Rule '{rule}' added.",
		OnSwearsFoundResponse:    "{count} swears found: {swears}",
		OnEmptyRankResponse:      "Rank is empty.",
		OnSwearNotifyOnResponse:  "Swear notification is on",
		OnSwearNotifyOffResponse: "Swear notification is off",
		MonthlyRankHeaderFormat:  "*Monthly Rank* - {month} {year}",
		TotalRankHeaderFormat:    "*Total Rank*",
		RankLineFormat:           "{index}. *{user}*: {count} swears",
		MonthNames:               []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},

		OnUserFetchErr:        "Error when fetching slack users!",
		OnDictFileReadErr:     "Error when reading database!",
		OnAddRuleConflictErr:  "Similar rule already exists!",
		OnAddRuleSaveErr:      "Error when saving to database!",
		OnInvalidWildcardErr:  "Invalid wildcard placement!",
		OnStatsFileReadErr:    "Error when reading stats file!",
		OnStatsSaveErr:        "Error when saving to stats file!",
		OnSettingsFileReadErr: "Error when reading settings file!",
		OnSettingsSaveErr:     "Error when saving to settings file!",
	}
}
