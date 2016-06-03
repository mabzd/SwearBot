package modswears

import (
	"../../mods"
	"../../utils"
	"strings"
	"testing"
)

func TestConfigIntegrity(t *testing.T) {
	var config ModSwearsConfig
	err := mods.LoadConfig("../../modswears-config-rename.json", &config)
	if err != nil {
		t.Fatalf("Loading config failed %v", err)
	}
	emptyFields := utils.GetEmptyFieldNames(config)
	if len(emptyFields) > 0 {
		t.Fatalf("Found empty fields in config: %s", strings.Join(emptyFields, ", "))
	}
}
