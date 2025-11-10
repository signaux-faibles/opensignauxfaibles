package parsing

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type HeaderIndexer struct {
	// Instance of destination type (struct)
	// If provided, indexing will validate that properties tagged with `input` have the tag
	// value in the headers
	// If nil, no validation will be performed
	Dest any
}

// Index validates then indexes the found columns
// en en-tête d'un fichier csv, à partir des noms de colonnes spécifiés dans le
// tag "input"  annotant les propriétés du type de destination du parseur.
//
// Si aucun type n'est précisé, les en-têtes ne sont pas validées.
func (i HeaderIndexer) Index(headers []string, caseSensitive bool) (ColIndex, error) {
	var colIndex = newColIndex(caseSensitive)

	for pos, name := range headers {
		if !caseSensitive {
			name = strings.ToLower(name)
		}
		colIndex.mapping[name] = pos
	}

	var err error

	if i.Dest != nil {
		requiredFields := ExtractInputHeaders(i.Dest)
		_, err = colIndex.HasFields(requiredFields)
	}

	return colIndex, err
}

// ColIndex provides the index of each column.
type ColIndex struct {
	mapping       map[string]int
	caseSensitive bool
}

// newColIndex initializes an empty ColIndex
func newColIndex(caseSensitive bool) ColIndex {
	c := ColIndex{}
	c.mapping = make(map[string]int)
	c.caseSensitive = caseSensitive
	return c
}

func (idx ColIndex) Get(name string) (int, bool) {
	if !idx.caseSensitive {
		name = strings.ToLower(name)
	}
	pos, found := idx.mapping[name]
	return pos, found
}

// HasFields checks the presence of a set of columns.
func (idx ColIndex) HasFields(requiredFields []string) (bool, error) {
	for _, name := range requiredFields {
		if _, found := idx.Get(name); !found {
			return false, errors.New("column " + name + " not found, aborting")
		}
	}
	return true, nil
}

// IndexRow returns a structure to facilitate data reading.
func (idx ColIndex) IndexRow(row []string) IndexedRow {
	return IndexedRow{idx, row}
}

// IndexedRow facilitates reading data by columns, within a row.
type IndexedRow struct {
	idx ColIndex
	row []string
}

// GetVal returns the value associated with the given column, on the current row.
// If the column does not exist, a fatal error is triggered.
func (indexedRow IndexedRow) GetVal(colName string) string {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		log.Fatal("Column not found in ColMapping: " + colName)
	}
	return indexedRow.row[index]
}

// GetOptionalVal returns the value associated with the given column, on the current row.
// If the column does not exist, the boolean will be false.
func (indexedRow IndexedRow) GetOptionalVal(colName string) (string, bool) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return "", false
	}
	return indexedRow.row[index], true
}

// GetFloat64 returns the decimal value associated with the given column, on the current row.
// A nil pointer is returned if the column does not exist or the value is an empty string.
func (indexedRow IndexedRow) GetFloat64(colName string) (*float64, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return nil, fmt.Errorf("GetFloat64 failed to find column: %v", colName)
	}
	return ParsePFloat(indexedRow.row[index])
}

// GetCommaFloat64 returns the decimal value with comma associated with the given column, on the current row.
// A nil pointer is returned if the column does not exist or the value is an empty string.
func (indexedRow IndexedRow) GetCommaFloat64(colName string) (*float64, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return nil, fmt.Errorf("GetCommaFloat64 failed to find column: %v", colName)
	}
	normalizedDecimalVal := strings.Replace(indexedRow.row[index], ",", ".", -1)
	return ParsePFloat(normalizedDecimalVal)
}

// GetInt returns the integer value associated with the given column, on the current row.
// A nil pointer is returned if the column does not exist or the value is an empty string.
func (indexedRow IndexedRow) GetInt(colName string) (*int, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return nil, fmt.Errorf("GetInt failed to find column: %v", colName)
	}
	return ParsePInt(indexedRow.row[index])
}

// GetIntFromFloat returns the integer value associated with the given column, on the current row.
// A nil pointer is returned if the column does not exist or the value is an empty string.
func (indexedRow IndexedRow) GetIntFromFloat(colName string) (*int, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return nil, fmt.Errorf("GetIntFromFloat failed to find column: %v", colName)
	}
	return ParsePIntFromFloat(indexedRow.row[index])
}

// GetBool returns the boolean value associated with the given column, on the current row.
// A nil pointer is returned if the column does not exist or the value is an empty string.
func (indexedRow IndexedRow) GetBool(colName string) (bool, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return false, fmt.Errorf("GetBool failed to find column: %v", colName)
	}
	return strconv.ParseBool(indexedRow.row[index])
}

// LowercaseFields normalizes column names to lowercase.
func LowercaseFields(headerFields []string) []string {
	var normalizedHeaderFields = make([]string, len(headerFields))
	for i, name := range headerFields {
		normalizedHeaderFields[i] = strings.ToLower(name)
	}
	return normalizedHeaderFields
}

// ExtractInputHeaders extracts input file headers via the "input" tag
func ExtractInputHeaders(object any) []string {
	return ExtractFieldsByTags(object, "input")
}

// ExtractFieldsByTags extrait les noms de colonnes depuis les valeurs du tag
// "tag" de chaque propriété de l'objet (de type "struct") fourni.
//
// Pour une propriété de type "struct", et en l'absence de ce tag, les noms des
// colonnes seront recursivement extraites du type de la propriété.
func ExtractFieldsByTags(object any, tag string) (fields []string) {
	t := reflect.TypeOf(object)
	for i := range t.NumField() {
		field := t.Field(i)

		fieldTag := field.Tag.Get(tag)

		if fieldTag != "" {
			fields = append(fields, fieldTag)
		}
	}
	return fields
}

// ExtractValuesByTags returns the values associated with a tag
// The values are in same order than the output of ExtractFieldsByTags
func ExtractValuesByTags(object any, tag string) (values []reflect.Value) {

	t := reflect.TypeOf(object)
	v := reflect.ValueOf(object)

	for i := range t.NumField() {
		field := t.Field(i)
		fieldValue := v.Field(i)
		tagWithSingleVal := field.Tag.Get(tag)

		if tagWithSingleVal != "" {
			values = append(values, fieldValue)
		}
	}
	return values
}
