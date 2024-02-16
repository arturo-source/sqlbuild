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
