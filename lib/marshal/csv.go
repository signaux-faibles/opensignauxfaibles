package marshal

import (
	"reflect"
	"strconv"
	"time"
)

// ExtractCSVHeaders extrait les en-têtes csv via le tag "csv"
func ExtractCSVHeaders(tuple Tuple) (header []string) {
	return extractFieldsByTags(tuple, "csv")
}

// ExtractCSVRow returns the tuple values, in same order as the header, and converted to strings
func ExtractCSVRow(tuple Tuple) (values []string) {
	rawValues := extractValuesByTags(tuple, "csv")
	for _, v := range rawValues {
		values = append(values, valueToString(v))
	}
	return values
}

func valueToString(v reflect.Value) string {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return ""
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			t := v.Interface().(time.Time)
			return t.Format(time.DateOnly)
		}
		return ""
	default:
		// Fallback to string representation
		return v.String()
	}
}
