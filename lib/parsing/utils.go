package parsing

import "strconv"

// ParsePInt parse un entier et retourne un pointeur
func ParsePInt(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	i, err := strconv.Atoi(s)
	return &i, err
}

// ParsePIntFromFloat parse un float, le transforme en int et retourne un pointeur sur l'int
func ParsePIntFromFloat(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	var i = int(f)
	return &i, err
}

// ParsePFloat parse un flottant et retourne un pointeur
func ParsePFloat(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}
	i, err := strconv.ParseFloat(s, 64)
	return &i, err
}
