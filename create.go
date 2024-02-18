package sqlbuild

import (
	"fmt"
)

// Create creates a 'create table' query from the struct name, and sets all the fields as columns, with the specific type of variable
func Create(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)

	// TODO
	query = fmt.Sprintf("create table %s ()", sName)
	return
}
