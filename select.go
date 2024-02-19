package sqlbuild

import (
	"fmt"
)

// SelectAll creates a 'select' query from the struct name
func SelectAll(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	query = fmt.Sprintf("select * from %s", sName)
	return
}

// SelectById creates a 'select' query from the struct name, but only for the specific id
func SelectById(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	fields := newFields(sval)
	idName, idValue, err := fields.getId()
	if err != nil {
		return query, err
	}

	query = sanitizeInput("select * from %s where %s = %v", sName, idName, idValue)
	return
}
