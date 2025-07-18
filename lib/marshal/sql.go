package marshal

import "fmt"

// ExtractTableColumns extrait les noms des colonnes pour une table SQL via le tag "sql"
func ExtractTableColumns(tuple Tuple) (header []string) {
	return ExtractFieldsByTags(tuple, "sql")
}

func ExtractTableRow(tuple Tuple) (row []any) {
	rawValues := ExtractValuesByTags(tuple, "sql")
	fmt.Println("N values", len(rawValues))
	for _, v := range rawValues {
		fmt.Printf("DEBUG: value %+v\n", v.Interface())
		row = append(row, v.Interface)
	}
	return row
}
