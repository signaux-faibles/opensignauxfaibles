package marshal

import "reflect"

// ExtractTableColumns extrait les noms des colonnes pour une table SQL via le tag "sql"
func ExtractTableColumns(tuple Tuple) (header []string) {
	return ExtractFieldsByTags(tuple, "sql")
}

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
			return reflect.Zero(val.Type().Elem()) // return zero value if nil
		}
		val = val.Elem()
	}
	return val
}
