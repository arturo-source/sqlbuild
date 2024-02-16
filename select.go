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
