package sqlbuild

// SelectAll creates a 'select' query from the struct name
func SelectAll(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	query = executeTemplate(`select * from {{s .tableName "\""}}`, args{"tableName": sName})
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

	query = executeTemplate(
		`select * from {{s .tableName "\""}} where {{s .id "\""}} = {{s .idValue "'"}}`,
		args{"tableName": sName, "id": idName, "idValue": idValue})
	return
}
