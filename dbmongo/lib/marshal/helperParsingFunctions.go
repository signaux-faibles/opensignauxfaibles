package marshal

import "strings"

var MappingDict = map[string]mappingFunc{}

type mappingFunc func(...interface{}) interface{}

func Departement(codePostal string) string {
	if len(codePostal) != 5 {
		return ""
	}
	return codePostal[0:2]
}

func StripPoint(APE string) string {
	APE = strings.Replace(APE, ".", "", -1)
	return APE
}

// Divide100 divides a float by 100
func Divide100(f *float64) *float64 {
	if f == nil {
		return nil
	}
	*f = *f / 100
	return f
}
