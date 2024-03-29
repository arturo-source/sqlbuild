package sqlbuild

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"text/template"
	"time"
)

var (
	ErrNoStruct = errors.New("provided value is not a struct")
	ErrNoId     = errors.New("need an id in the structure")
)

type ErrNoValidType struct {
	t string
}

func (e ErrNoValidType) Error() string {
	return e.t + " is not a valid type (valid types are bool, int, int8..., uint, uint8..., float32, float64, string and time.Time)"
}

type fields struct {
	namesOrdered []string
	nameValues   map[string]any
}

// newFields (from a struct) saves each field into a map[string]any, it will set the Tag `db:""` as key of the map if it is set
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
		nameValues:   make(map[string]any),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("db")
		if fieldName == "" {
			fieldName = field.Name
		}

		fields.set(fieldName, val.Field(i).Interface())
	}

	return fields
}

func (f *fields) len() int {
	return len(f.namesOrdered)
}

func (f *fields) set(k string, v any) {
	f.namesOrdered = append(f.namesOrdered, k)
	f.nameValues[k] = v
}

func (f *fields) get(k string) any {
	return f.nameValues[k]
}

// getId (from fields) finds id case insensitive inside the fields. Returns the original id key, and its value
func (f *fields) getId() (key string, value any, err error) {
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

type args map[string]any

// executeTemplate adds "sV" and "sK" funcs to avoid sql injections, and executes the query as text/template
func executeTemplate(query string, arguments args) string {
	t := template.New("query")
	t = t.Funcs(template.FuncMap{
		"sV": func(v any) any { return sanitize(v, "'") },
		"sK": func(v any) any { return sanitize(v, "\"") },
	})
	t = template.Must(t.Parse(query))

	buf := &bytes.Buffer{}
	t.Execute(buf, arguments)

	return buf.String()
}

// sanitize checks if the `thing` contains `quoteMark`, to avoid sql injections
func sanitize(thing any, quoteMark string) any {
	if thing == nil {
		return "null"
	}

	// https://codefibershq.com/blog/golang-why-nil-is-not-always-nil
	if val := reflect.ValueOf(thing); val.Kind() == reflect.Pointer && val.IsNil() {
		return "null"
	}

	if thingStr, ok := thing.(string); ok {
		if strings.Contains(thingStr, quoteMark) {
			thingStr = strings.ReplaceAll(thingStr, quoteMark, quoteMark+quoteMark)
		}

		thing = any(quoteMark + thingStr + quoteMark)
	}

	return thing
}

// getStructFromPointer unreferences the pointer until it gets a struct
func getStructFromPointer(s any) (val reflect.Value, err error) {
	val = reflect.ValueOf(s)
	kind := val.Kind()

	for kind == reflect.Pointer {
		s = val.Elem().Interface()
		val = reflect.ValueOf(s)
		kind = val.Kind()
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

// getVarType extracts the variable type and returns the corresponding sql type
func getVarType(v any) (string, error) {
	nullable := " NOT NULL"
	val := reflect.ValueOf(v)
	kind := val.Kind()

	for kind == reflect.Pointer {
		nullable = ""
		if val.IsNil() {
			kind = reflect.TypeOf(v).Elem().Kind()
			break
		}

		v = val.Elem().Interface()
		val = reflect.ValueOf(v)
		kind = val.Kind()
	}

	types := map[reflect.Kind]string{
		reflect.Bool:    "BOOLEAN",
		reflect.Int:     "INT",
		reflect.Int8:    "INT",
		reflect.Int16:   "INT",
		reflect.Int32:   "INT",
		reflect.Int64:   "INT",
		reflect.Uint:    "INT UNSIGNED",
		reflect.Uint8:   "INT UNSIGNED",
		reflect.Uint16:  "INT UNSIGNED",
		reflect.Uint32:  "INT UNSIGNED",
		reflect.Uint64:  "INT UNSIGNED",
		reflect.Float32: "DOUBLE",
		reflect.Float64: "DOUBLE",
		reflect.String:  "TEXT",
	}

	if t, ok := types[kind]; ok {
		return t + nullable, nil
	}

	if _, ok := v.(time.Time); ok {
		return "DATETIME" + nullable, nil
	}

	return "", ErrNoValidType{kind.String()}
}
