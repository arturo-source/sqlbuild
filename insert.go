package sqlbuild

import (
	"fmt"
	"strings"
)

// Insert creates a 'insert' query from the struct name, and sets all values in the struct
func Insert(s any) (query string, err error) {
	return InsertMultiple([]any{s})
}

// InsertMultiple creates a 'insert' query from the struct name, and sets all values in the struct, for all the structs in the array
func InsertMultiple[T any](structs []T) (query string, err error) {
	sval, err := getStructFromPointer(structs[0])
	if err != nil {
		return query, err
	}

	addCommasAndParenthesis := func(things []string) string {
		withCommas := strings.Join(things, ", ")
		withParenthesis := fmt.Sprint("(", withCommas, ")")
		return withParenthesis
	}

	sName := getStructName(sval)
	fields := newFields(sval)
	fieldsNames := fields.getNames()
	keys := addCommasAndParenthesis(sanitizedSlice(fieldsNames, "\""))

	multipleValues := make([]string, 0, fields.len())
	for _, s := range structs {
		sval, err := getStructFromPointer(s)
		if err != nil {
			return query, err
		}

		fields := newFields(sval)
		values := make([]any, 0, fields.len())
		for _, fieldName := range fieldsNames {
			values = append(values, fields.get(fieldName).Interface())
		}

		multipleValues = append(multipleValues, addCommasAndParenthesis(sanitizedSlice(values, "'")))
	}
	multipleValuesStr := strings.Join(multipleValues, ", ")

	query = "insert into %s __keys__ values __multiple_values__"
	query = strings.ReplaceAll(query, "__keys__", keys)
	query = strings.ReplaceAll(query, "__multiple_values__", multipleValuesStr)
	query = sanitizeInput(query, sName)
	return
}
