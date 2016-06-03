package modswears

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestSwears(t *testing.T) {
	tmpFileName := createTmpDict(t)
	defer os.Remove(tmpFileName)

	mod := createSwears(t, tmpFileName)
	expected := []string{"a", "abcd", "abba"}
	assertFindSwears(t, mod, "Test A abCD abcde abBA abc", expected)
}

func TestAddRule(t *testing.T) {
	tmpFileName := createTmpDict(t)
	defer os.Remove(tmpFileName)

	mod := createSwears(t, tmpFileName)
	assertAddRule(t, mod, "Fgh*")
	expected := []string{"fgh", "abcd", "fghi"}
	assertFindSwears(t, mod, "Test FGH abcd fghi fgi", expected)
}

func TestAddRuleFileReadErr(t *testing.T) {
	tmpFileName := createTmpDict(t)
	defer os.Remove(tmpFileName)

	mod := createSwears(t, tmpFileName)
	os.Remove(tmpFileName)
	assertAddRuleErr(t, mod, "xxx*", DictFileReadErr)
}

func TestAddRuleConflictErr(t *testing.T) {
	tmpFileName := createTmpDict(t)
	defer os.Remove(tmpFileName)

	mod := createSwears(t, tmpFileName)
	assertAddRuleErr(t, mod, "abc*", AddRuleConflictErr)
	assertAddRuleErr(t, mod, "ab*", AddRuleConflictErr)
	assertAddRuleErr(t, mod, "*", AddRuleConflictErr)
}

func TestAddRuleInvalidWildcardErr(t *testing.T) {
	tmpFileName := createTmpDict(t)
	defer os.Remove(tmpFileName)

	mod := createSwears(t, tmpFileName)
	assertAddRuleErr(t, mod, "xx*x", InvalidWildcardErr)
	assertAddRuleErr(t, mod, "*dd*", InvalidWildcardErr)
	assertAddRuleErr(t, mod, "**x1", InvalidWildcardErr)
	assertAddRuleErr(t, mod, "x2**", InvalidWildcardErr)
	assertAddRuleErr(t, mod, "**", InvalidWildcardErr)
}

func createSwears(t *testing.T, tmpFilePath string) *ModSwears {
	mod := NewModSwears()
	mod.dictFileName = tmpFilePath
	err := mod.LoadSwears()
	if err != Success {
		t.Fatalf("Expected to load dictionary without errors, got %v", err)
	}

	return mod
}

func createTmpDict(t *testing.T) string {
	tmpFile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}
	addDictData(tmpFile)
	tmpFile.Close()
	return tmpFile.Name()
}

func addDictData(file *os.File) {
	rules := []string{"a", "abcd", "abb*"}
	for _, rule := range rules {
		file.WriteString(fmt.Sprintf("%s\n", rule))
	}
}

func assertFindSwears(t *testing.T, mod *ModSwears, m string, expected []string) {
	actual := mod.FindSwears(m)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected swears %#v, got %#v", expected, actual)
	}
}

func assertAddRule(t *testing.T, mod *ModSwears, r string) {
	err := mod.AddRule(r)
	if err != Success {
		t.Fatalf("Expected no errors when adding rule '%s', got: %v", r, err)
	}
}

func assertAddRuleErr(t *testing.T, mod *ModSwears, r string, expected int) {
	err := mod.AddRule(r)
	if err != expected {
		t.Fatalf("Expected error %v when adding rule '%s', got %v", expected, r, err)
	}
}
