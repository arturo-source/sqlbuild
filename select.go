package sqlbuild

import (
	"fmt"
)

func SelectAll(s any) (query string, err error) {
	sName, err := getStructName(s)
	if err != nil {
		return query, err
	}

	query = fmt.Sprintf("select * from %s", sName)
	return
}

func SelectById(s any) (query string, err error) {
	sName, err := getStructName(s)
	if err != nil {
		return query, err
	}

	fields := getStructFields(s)
	idName, idValue, err := getIdFromFields(fields)
	if err != nil {
		return query, err
	}

	query = sanitizeInput("select * from %s where %s = %v", sName, idName, idValue)
	return
}
