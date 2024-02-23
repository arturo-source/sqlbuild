package sqlbuild

// SelectAll creates a 'select' query from the struct name
func SelectAll(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	query = executeTemplate(`SELECT * FROM {{sK .tableName}}`, args{"tableName": sName})
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

	queryTemplate := `SELECT * FROM {{sK .tableName}} WHERE {{sK .id}} = {{sV .idValue}}`
	query = executeTemplate(queryTemplate, args{"tableName": sName, "id": idName, "idValue": idValue})
	return
}
