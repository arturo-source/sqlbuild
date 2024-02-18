package sqlbuild

import (
	"strings"
)

// Update creates as 'update' query from the struct name, and sets all the values in the struct, only for the specific id
func Update(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	fields := getStructFields(sval)
	idName, idValue, err := getIdFromFields(fields)
	if err != nil {
		return query, err
	}

	keyValues := make([]string, 0, len(fields))
	for k, v := range fields {
		keyValues = append(keyValues, sanitizeInput("%s = %v", k, v))
	}
	keyValuesStr := strings.Join(keyValues, ", ")

	query = "update %s set __keys_values__ where %s = %v"
	query = strings.ReplaceAll(query, "__keys_values__", keyValuesStr)
	query = sanitizeInput(query, sName, idName, idValue)
	return
}
