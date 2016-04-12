package utils

import (
	"testing"
)

func TestParamFormat(t *testing.T) {
	params := map[string]string{
		"Param1": "value 1",
		"Param2": "value 2",
	}

	assertEq(t, ParamFormat("{Param1}", params), "value 1")
	assertEq(t, ParamFormat("Test {Param1}", params), "Test value 1")
	assertEq(t, ParamFormat("Test {Param1}!", params), "Test value 1!")
	assertEq(t, ParamFormat("Test {Param1} {Param1}!", params), "Test value 1 value 1!")
	assertEq(t, ParamFormat("Test {param1}", params), "Test {param1}")
	assertEq(t, ParamFormat("{Param2}{Param1}", params), "value 2value 1")
	assertEq(t, ParamFormat("{{Param1}", params), "{value 1")
	assertEq(t, ParamFormat("{Param1}}", params), "value 1}")
	assertEq(t, ParamFormat("{{Param1}}", params), "{value 1}")
}

func assertEq(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}
