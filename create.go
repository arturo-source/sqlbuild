package sqlbuild

import (
	"fmt"
)

func CreateTable(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)

	// TODO
	query = fmt.Sprintf("create table %s ()", sName)
	return
}
