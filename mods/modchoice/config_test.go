package modchoice

import (
	"../../utils"
	"strings"
	"testing"
)

func TestConfigIntegrity(t *testing.T) {
	config := NewModChoiceConfig()
	emptyFields := utils.GetEmptyFieldNames(*config)
	if len(emptyFields) > 0 {
		t.Fatalf("Found empty fields in config: %s", strings.Join(emptyFields, ", "))
	}
}
