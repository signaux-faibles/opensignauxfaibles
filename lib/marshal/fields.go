package marshal

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"opensignauxfaibles/lib/misc"
)

// IndexColumnsFromCsvHeader extrait les noms de colonnes depuis l'en-tête d'un
// flux CSV puis les indexe à l'aide de ValidateAndIndexColumnsFromColTags().
func IndexColumnsFromCsvHeader(reader *csv.Reader, destObject interface{}) (ColMapping, error) {
	header, err := reader.Read()
	if err != nil {
		return ColMapping{}, err
	}
	return ValidateAndIndexColumnsFromInputTags(header, destObject)
}

// ValidateAndIndexColumnsFromInputTags valide puis indexe les colonnes trouvées
// en en-tête d'un fichier csv, à partir des noms de colonnes spécifiés dans le
// tag "input"  annotant les propriétés du type de destination du parseur.
func ValidateAndIndexColumnsFromInputTags(headerRow []string, destObject interface{}) (ColMapping, error) {
	requiredFields := ExtractInputHeaders(destObject)
	idx := indexFields(headerRow)
	_, err := idx.HasFields(requiredFields)
	return idx, err
}

// indexFields indexe la position de chaque colonne par son nom,
// à partir de la liste ordonnée des noms de colonne, telle que lue en en-tête.
func indexFields(headerFields []string) ColMapping {
	var colMapping = ColMapping{index: map[string]int{}}
	for idx, name := range headerFields {
		colMapping.index[name] = idx
	}
	return colMapping
}

// CreateColMapping créée un index de colonnes. Utile pour constituer des données de test.
func CreateColMapping(index map[string]int) ColMapping {
	return ColMapping{index}
}

// ColMapping fournit l'indice de chaque colonne.
type ColMapping struct {
	index map[string]int
}

// HasFields vérifie la présence d'un ensemble de colonnes.
func (colMapping ColMapping) HasFields(requiredFields []string) (bool, error) {
	for _, name := range requiredFields {
		if _, found := colMapping.index[name]; !found {
			return false, errors.New("Colonne " + name + " non trouvée. Abandon.")
		}
	}
	return true, nil
}

// IndexRow retourne une structure pour faciliter la lecture de données.
func (colMapping ColMapping) IndexRow(row []string) IndexedRow {
	return IndexedRow{colMapping, row}
}

// IndexedRow facilite la lecture de données par colonnes, dans une ligne.
type IndexedRow struct {
	colMaping ColMapping
	row       []string
}

// GetVal retourne la valeur associée à la colonne donnée, sur la ligne en cours.
// Dans le cas où la colonne n'existe pas, une erreur fatale est déclenchée.
func (indexedRow IndexedRow) GetVal(colName string) string {
	index, ok := indexedRow.colMaping.index[colName]
	if !ok {
		log.Fatal("Column not found in ColMapping: " + colName)
	}
	return indexedRow.row[index]
}

// GetOptionalVal retourne la valeur associée à la colonne donnée, sur la ligne en cours.
// Dans le cas où la colonne n'existe pas, le booléen sera faux.
func (indexedRow IndexedRow) GetOptionalVal(colName string) (string, bool) {
	index, ok := indexedRow.colMaping.index[colName]
	if !ok {
		return "", false
	}
	return indexedRow.row[index], true
}

// GetFloat64 retourne la valeur décimale associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetFloat64(colName string) (*float64, error) {
	index, ok := indexedRow.colMaping.index[colName]
	if !ok {
		return nil, fmt.Errorf("GetFloat64 failed to find column: %v", colName)
	}
	return misc.ParsePFloat(indexedRow.row[index])
}

// GetCommaFloat64 retourne la valeur décimale avec virgule associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetCommaFloat64(colName string) (*float64, error) {
	index, ok := indexedRow.colMaping.index[colName]
	if !ok {
		return nil, fmt.Errorf("GetCommaFloat64 failed to find column: %v", colName)
	}
	normalizedDecimalVal := strings.Replace(indexedRow.row[index], ",", ".", -1)
	return misc.ParsePFloat(normalizedDecimalVal)
}

// GetInt retourne la valeur entière associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetInt(colName string) (*int, error) {
	index, ok := indexedRow.colMaping.index[colName]
	if !ok {
		return nil, fmt.Errorf("GetInt failed to find column: %v", colName)
	}
	return misc.ParsePInt(indexedRow.row[index])
}

// GetIntFromFloat retourne la valeur entière associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetIntFromFloat(colName string) (*int, error) {
	index, ok := indexedRow.colMaping.index[colName]
	if !ok {
		return nil, fmt.Errorf("GetIntFromFloat failed to find column: %v", colName)
	}
	return misc.ParsePIntFromFloat(indexedRow.row[index])
}

// GetBool retourne la valeur booléenne associée à la colonne donnée, sur la ligne en cours.
// Un pointeur nil est retourné si la colonne n'existe pas ou la valeur est une chaine vide.
func (indexedRow IndexedRow) GetBool(colName string) (bool, error) {
	index, ok := indexedRow.colMaping.index[colName]
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
