package marshal

// ExtractTableColumns extrait les noms des colonnes pour une table SQL via le tag "sql"
func ExtractTableColumns(tuple Tuple) (header []string) {
	return extractFieldsByTags(tuple, "sql")
}

// ExtractTableRow extrait les valeurs des colonnes pour une table SQL via le tag "sql"
func ExtractTableRow(tuple Tuple) (row []any) {
	rawValues := extractValuesByTags(tuple, "sql")
	for _, v := range rawValues {
		row = append(row, v.Interface())
	}
	return row
}
