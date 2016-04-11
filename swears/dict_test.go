package swears

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
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(tmpFileName)
	expected := []string{"a", "abcd", "abba"}
	assertFindSwears(t, sw, "Test A abCD abcde abBA abc", expected)
}

func TestAddRule(t *testing.T) {
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(tmpFileName)
	assertAddRule(t, sw, "Fgh*")
	expected := []string{"fgh", "abcd", "fghi"}
	assertFindSwears(t, sw, "Test FGH abcd fghi fgi", expected)
}

func TestAddRuleFileReadErr(t *testing.T) {
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(tmpFileName)
	os.Remove(tmpFileName)
	assertAddRuleErr(t, sw, "xxx*", "FileReadErr")
}

func TestAddRuleConflictErr(t *testing.T) {
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(tmpFileName)
	assertAddRuleErr(t, sw, "abc*", "ConflictErr")
	assertAddRuleErr(t, sw, "ab*", "ConflictErr")
	assertAddRuleErr(t, sw, "*", "ConflictErr")
}

func TestAddRuleInvalidWildcardErr(t *testing.T) {
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(tmpFileName)
	assertAddRuleErr(t, sw, "xx*x", "InvalidWildcard")
	assertAddRuleErr(t, sw, "*dd*", "InvalidWildcard")
	assertAddRuleErr(t, sw, "**x1", "InvalidWildcard")
	assertAddRuleErr(t, sw, "x2**", "InvalidWildcard")
	assertAddRuleErr(t, sw, "**", "InvalidWildcard")
}

func createSwears(tmpFilePath string) *Swears {
	config := SwearsConfig{
		DictFileName:         tmpFilePath,
		OnAddRuleFileReadErr: "FileReadErr",
		OnAddRuleConflictErr: "ConflictErr",
		OnAddRuleSaveErr:     "SaveErr",
		OnIvalidWildcardErr:  "InvalidWildcard",
	}
	sw := NewSwears(nil, config)
	sw.LoadSwears()
	return sw
}

func createTmpDict() string {
	tmpFile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
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

func assertFindSwears(t *testing.T, sw *Swears, m string, expected []string) {
	actual := sw.FindSwears(m)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected swears %#v, got %#v", expected, actual)
	}
}

func assertAddRule(t *testing.T, sw *Swears, r string) {
	err := sw.AddRule(r)
	if err != nil {
		t.Fatalf("Expected no errors when adding rule '%s', got: %v", r, err)
	}
}

func assertAddRuleErr(t *testing.T, sw *Swears, r string, expected string) {
	err := sw.AddRule(r)
	if err == nil {
		t.Fatalf("Expected error %v when adding rule '%s', got no errors", expected, r)
	} else if err.Error() != expected {
		t.Fatalf("Expected error %v when adding rule '%s', got %v", expected, r, err)
	}
}
