package settings

import (
	"../utils"
	"os"
	"reflect"
	"testing"
)

func TestReadEmptySettings(t *testing.T) {
	fileName := createTmpSettingsPath(t)
	defer os.Remove(fileName)

	expected := &AllSettings{
		UserSettings: map[string]*UserSettings{},
		ChanSettings: map[string]*ChanSettings{},
		Settings:     map[string]string{},
	}

	assertLoadSettings(t, fileName, expected)
}

func TestWriteReadSettings(t *testing.T) {
	fileName := createTmpSettingsPath(t)
	defer os.Remove(fileName)

	settings := &AllSettings{
		UserSettings: map[string]*UserSettings{
			"user1": &UserSettings{
				UserId: "user1",
				ChanSettings: map[string]*ChanSettings{
					"chan1": &ChanSettings{
						ChannelId: "chan1",
						Settings:  map[string]string{"key1": "val1", "key2": "val2"},
					},
					"chan2": &ChanSettings{
						ChannelId: "chan2",
						Settings:  map[string]string{},
					},
				},
				Settings: map[string]string{"key4": "val4", "key5": "val5"},
			},
			"user2": &UserSettings{
				UserId:       "user2",
				ChanSettings: map[string]*ChanSettings{},
				Settings:     map[string]string{},
			},
		},
		ChanSettings: map[string]*ChanSettings{
			"chan1": &ChanSettings{
				ChannelId: "chan1",
				Settings:  map[string]string{"key6": "val6", "key7": "val7"},
			},
			"chan2": &ChanSettings{
				ChannelId: "chan2",
				Settings:  map[string]string{},
			},
		},
		Settings: map[string]string{"key8": "val8", "key9": "val9"},
	}

	assertSaveSettings(t, fileName, settings)
	assertLoadSettings(t, fileName, settings)
}

func TestSettingSettings(t *testing.T) {
	settings := NewSettings()
	settings.SetSetting("k1", "v1")
	settings.SetSetting("k2", "v2")
	settings.SetSetting("k1", "v1u")

	assertGetSetting(t, settings, "k1", "v1u")
	assertGetSetting(t, settings, "k2", "v2")
	assertNotGetSetting(t, settings, "k3")

	assertRemoveSetting(t, settings, "k1")
	assertNotGetSetting(t, settings, "k1")
	assertGetSetting(t, settings, "k2", "v2")
	assertNotRemoveSetting(t, settings, "k3")
}

func TestSettingUserSettings(t *testing.T) {
	settings := NewSettings()
	settings.SetUserSetting("u1", "k1", "u1v1")
	settings.SetUserSetting("u1", "k2", "u1v2")
	settings.SetUserSetting("u1", "k1", "u1v1u")
	settings.SetUserSetting("u2", "k1", "u2v1")

	assertGetUserSetting(t, settings, "u1", "k1", "u1v1u")
	assertGetUserSetting(t, settings, "u1", "k2", "u1v2")
	assertGetUserSetting(t, settings, "u2", "k1", "u2v1")
	assertNotGetUserSetting(t, settings, "u2", "k2")
	assertNotGetUserSetting(t, settings, "u3", "k1")

	assertRemoveUserSetting(t, settings, "u1", "k1")
	assertNotGetUserSetting(t, settings, "u1", "k1")
	assertGetUserSetting(t, settings, "u1", "k2", "u1v2")
	assertGetUserSetting(t, settings, "u2", "k1", "u2v1")
	assertNotRemoveUserSetting(t, settings, "u1", "k3")
	assertNotRemoveUserSetting(t, settings, "u3", "k1")
}

func TestSettingChanSettings(t *testing.T) {
	settings := NewSettings()
	settings.SetChanSetting("c1", "k1", "c1v1")
	settings.SetChanSetting("c1", "k2", "c1v2")
	settings.SetChanSetting("c1", "k1", "c1v1u")
	settings.SetChanSetting("c2", "k1", "c2v1")

	assertGetChanSetting(t, settings, "c1", "k1", "c1v1u")
	assertGetChanSetting(t, settings, "c1", "k2", "c1v2")
	assertGetChanSetting(t, settings, "c2", "k1", "c2v1")
	assertNotGetChanSetting(t, settings, "c2", "k2")
	assertNotGetChanSetting(t, settings, "c3", "k1")

	assertRemoveChanSetting(t, settings, "c1", "k1")
	assertNotGetChanSetting(t, settings, "c1", "k1")
	assertGetChanSetting(t, settings, "c1", "k2", "c1v2")
	assertGetChanSetting(t, settings, "c2", "k1", "c2v1")
	assertNotRemoveChanSetting(t, settings, "c1", "k3")
	assertNotRemoveChanSetting(t, settings, "c3", "k1")
}

func TestSettingUserChanSettings(t *testing.T) {
	settings := NewSettings()
	settings.SetUserChanSetting("u1", "c1", "k1", "u1c1v1")
	settings.SetUserChanSetting("u1", "c1", "k2", "u1c1v2")
	settings.SetUserChanSetting("u1", "c1", "k1", "u1c1v1u")
	settings.SetUserChanSetting("u1", "c2", "k1", "u1c2v1")
	settings.SetUserChanSetting("u2", "c1", "k1", "u2c1v1")

	assertGetUserChanSetting(t, settings, "u1", "c1", "k1", "u1c1v1u")
	assertGetUserChanSetting(t, settings, "u1", "c1", "k2", "u1c1v2")
	assertGetUserChanSetting(t, settings, "u1", "c2", "k1", "u1c2v1")
	assertGetUserChanSetting(t, settings, "u2", "c1", "k1", "u2c1v1")
	assertNotGetUserChanSetting(t, settings, "u1", "c2", "k2")
	assertNotGetUserChanSetting(t, settings, "u2", "c1", "k2")
	assertNotGetUserChanSetting(t, settings, "u3", "c1", "k1")

	assertRemoveUserChanSetting(t, settings, "u1", "c1", "k1")
	assertNotGetUserChanSetting(t, settings, "u1", "c1", "k1")
	assertGetUserChanSetting(t, settings, "u1", "c1", "k2", "u1c1v2")
	assertGetUserChanSetting(t, settings, "u1", "c2", "k1", "u1c2v1")
	assertGetUserChanSetting(t, settings, "u2", "c1", "k1", "u2c1v1")
	assertNotRemoveUserChanSetting(t, settings, "u1", "c1", "k3")
	assertNotRemoveUserChanSetting(t, settings, "u1", "c3", "k1")
	assertNotRemoveUserChanSetting(t, settings, "u3", "c1", "k1")
}

func createTmpSettingsPath(t *testing.T) string {
	fileName := utils.CreateTmpFileName("Settings")
	if fileName == "" {
		t.Fatal("Cannot create temp settings file path")
	}
	return fileName
}

func assertLoadSettings(t *testing.T, fileName string, expected *AllSettings) {
	actual := NewSettings()
	err := actual.Load(fileName)
	if err != Success {
		t.Fatalf("Expected no errors when reading settings, got %v", err)
	} else if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Settings deep equal failed.")
	}
}

func assertSaveSettings(t *testing.T, fileName string, settings *AllSettings) {
	err := settings.Save(fileName)
	if err != Success {
		t.Fatalf("Expected no errors when writing settings, got %v", err)
	}
}

func assertGetSetting(
	t *testing.T,
	settings *AllSettings,
	key string,
	expected string) {

	actual, ok := settings.GetSetting(key)
	if !ok {
		t.Fatalf("No setting key '%v' found", key)
	}
	if actual != expected {
		t.Fatalf("Setting key '%v': expected value '%v', got '%v'", key, expected, actual)
	}
}

func assertNotGetSetting(
	t *testing.T,
	settings *AllSettings,
	key string) {

	actual, ok := settings.GetSetting(key)
	if ok {
		t.Fatalf("Setting key '%v': expected no setting value, got '%v'", key, actual)
	}
}

func assertGetUserSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	key string,
	expected string) {

	actual, ok := settings.GetUserSetting(userId, key)
	if !ok {
		t.Fatalf("No user '%v' setting key '%v' found", userId, key)
	}
	if actual != expected {
		t.Fatalf(
			"User '%v' setting key '%v': expected value '%v', got '%v'",
			userId,
			key,
			expected,
			actual)
	}
}

func assertNotGetUserSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	key string) {

	actual, ok := settings.GetUserSetting(userId, key)
	if ok {
		t.Fatalf(
			"User '%v' setting key '%v': expected no setting value, got '%v'",
			userId,
			key,
			actual)
	}
}

func assertGetChanSetting(
	t *testing.T,
	settings *AllSettings,
	channelId string,
	key string,
	expected string) {

	actual, ok := settings.GetChanSetting(channelId, key)
	if !ok {
		t.Fatalf("No channel '%v' setting key '%v' found", channelId, key)
	}
	if actual != expected {
		t.Fatalf(
			"Channel '%v' setting key '%v': expected value '%v', got '%v'",
			channelId,
			key,
			expected,
			actual)
	}
}

func assertNotGetChanSetting(
	t *testing.T,
	settings *AllSettings,
	channelId string,
	key string) {

	actual, ok := settings.GetChanSetting(channelId, key)
	if ok {
		t.Fatalf(
			"Channel '%v' setting key '%v': expected no setting value, got '%v'",
			channelId,
			key,
			actual)
	}
}

func assertGetUserChanSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	channelId string,
	key string,
	expected string) {

	actual, ok := settings.GetUserChanSetting(userId, channelId, key)
	if !ok {
		t.Fatalf("No user '%v' channel '%v' setting key '%v' found", userId, channelId, key)
	}
	if actual != expected {
		t.Fatalf(
			"User '%v' channel '%v' setting key '%v': expected value '%v', got '%v'",
			userId,
			channelId,
			key,
			expected,
			actual)
	}
}

func assertNotGetUserChanSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	channelId string,
	key string) {

	actual, ok := settings.GetUserChanSetting(userId, channelId, key)
	if ok {
		t.Fatalf(
			"User '%v' channel '%v' setting key '%v': expected no setting value, got '%v'",
			userId,
			channelId,
			key,
			actual)
	}
}

func assertRemoveSetting(t *testing.T, settings *AllSettings, key string) {
	if !settings.RemoveSetting(key) {
		t.Fatalf("Setting key '%s': key expected to exist", key)
	}
}

func assertNotRemoveSetting(t *testing.T, settings *AllSettings, key string) {
	if settings.RemoveSetting(key) {
		t.Fatalf("Setting key '%s': key expected to not exist", key)
	}
}

func assertRemoveUserSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	key string) {

	if !settings.RemoveUserSetting(userId, key) {
		t.Fatalf("User '%s' setting key '%s': key expected to exist", userId, key)
	}
}

func assertNotRemoveUserSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	key string) {

	if settings.RemoveUserSetting(userId, key) {
		t.Fatalf("User '%s' setting key '%s': key expected to not exist", userId, key)
	}
}

func assertRemoveChanSetting(
	t *testing.T,
	settings *AllSettings,
	channelId string,
	key string) {

	if !settings.RemoveChanSetting(channelId, key) {
		t.Fatalf("Channel '%s' setting key '%s': key expected to exist", channelId, key)
	}
}

func assertNotRemoveChanSetting(
	t *testing.T,
	settings *AllSettings,
	channelId string,
	key string) {

	if settings.RemoveChanSetting(channelId, key) {
		t.Fatalf("Channel '%s' setting key '%s': key expected to not exist", channelId, key)
	}
}

func assertRemoveUserChanSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	channelId string,
	key string) {

	if !settings.RemoveUserChanSetting(userId, channelId, key) {
		t.Fatalf(
			"User '%s' channel '%s' setting key '%s': key expected to exist",
			userId,
			channelId,
			key)
	}
}

func assertNotRemoveUserChanSetting(
	t *testing.T,
	settings *AllSettings,
	userId string,
	channelId string,
	key string) {

	if settings.RemoveUserChanSetting(userId, channelId, key) {
		t.Fatalf(
			"User '%s' channel '%s' setting key '%s': key expected to not exist",
			userId,
			channelId,
			key)
	}
}
