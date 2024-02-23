package sqlbuild

// DeleteAll creates a 'delete' query from the struct name
func DeleteAll(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	query = executeTemplate(`DELETE FROM {{sK .tableName}}`, args{"tableName": sName})
	return
}

// DeleteById creates a 'delete' query from the struct name, but only for the specific id
func DeleteById(s any) (query string, err error) {
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

	queryTemplate := `DELETE FROM {{sK .tableName}} WHERE {{sK .id}} = {{sV .idValue}}`
	query = executeTemplate(queryTemplate, args{"tableName": sName, "id": idName, "idValue": idValue})
	return
}
