package marshal

import "errors"

// GetFieldBindings indexe la position de chaque colonne par son nom,
// à partir de la liste ordonnée des noms de colonne, telle que lue en en-tête.
func GetFieldBindings(orderedFields []string) ColMapping {
	var f = ColMapping{}
	for i, k := range orderedFields {
		f[k] = i
	}
	return f
}

// ColMapping fournit l'indice de chaque colonne.
type ColMapping map[string]int

// HasFields vérifie la présence d'un ensemble de colonnes.
func (colMapping ColMapping) HasFields(requiredFields []string) (bool, error) {
	for _, field := range requiredFields {
		if _, found := colMapping[field]; !found {
			return false, errors.New("Colonne " + field + " non trouvée. Abandon.")
		}
	}
	return true, nil
}
