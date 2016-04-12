package main

import (
	"reflect"
	"testing"
)

// Checks whether all fields in config-rename.json are consistent with fields in BotConfig struct.
// Due to shitty golang JSON package this cannot be asserted without test. Thanks Obama.
func TestConfigFileIntegrity(t *testing.T) {

	config := readConfig("config-rename.json")
	assertAllFieldsNotEmpty(t, config.SwearsConfig)
}

func assertAllFieldsNotEmpty(t *testing.T, obj interface{}) {
	objVal := reflect.ValueOf(obj)
	objType := objVal.Type()

	for i := 0; i < objVal.NumField(); i++ {
		fieldVal := objVal.Field(i)
		fieldName := objType.Field(i).Name
		kind := fieldVal.Kind()

		if kind == reflect.Func ||
			kind == reflect.Interface ||
			kind == reflect.Ptr ||
			kind == reflect.Struct ||
			kind == reflect.UnsafePointer {

			t.Errorf("Can't check field '%s'", fieldName)
		} else if kind == reflect.Array ||
			kind == reflect.Slice ||
			kind == reflect.Map ||
			kind == reflect.Chan {

			assertNotZeroLen(t, fieldVal, fieldName)
		} else {
			assertNotZeroVal(t, fieldVal.Interface(), fieldName)
		}
	}
}

func assertNotZeroLen(t *testing.T, val reflect.Value, fieldName string) {
	if val.Len() == 0 {
		t.Errorf("Field %s is zero length.", fieldName)
	}
}

func assertNotZeroVal(t *testing.T, in interface{}, fieldName string) {
	if in == reflect.Zero(reflect.TypeOf(in)).Interface() {
		t.Errorf("Field '%s' is zero value.", fieldName)
	}
}
