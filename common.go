package sqlbuild

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrNoStruct = errors.New("provided value is not a struct")
	ErrNoId     = errors.New("need an id in the structure")
)

type fields struct {
	namesOrdered []string
	nameValues   map[string]reflect.Value
}

// newFields (from a struct) saves each field into a map[string]reflect.Value, it will set the Tag `db:""` as key of the map if it is set
//
// Example:
//
//	type Person struct {
//		Name string `db:"name"`
//		Age  int
//	}
func newFields(val reflect.Value) fields {
	t := reflect.TypeOf(val.Interface())
	fields := fields{
		namesOrdered: make([]string, 0, t.NumField()),
		nameValues:   make(map[string]reflect.Value),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("db")
		if fieldName == "" {
			fieldName = field.Name
		}

		fields.set(fieldName, val.Field(i))
	}

	return fields
}

func (f *fields) len() int {
	return len(f.namesOrdered)
}

func (f *fields) set(k string, v reflect.Value) {
	f.namesOrdered = append(f.namesOrdered, k)
	f.nameValues[k] = v
}

func (f *fields) get(k string) reflect.Value {
	return f.nameValues[k]
}

func (f *fields) getNames() []string {
	return f.namesOrdered
}

// getId (from fields) finds id case insensitive inside the fields. Returns the original id key, and its value
func (f *fields) getId() (key string, value reflect.Value, err error) {
	for fieldName, fieldValue := range f.nameValues {
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
		args[i] = reflect.ValueOf(sanitize(args[i], "'")) // <-- Change
	}

	return fmt.Sprintf(query, args...)
}

func sanitize(thing any, quoteMark string) string {
	if thingStr, ok := any(thing).(string); ok {
		if strings.Contains(thingStr, quoteMark) {
			thingStr = strings.ReplaceAll(thingStr, quoteMark, quoteMark+quoteMark)
		}

		return fmt.Sprint(quoteMark, thingStr, quoteMark)
	}

	return fmt.Sprint(thing)
}

func sanitizedSlice[T any](things []T, quoteMark string) []string {
	thingsSanitized := make([]string, 0, len(things))
	for i := range things {
		thingsSanitized = append(thingsSanitized, sanitize(things[i], quoteMark))
	}

	return thingsSanitized
}

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
