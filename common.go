package sqlbuild

import (
	"errors"
	"reflect"
	"strings"
)

type Fields map[string]reflect.Value

var (
	ErrNoStruct = errors.New("provided value is not a struct")
	ErrNoId     = errors.New("need an id in the structure")
)

func getStructName(s any) (sName string, err error) {
	val := reflect.ValueOf(s)
	kind := val.Kind()
	if kind == reflect.Pointer {
		return getStructName(val.Elem().Interface())
	}

	sName = reflect.TypeOf(s).Name()
	if kind != reflect.Struct {
		err = ErrNoStruct
	}

	return
}

func getStructFields(s any) Fields {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	fields := make(Fields)

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

func getIdFromFields(fields Fields) (key string, value reflect.Value, err error) {
	for fieldName, fieldValue := range fields {
		if strings.ToLower(fieldName) == "id" {
			key = fieldName
			value = fieldValue
			return
		}
	}

	err = ErrNoId
	return
}
