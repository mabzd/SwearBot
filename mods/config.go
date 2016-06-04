package mods

type ModInfo struct {
	Name     string
	Enabled  bool
	Priority int
	Instance Mod
}

func NewModInfos() []*ModInfo {
	return []*ModInfo{
		&ModInfo{
			Name:     "modswears",
			Enabled:  true,
			Priority: 0,
		},
		&ModInfo{
			Name:     "modchoice",
			Enabled:  true,
			Priority: 0,
		},
	}
}
