package marshal

import (
	"encoding/csv"
	"errors"
	"reflect"
	"strings"
)

// IndexColumnsFromCsvHeader extrait les noms de colonnes depuis l'en-tête d'un
// flux CSV puis les indexe à l'aide de ValidateAndIndexColumnsFromColTags().
func IndexColumnsFromCsvHeader(reader *csv.Reader, destObject interface{}) (ColMapping, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return ValidateAndIndexColumnsFromColTags(header, destObject)
}

// ValidateAndIndexColumnsFromColTags valide puis indexe les colonnes trouvées
// en en-tête d'un fichier csv, à partir des noms de colonnes spécifiés dans le
// tag "col" annotant les propriétés du type de destination du parseur.
func ValidateAndIndexColumnsFromColTags(headerRow []string, destObject interface{}) (ColMapping, error) {
	requiredFields := extractColTags(destObject)
	idx := indexFields(headerRow)
	_, err := idx.HasFields(requiredFields)
	return idx, err
}

// indexFields indexe la position de chaque colonne par son nom,
// à partir de la liste ordonnée des noms de colonne, telle que lue en en-tête.
func indexFields(headerFields []string) ColMapping {
	var colMapping = ColMapping{}
	for idx, name := range headerFields {
		colMapping[name] = idx
	}
	return colMapping
}

// ColMapping fournit l'indice de chaque colonne.
type ColMapping map[string]int

// HasFields vérifie la présence d'un ensemble de colonnes.
func (colMapping ColMapping) HasFields(requiredFields []string) (bool, error) {
	for _, name := range requiredFields {
		if _, found := colMapping[name]; !found {
			return false, errors.New("Colonne " + name + " non trouvée. Abandon.")
		}
	}
	return true, nil
}

// LowercaseFields normalise les noms de colonnes en minuscules.
func LowercaseFields(headerFields []string) []string {
	var normalizedHeaderFields = make([]string, len(headerFields))
	for i, name := range headerFields {
		normalizedHeaderFields[i] = strings.ToLower(name)
	}
	return normalizedHeaderFields
}

// extractColTags extraie les noms de colonnes depuis les valeurs du tag "col"
// de chaque propriété de l'objet fourni.
// Il est possible d'associer plusieurs colonnes en séparant par des virgules.
func extractColTags(object interface{}) (expectedFields []string) {
	structure := reflect.TypeOf(object)
	for i := 0; i < structure.NumField(); i++ {
		tag := structure.Field(i).Tag.Get("col")
		if tag != "" {
			expectedFields = append(expectedFields, strings.Split(tag, ",")...)
		}
	}
	return expectedFields
}
