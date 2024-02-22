package sqlbuild

// Drop creates a 'drop table' query from the struct name
func Drop(s any) (query string, err error) {
	sval, err := getStructFromPointer(s)
	if err != nil {
		return query, err
	}

	sName := getStructName(sval)
	query = executeTemplate(`drop table {{sK .tableName}}`, args{"tableName": sName})
	return
}
