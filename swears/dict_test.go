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

	sw := createSwears(t, tmpFileName)
	expected := []string{"a", "abcd", "abba"}
	assertFindSwears(t, sw, "Test A abCD abcde abBA abc", expected)
}

func TestAddRule(t *testing.T) {
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(t, tmpFileName)
	assertAddRule(t, sw, "Fgh*")
	expected := []string{"fgh", "abcd", "fghi"}
	assertFindSwears(t, sw, "Test FGH abcd fghi fgi", expected)
}

func TestAddRuleFileReadErr(t *testing.T) {
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(t, tmpFileName)
	os.Remove(tmpFileName)
	assertAddRuleErr(t, sw, "xxx*", DictFileReadErr)
}

func TestAddRuleConflictErr(t *testing.T) {
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(t, tmpFileName)
	assertAddRuleErr(t, sw, "abc*", AddRuleConflictErr)
	assertAddRuleErr(t, sw, "ab*", AddRuleConflictErr)
	assertAddRuleErr(t, sw, "*", AddRuleConflictErr)
}

func TestAddRuleInvalidWildcardErr(t *testing.T) {
	tmpFileName := createTmpDict()
	defer os.Remove(tmpFileName)

	sw := createSwears(t, tmpFileName)
	assertAddRuleErr(t, sw, "xx*x", InvalidWildcardErr)
	assertAddRuleErr(t, sw, "*dd*", InvalidWildcardErr)
	assertAddRuleErr(t, sw, "**x1", InvalidWildcardErr)
	assertAddRuleErr(t, sw, "x2**", InvalidWildcardErr)
	assertAddRuleErr(t, sw, "**", InvalidWildcardErr)
}

func createSwears(t *testing.T, tmpFilePath string) *Swears {
	config := SwearsConfig{
		DictFileName: tmpFilePath,
	}
	sw := NewSwears(nil, config)
	err := sw.LoadSwears()
	if err != Success {
		t.Fatalf("Expected to load dictionary without errors, got %v", err)
	}

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
	if err != Success {
		t.Fatalf("Expected no errors when adding rule '%s', got: %v", r, err)
	}
}

func assertAddRuleErr(t *testing.T, sw *Swears, r string, expected int) {
	err := sw.AddRule(r)
	if err != expected {
		t.Fatalf("Expected error %v when adding rule '%s', got %v", expected, r, err)
	}
}
