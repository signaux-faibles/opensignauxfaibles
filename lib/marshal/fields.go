package marshal

// GetFieldBindings indexe la position de chaque colonne par son nom,
// à partir de la liste ordonnée des noms de colonne, telle que lue en en-tête.
func GetFieldBindings(orderedFields []string) map[string]int {
	var f = map[string]int{}
	for i, k := range orderedFields {
		f[k] = i
	}
	return f
}
