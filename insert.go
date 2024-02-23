package sqlbuild

// Insert creates a 'insert' query from the struct name, and sets all values in the struct
func Insert(s any) (query string, err error) {
	return InsertMultiple([]any{s})
}

// InsertMultiple creates a 'insert' query from the struct name, and sets all values in the struct, for all the structs in the array
func InsertMultiple[T any](structs []T) (query string, err error) {
	sval, err := getStructFromPointer(structs[0])
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	fields := newFields(sval)

	multipleFieldsValues := make([][]any, 0, len(structs))
	for _, s := range structs {
		sval, err := getStructFromPointer(s)
		if err != nil {
			return query, err
		}

		fields := newFields(sval)
		fieldsValues := make([]any, 0, fields.len())
		for _, fieldName := range fields.namesOrdered {
			fieldsValues = append(fieldsValues, fields.get(fieldName))
		}

		multipleFieldsValues = append(multipleFieldsValues, fieldsValues)
	}

	queryTemplate := `INSERT INTO {{sK .tableName}} ({{range $i, $fName := .fieldsNames}}{{if $i}}, {{end}}{{sK $fName}}{{end}}) VALUES {{range $i, $fieldsValues := .multipleFieldsValues}}{{if $i}}, {{end}}({{range $j, $value := $fieldsValues}}{{if $j}}, {{end}}{{sV $value}}{{end}}){{end}}`
	query = executeTemplate(queryTemplate, args{"tableName": sName, "fieldsNames": fields.namesOrdered, "multipleFieldsValues": multipleFieldsValues})
	return
}
