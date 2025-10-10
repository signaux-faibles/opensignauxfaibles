package engine

// SirenFilter décrit le périmètre d'import.
// Les clés sont des numéros SIREN.
// Seules les clés présentes et dont la valeur est `true` sont dans le périmètre.
type SirenFilter interface {
	ShouldSkip(string) bool
}

// func (f SirenFilter) Add(siren string) error {
// 	if !sfregexp.RegexpDict["siren"].MatchString(siren) {
// 		return fmt.Errorf("format SIREN invalide: %s", siren)
// 	}
// 	f[siren] = true
// 	return nil
// }

// FilterReader defines the interface for reading SIREN filters from various sources.
type FilterReader interface {
	Read() (SirenFilter, error)
	SuccessStr() string
}

var NoFilter SirenFilter = noFilter{}

type noFilter struct{}

func (f noFilter) ShouldSkip(string) bool { return false }
