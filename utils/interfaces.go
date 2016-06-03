package utils

import (
	"fmt"
	"reflect"
)

func GetEmptyFieldNames(obj interface{}) []string {
	result := []string{}
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

			panic(fmt.Sprintf("Cannot check field %v", fieldName))
		} else if kind == reflect.Array ||
			kind == reflect.Slice ||
			kind == reflect.Map ||
			kind == reflect.Chan {

			if fieldVal.Len() == 0 {
				result = append(result, fieldName)
			}
		} else {
			in := fieldVal.Interface()
			if in == reflect.Zero(reflect.TypeOf(in)).Interface() {
				result = append(result, fieldName)
			}
		}
	}
	return result
}
