package sqlbuild

// Create creates a 'create table' query from the struct name, and sets all the fields as columns, with the specific type of variable.
//
// Create also asigns "id" as primary key, and no-pointers as "not null"
func Create(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	fields := newFields(sval)
	idName, _, _ := fields.getId()

	varTypes := make([]string, 0, fields.len())
	for _, t := range fields.namesOrdered {
		tValid, err := getVarType(fields.get(t))
		if err != nil {
			return query, err
		}

		varTypes = append(varTypes, tValid)
	}

	queryTemplate := `create table {{sK .tableName}} ({{range $i, $name := .namesOrdered}}{{if $i}}, {{end}}{{sK $name}} {{index $.varTypes $i}}{{if eq $name $.idName}} auto_increment primary key{{end}}{{end}})`
	query = executeTemplate(queryTemplate, args{"tableName": sName, "idName": idName, "namesOrdered": fields.namesOrdered, "varTypes": varTypes})
	return
}
