package sfregexp

import "regexp"

// RegexpDict fournit des expressions régulières communes.
var RegexpDict = map[string]*regexp.Regexp{
	"siret":    regexp.MustCompile("^[0-9]{14}$"),
	"siren":    regexp.MustCompile("^[0-9]{9}$"),
	"notDigit": regexp.MustCompile("[^0-9]"),
	"nil":      regexp.MustCompile(".*"),
}
