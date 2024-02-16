package sqlbuild

import (
	"fmt"
)

func CreateTable(s any) (query string, err error) {
	sName, err := getStructName(s)
	if err != nil {
		return query, err
	}

	// TODO
	query = fmt.Sprintf("create table %s ()", sName)
	return
}
