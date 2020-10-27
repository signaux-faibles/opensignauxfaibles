package sfregexp

import (
	"regexp"
)

// RegexpDict fournit des expressions régulières communes.
var RegexpDict = map[string]*regexp.Regexp{
	"siret":    regexp.MustCompile("^[0-9]{14}$"),
	"siren":    regexp.MustCompile("^[0-9]{9}$"),
	"notDigit": regexp.MustCompile("[^0-9]"),
}

// ValidSiret retourne `true` si le numéro SIRET est valide.
func ValidSiret(siret string) bool {
	return RegexpDict["siret"].MatchString(siret)
}

// ValidSiren retourne `true` si le numéro SIREN est valide.
func ValidSiren(siren string) bool {
	return RegexpDict["siren"].MatchString(siren)
}
