package sfregexp

import "regexp"

var RegexpDict = map[string]*regexp.Regexp{
	"siret": regexp.MustCompile("^[0-9]{14}$"),
	"siren": regexp.MustCompile("^[0-9]{9}$"),
	"nil":   regexp.MustCompile(".*"),
}

func PossibleRegexp() []string {
	possibleRegexp := make([]string, len(RegexpDict))
	i := 0
	for k := range RegexpDict {
		possibleRegexp[i] = k
		i++
	}
	return possibleRegexp
}
