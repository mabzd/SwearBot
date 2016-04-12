package swears

import (
	"../utils"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestReadEmptySettings(t *testing.T) {
	fileName := createTmpSettingsPath(t)
	defer os.Remove(fileName)

	sw := createSettings(fileName)
	expected := &AllSettings{UserSettings: make(map[string][]*UserSettings)}
	assertReadSettings(t, sw, expected)
}

func TestWriteReadSettings(t *testing.T) {
	fileName := createTmpSettingsPath(t)
	defer os.Remove(fileName)

	settings := &AllSettings{
		UserSettings: map[string][]*UserSettings{
			"user1": []*UserSettings{
				&UserSettings{
					UserId:   "user1",
					Channel:  "#channel1",
					Settings: map[string]string{"key": "value1"},
				},
				&UserSettings{
					UserId:   "user1",
					Channel:  "#channel2",
					Settings: map[string]string{"key": "value2"},
				},
			},
			"user2": []*UserSettings{
				&UserSettings{
					UserId:   "user2",
					Channel:  "#channel2",
					Settings: map[string]string{"key": "value2"},
				},
			},
		},
	}

	sw := createSettings(fileName)
	assertWriteSettings(t, sw, settings)
	assertReadSettings(t, sw, settings)
}

func TestSettingSettings(t *testing.T) {
	s := &AllSettings{UserSettings: make(map[string][]*UserSettings)}

	s.SetSetting("u1", "c1", "key1", "value1")
	s.SetSetting("u1", "c1", "key2", "value2")
	s.SetSetting("u1", "c2", "key3", "value3")
	s.SetSetting("u1", "c1", "key2", "value2-changed")
	s.SetSetting("u2", "c1", "key4", "value4")

	assertGetSetting(t, s, "u1", "c1", "key1", "value1")
	assertGetSetting(t, s, "u1", "c1", "key2", "value2-changed")
	assertGetSetting(t, s, "u1", "c2", "key3", "value3")
	assertGetSetting(t, s, "u2", "c1", "key4", "value4")
	assertNotGetString(t, s, "u1", "c1", "key3")
	assertNotGetString(t, s, "u1", "c1", "key4")
}

func createSettings(tmpFilePath string) *Swears {
	config := SwearsConfig{
		SettingsFileName: tmpFilePath,
	}
	return NewSwears(nil, config)
}

func createTmpSettingsPath(t *testing.T) string {
	fileName := utils.CreateTmpFileName("Settings")
	if fileName == "" {
		t.Fatal("Cannot create temp settings file path")
	}
	return fileName
}

func assertReadSettings(t *testing.T, sw *Swears, expected *AllSettings) {
	settings, err := sw.ReadSettings()
	if err != Success {
		t.Fatalf("Expected no errors when reading settings, got %v", err)
	} else if !reflect.DeepEqual(settings, expected) {
		t.Fatalf("Settings deep equal failed. Expected %#v, got %#v", expected, settings)
	}
}

func assertWriteSettings(t *testing.T, sw *Swears, settings *AllSettings) {
	err := sw.WriteSettings(settings)
	if err != Success {
		t.Fatalf("Expected no errors when writing settings, got %v", err)
	}
}

func assertGetSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	channel string,
	key string,
	expected string) {

	value, ok := settings.GetSetting(userId, channel, key)
	if !ok || value != expected {
		if !ok {
			value = "nothing"
		} else {
			value = fmt.Sprintf("'%s'", value)
		}

		t.Fatalf(
			"Expected value '%s' for userId='%s', channel='%s', key='%s', got %s.",
			expected,
			userId,
			channel,
			key,
			value)
	}
}

func assertNotGetString(
	t *testing.T,
	settings *AllSettings,
	userId string,
	channel string,
	key string) {

	value, ok := settings.GetSetting(userId, channel, key)
	if ok {
		t.Fatalf(
			"Expected no value for userId='%s', channel='%s', key='%s', got '%s'.",
			userId,
			channel,
			key,
			value)
	}
}
