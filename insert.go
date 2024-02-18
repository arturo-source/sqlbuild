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
func InsertMultiple(structs []any) (query string, err error) {
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
	fields := getStructFields(sval)
	fieldsNames := fields.GetNames()
	keys := addCommasAndParenthesis(fieldsNames)

	multipleValues := make([]string, 0, len(fields))
	for _, s := range structs {
		sval, err := getStructFromPointer(s)
		if err != nil {
			return query, err
		}

		fields := getStructFields(sval)
		values := make([]string, 0, len(fields))

		for _, fieldName := range fieldsNames {
			value := sanitizeInput("%v", fields[fieldName].Interface())
			values = append(values, value)
		}

		multipleValues = append(multipleValues, addCommasAndParenthesis(values))
	}
	multipleValuesStr := strings.Join(multipleValues, ", ")

	query = "insert into %s __keys__ values __multiple_values__"
	query = strings.ReplaceAll(query, "__keys__", keys)
	query = strings.ReplaceAll(query, "__multiple_values__", multipleValuesStr)
	query = sanitizeInput(query, sName)
	return
}
