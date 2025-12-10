package sfregexp

import (
	"regexp"
)

// RegexpDict fournit des expressions régulières communes.
var RegexpDict = map[string]*regexp.Regexp{
	"siret":      regexp.MustCompile("^[0-9]{14}$"),
	"siren":      regexp.MustCompile("^[0-9]{9}$"),
	"notDigit":   regexp.MustCompile("[^0-9]"),
	"fiveDigits": regexp.MustCompile(`^\d{5}$`),
}

// ValidSiret retourne `true` si le numéro SIRET est valide.
func ValidSiret(siret string) bool {
	return RegexpDict["siret"].MatchString(siret)
}

// ValidSiren retourne `true` si le numéro SIREN est valide.
func ValidSiren(siren string) bool {
	return RegexpDict["siren"].MatchString(siren)
}

// ValidCodePostal retourne `true` si le code postal est composé de 5 chiffres.
func ValidCodePostal(codePostal string) bool {
	return RegexpDict["fiveDigits"].MatchString(codePostal)
}

// ValidCodeCommune retourne `true` si le code postal est composé de 5 chiffres.
func ValidCodeCommune(codeCommune string) bool {
	return RegexpDict["fiveDigits"].MatchString(codeCommune)
}
