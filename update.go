package sqlbuild

// Update creates as 'update' query from the struct name, and sets all the values in the struct, only for the specific id
func Update(s any) (query string, err error) {
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

	queryTemplate := `update {{sK .tableName}} set {{range $i, $name := .namesOrdered}}{{if $i}}, {{end}}{{sK $name}} = {{sV (index $.nameValues $name)}}{{end}} WHERE {{sK .id}} = {{sV .idValue}}`
	query = executeTemplate(queryTemplate, args{"tableName": sName, "id": idName, "idValue": idValue, "nameValues": fields.nameValues, "namesOrdered": fields.namesOrdered})
	return
}
