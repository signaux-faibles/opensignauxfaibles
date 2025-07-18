package marshal

// ExtractTableColumns extrait les noms des colonnes pour une table SQL via le tag "sql"
func ExtractTableColumns(tuple Tuple) (header []string) {
	return ExtractFieldsByTags(tuple, "sql")
}

func ExtractTableRow(tuple Tuple) (row []any) {
	rawValues := ExtractValuesByTags(tuple, "sql")
	for _, v := range rawValues {
		row = append(row, v.Interface)
	}
	return row
}
