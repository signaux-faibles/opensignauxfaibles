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

// FilterResolver encapsulates the complete filter resolution logic:
// checking requirements, updating state, and reading the filter.
// This abstraction simplifies handling the --no-filter flag and
// maintains clean separation between filter orchestration and
// individual read/write operations
type FilterResolver interface {
	// Resolve performs all filter operations (check, update, read) and
	// returns the final SirenFilter to use for the import
	Resolve(batchFiles BatchFiles) (SirenFilter, error)
}

var NoFilter SirenFilter = noFilter{}

type noFilter struct{}

func (f noFilter) ShouldSkip(string) bool { return false }

func (f noFilter) All() map[string]struct{} { return nil }
