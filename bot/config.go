package bot

import (
	"../mods"
	"../mods/modchoice"
	"../mods/modicm"
	"../mods/modmention"
	"../mods/modswears"
)

func registerMods(container *mods.ModContainer) {
	container.AddMod(modswears.NewModSwears())
	container.AddMod(modchoice.NewModChoice())
	container.AddMod(modmention.NewModMention())
	container.AddMod(modicm.NewModIcm())
}
