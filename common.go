package sqlbuild

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Fields map[string]reflect.Value

var (
	ErrNoStruct = errors.New("provided value is not a struct")
	ErrNoId     = errors.New("need an id in the structure")
)

// getStructFromPointer unreferences the pointer until it gets a struct
func getStructFromPointer(s any) (val reflect.Value, err error) {
	val = reflect.ValueOf(s)

	kind := val.Kind()
	if kind == reflect.Pointer {
		return getStructFromPointer(val.Elem().Interface())
	}

	if kind != reflect.Struct {
		err = ErrNoStruct
	}

	return
}

// getStructName returns the name of the struct as a string
func getStructName(val reflect.Value) string {
	return reflect.TypeOf(val.Interface()).Name()
}

// getStructFields expects any struct and saves each field into a map[string]reflect.Value, it will set the Tag `db:""` as key of the map if it is set
//
// Example:
//
//	type Person struct {
//		Name string `db:"name"`
//		Age  int
//	}
func getStructFields(val reflect.Value) Fields {
	t := reflect.TypeOf(val.Interface())
	fields := make(Fields)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("db")
		if fieldName == "" {
			fieldName = field.Name
		}

		fields[fieldName] = val.Field(i)
	}

	return fields
}

// getIdFromFields finds id case insensitive inside the fields. Returns the original id key, and its value
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

// sanitizeInput format like fmt.Sprintf, but sanitizes args to avoid sql injections
func sanitizeInput(query string, args ...any) string {
	for i := range args {
		if argStr, ok := args[i].(string); ok {
			argStr = strings.ReplaceAll(argStr, "'", "''")
			argStr = fmt.Sprint("'", argStr, "'")
			args[i] = reflect.ValueOf(argStr)
		}
	}

	return fmt.Sprintf(query, args...)
}
