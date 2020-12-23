package marshal

import (
	"errors"
	"strings"

	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
)

// IndexFields indexe la position de chaque colonne par son nom,
// à partir d'un en-tête ordonné et d'une liste de colonnes attendues.
func IndexFields(headerFields []string, expectedFields []string) ColMapping {
	var colMapping = ColMapping{}
	for _, name := range expectedFields {
		idx := misc.SliceIndex(len(headerFields), func(i int) bool { return headerFields[i] == name })
		if idx != -1 {
			colMapping[name] = idx
		}
	}
	return colMapping
}

// GetFieldBindings indexe la position de chaque colonne par son nom,
// à partir de la liste ordonnée des noms de colonne, telle que lue en en-tête.
func GetFieldBindings(orderedFields []string) ColMapping {
	var colMapping = ColMapping{}
	for idx, name := range orderedFields {
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
