package engine

// SirenFilter décrit le périmètre d'import.
// Les clés sont des numéros SIREN.
// Seules les clés présentes et dont la valeur est `true` sont dans le périmètre.
type SirenFilter interface {
	ShouldSkip(string) bool

	// Retourne tous les sirets du périmètre
	All() map[string]struct{}
}

// func (f SirenFilter) Add(siren string) error {
// 	if !sfregexp.RegexpDict["siren"].MatchString(siren) {
// 		return fmt.Errorf("format SIREN invalide: %s", siren)
// 	}
// 	f[siren] = true
// 	return nil
// }

type FilterWriter interface {
	Write(SirenFilter) error
}

// FilterReader retrieves a SirenFilter for a given batch.
// Implementations may read from files, databases, or other sources.
type FilterReader interface {
	Read() (SirenFilter, error)
}

var NoFilter SirenFilter = noFilter{}

type noFilter struct{}

func (f noFilter) ShouldSkip(string) bool { return false }

func (f noFilter) All() map[string]struct{} { return nil }
