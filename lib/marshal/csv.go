package marshal

import (
	"opensignauxfaibles/lib/misc"
	"reflect"
	"strconv"
	"time"
)

// CSVMarshaller is a handy wrapper to extract headers and values from a tuple.
// The tuple is expected to be a "struct" type with "csv" tags for values to be
// included in the csv export.
type CSVMarshaller struct {
	objectType  reflect.Type
	objectValue reflect.Value
}

func NewCSVMarshaller(tuple any) CSVMarshaller {
	return CSVMarshaller{
		objectType:  reflect.TypeOf(tuple),
		objectValue: reflect.ValueOf(tuple),
	}
}

// Headers returns the CSV headers
func (m CSVMarshaller) Headers() []string {
	return recursiveExtractTags(m.objectType, "csv")
}

// Values returns the tuple values, in same order as the header, and converted to strings
func (m CSVMarshaller) Values() (values []string) {
	rawValues := m.recursiveExtractValues(m.objectType, m.objectValue, "csv")
	for _, v := range rawValues {
		values = append(values, m.valueToCSV(v))
	}
	return values
}

func (m CSVMarshaller) valueToCSV(v reflect.Value) string {
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
		if v.Type() == reflect.TypeOf(misc.Periode{}) {
			t := v.Interface().(misc.Periode)
			return t.String()
		}
		return ""
	default:
		// Fallback to string representation
		return v.String()
	}
}
