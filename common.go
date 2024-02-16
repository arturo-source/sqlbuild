package sqlbuild

import (
	"errors"
	"reflect"
)

var (
	ErrNoStruct = errors.New("provided value is not a struct")
)

func getStructName(s any) (sName string, err error) {
	val := reflect.ValueOf(s)
	kind := val.Kind()
	if kind == reflect.Pointer {
		return getStructName(val.Elem().Interface())
	}

	sName = reflect.TypeOf(s).Name()
	if kind == reflect.Struct {
		err = ErrNoStruct
	}

	return
}

func getStructFields(s any) map[string]reflect.Value {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	fields := make(map[string]reflect.Value)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("db")
		if fieldName == "" {
			fieldName = field.Name
		}

		fields[fieldName] = v.Field(i)
	}

	return fields
}
