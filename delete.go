package sqlbuild

import "fmt"

func DeleteAll(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	query = fmt.Sprintf("delete from %s", sName)
	return
}

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
