package utils

import (
	"regexp"
	"strings"
)

var paramRegexp *regexp.Regexp = regexp.MustCompile("\\{[a-zA-Z0-9]+\\}")

func ParamFormat(format string, params map[string]string) string {
	paramNames := paramRegexp.FindAllString(format, -1)
	for _, paramName := range paramNames {
		paramValue, ok := params[strings.Trim(paramName, "{}")]
		if ok {
			format = strings.Replace(format, paramName, paramValue, -1)
		}
	}
	return format
}

func ContainsCaseIns(needle string, haystack []string) bool {
	needle = strings.ToLower(needle)
	for _, value := range haystack {
		if needle == strings.ToLower(value) {
			return true
		}
	}
	return false
}
