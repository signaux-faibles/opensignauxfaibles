package marshal

import "reflect"

// ExtractTableColumns extrait les noms des colonnes pour une table SQL via le tag "sql"
func ExtractTableColumns(tuple Tuple) (header []string) {
	return ExtractFieldsByTags(tuple, "sql")
}

// ExtractTableRow extrait les valeurs des colonnes pour une table SQL via le tag "sql"
func ExtractTableRow(tuple Tuple) (row []any) {
	rawValues := ExtractValuesByTags(tuple, "sql")
	for _, v := range rawValues {
		row = append(row, deref(v).Interface())
	}
	return row
}

func deref(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return val
		}
		val = val.Elem()
	}
	return val
}
