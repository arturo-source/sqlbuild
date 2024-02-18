package sqlbuild

import "fmt"

// DeleteAll creates a 'delete' query from the struct name
func DeleteAll(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	query = fmt.Sprintf("delete from %s", sName)
	return
}

// DeleteById creates a 'delete' query from the struct name, but only for the specific id
func DeleteById(s any) (query string, err error) {
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

	query = sanitizeInput("delete from %s where %s = %v", sName, idName, idValue)
	return
}
