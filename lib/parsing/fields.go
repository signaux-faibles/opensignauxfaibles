package parsing

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"opensignauxfaibles/lib/misc"
)

type HeaderIndexer struct {
	// Instance of destination type (struct)
	// If provided, indexing will validate that properties tagged with `input` have the tag
	// value in the headers
	// If nil, no validation will be performed
	Dest any
}

// IndexColumnsFromCsvHeader valide puis indexe les colonnes trouvées
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

// ColIndex fournit l'indice de chaque colonne.
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

// HasFields vérifie la présence d'un ensemble de colonnes.
func (idx ColIndex) HasFields(requiredFields []string) (bool, error) {
	for _, name := range requiredFields {
		if _, found := idx.Get(name); !found {
			return false, errors.New("Colonne " + name + " non trouvée. Abandon.")
		}
	}
	return true, nil
}

// IndexRow retourne une structure pour faciliter la lecture de données.
func (idx ColIndex) IndexRow(row []string) IndexedRow {
	return IndexedRow{idx, row}
}

// IndexedRow facilite la lecture de données par colonnes, dans une ligne.
type IndexedRow struct {
	idx ColIndex
	row []string
}

// GetVal retourne la valeur associée à la colonne donnée, sur la ligne en cours.
// Dans le cas où la colonne n'existe pas, une erreur fatale est déclenchée.
func (indexedRow IndexedRow) GetVal(colName string) string {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		log.Fatal("Column not found in ColMapping: " + colName)
	}
	return indexedRow.row[index]
}

// GetOptionalVal retourne la valeur associée à la colonne donnée, sur la ligne en cours.
// Dans le cas où la colonne n'existe pas, le booléen sera faux.
func (indexedRow IndexedRow) GetOptionalVal(colName string) (string, bool) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return "", false
	}
	return indexedRow.row[index], true
}

// GetFloat64 retourne la valeur décimale associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetFloat64(colName string) (*float64, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return nil, fmt.Errorf("GetFloat64 failed to find column: %v", colName)
	}
	return misc.ParsePFloat(indexedRow.row[index])
}

// GetCommaFloat64 retourne la valeur décimale avec virgule associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetCommaFloat64(colName string) (*float64, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return nil, fmt.Errorf("GetCommaFloat64 failed to find column: %v", colName)
	}
	normalizedDecimalVal := strings.Replace(indexedRow.row[index], ",", ".", -1)
	return misc.ParsePFloat(normalizedDecimalVal)
}

// GetInt retourne la valeur entière associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetInt(colName string) (*int, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return nil, fmt.Errorf("GetInt failed to find column: %v", colName)
	}
	return misc.ParsePInt(indexedRow.row[index])
}

// GetIntFromFloat retourne la valeur entière associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetIntFromFloat(colName string) (*int, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return nil, fmt.Errorf("GetIntFromFloat failed to find column: %v", colName)
	}
	return misc.ParsePIntFromFloat(indexedRow.row[index])
}

// GetBool retourne la valeur booléenne associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetBool(colName string) (bool, error) {
	index, ok := indexedRow.idx.Get(colName)
	if !ok {
		return false, fmt.Errorf("GetBool failed to find column: %v", colName)
	}
	return strconv.ParseBool(indexedRow.row[index])
}

// LowercaseFields normalise les noms de colonnes en minuscules.
func LowercaseFields(headerFields []string) []string {
	var normalizedHeaderFields = make([]string, len(headerFields))
	for i, name := range headerFields {
		normalizedHeaderFields[i] = strings.ToLower(name)
	}
	return normalizedHeaderFields
}

// ExtractInputHeaders extrait les en-têtes des fichiers d'entrée via le tag "input"
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
