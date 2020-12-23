package marshal

// GetFieldBindings indexe la position de chaque colonne par nom.
func GetFieldBindings(fields []string) map[string]int {
	var f = map[string]int{}
	for i, k := range fields {
		f[k] = i
	}
	return f
}
