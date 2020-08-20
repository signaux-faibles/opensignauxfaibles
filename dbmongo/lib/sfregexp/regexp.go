package sfregexp

import "regexp"

var RegexpDict = map[string]*regexp.Regexp{
	"siret": regexp.MustCompile("^[0-9]{14}$"),
	"siren": regexp.MustCompile("^[0-9]{9}$"),
	"nil":   regexp.MustCompile(".*"),
}
