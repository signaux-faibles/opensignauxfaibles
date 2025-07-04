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
	return ValidateAndIndexColumnsFromColTags(header, destObject)
}

// ValidateAndIndexColumnsFromColTags valide puis indexe les colonnes trouvées
// en en-tête d'un fichier csv, à partir des noms de colonnes spécifiés dans le
// tag "col" ou "cols" annotant les propriétés du type de destination du parseur.
func ValidateAndIndexColumnsFromColTags(headerRow []string, destObject interface{}) (ColMapping, error) {
	requiredFields := ExtractColTags(destObject)
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

// ExtractColTags extraie les noms de colonnes depuis les valeurs du tag "col"
// (ou "cols", séparés par des virgules) de chaque propriété de l'objet fourni.
// Il est possible d'associer plusieurs colonnes en séparant par des virgules.
func ExtractColTags(object interface{}) (expectedFields []string) {
	structure := reflect.TypeOf(object)
	for i := 0; i < structure.NumField(); i++ {
		tagWithSingleVal := structure.Field(i).Tag.Get("col")
		if tagWithSingleVal != "" {
			expectedFields = append(expectedFields, tagWithSingleVal)
		}
		tagWithMultiVal := structure.Field(i).Tag.Get("cols")
		if tagWithMultiVal != "" {
			expectedFields = append(expectedFields, strings.Split(tagWithMultiVal, ",")...)
		}
	}
	return expectedFields
}
