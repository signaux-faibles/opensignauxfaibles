package marshal

import (
	"errors"
	"opensignauxfaibles/dbmongo/lib/engine"
	"opensignauxfaibles/dbmongo/lib/misc"
	"regexp"
	"strconv"
	"time"
)

var ParserDict = map[string]fullParser{
	"time":          ParsePTime,
	"urssafDate":    ParseUrssafPDate,
	"urssafPeriode": ParseUrssafPPeriod,
	"string":        ParseString,
	"float":         ParsePFloat,
	"int":           ParsePInt,
	"bool":          ParseBool,
}

func PossibleParsers() []string {
	var possibleParsers = make([]string, len(ParserDict))
	var i = 0
	for k := range ParserDict {
		possibleParsers[i] = k
		i++
	}
	return possibleParsers
}

// Interface to parse strings that are not empty and seem valid
type parser func(string, ...string) (interface{}, error)
type fullParser func(string, string, string, *regexp.Regexp, ...string) (interface{}, error)

// Parse field, with specific behaviors if the field is missing or not valid.
// A field is considered missing if it is equal to ""
// Returns as an interface{} the parsed field, or nil if no field was parsed,
// as well as possible errors (of type engine.CriticError)
// If validityRegex==nil, then no additional validity check is made.
func parseGeneric(
	parser parser,
	field string,
	ifEmpty string,
	ifInvalid string,
	validityRegex *regexp.Regexp,
	options ...string,
) (interface{}, error) {

	err := CheckParameter(ifEmpty)
	if err != nil {
		return nil, engine.NewCriticError(err, "fatal")
	}

	// if empty field
	if field == "" {
		switch ifEmpty {
		case "ignore":
			return nil, nil
		default:
			return nil, engine.NewCriticError(errors.New("empty field"), ifEmpty)
		}
	}

	err = CheckParameter(ifInvalid)
	if err != nil {
		return nil, engine.NewCriticError(err, "fatal")
	}
	if validityRegex == nil {
		validityRegex = regexp.MustCompile(".*")
	}
	isValid := validityRegex.MatchString(field)

	i, err := parser(field, options...)

	// if invalid field
	if !isValid || err != nil {
		switch ifInvalid {
		case "ignore":
			return nil, nil
		default:
			return nil, engine.NewCriticError(errors.New("Invalid field: "+field), ifInvalid)
		}
	}

	return i, nil
}

// CheckParameter checks parameters and if parameter is valid returns prefix
func CheckParameter(parameter string) error {

	validBehavior := map[string]struct{}{
		"ignore": {},
		"filter": {},
		"error":  {},
		"fatal":  {},
	}

	_, ok := validBehavior[parameter]
	if !ok {
		errString := "Parameter must be equal to one of"
		for k := range validBehavior {
			errString = errString + " " + k
		}
		return errors.New(errString)
	}
	return nil
}

//
// Parse floats (pointers)
//

// pFloatParser parses a valid float as a pointer. No option.
func pFloatParser(s string, options ...string) (interface{}, error) {
	f, err := strconv.ParseFloat(s, 64)
	var i interface{}
	i = &f
	return i, err
}

// ParsePFloat parses a float as a pointer and deals with missing and invalid
// values
func ParsePFloat(s string, ifEmpty string, ifInvalid string, validityRegex *regexp.Regexp, notUsed ...string) (interface{}, error) {
	i, err := parseGeneric(pFloatParser, s, ifEmpty, ifInvalid, validityRegex)
	f, _ := i.(*float64)
	return f, err
}

//
// Parse ints
//
// pIntParser parses a valid int as a pointer. No option.
func pIntParser(s string, options ...string) (interface{}, error) {
	in, err := strconv.Atoi(s)
	var i interface{}
	i = &in
	return i, err
}

// ParsePInt parses a int as a pointer and deals with missing and invalid
// values
func ParsePInt(s string, ifEmpty string, ifInvalid string, validityRegex *regexp.Regexp, notUsed ...string) (interface{}, error) {
	i, err := parseGeneric(pIntParser, s, ifEmpty, ifInvalid, validityRegex)
	return i, err
}

//
// Parse strings
//
// stringParser returns the string unchanged. No option.
func stringParser(s string, options ...string) (interface{}, error) {
	return s, nil
}

// ParseString parses a string and deals with missing and invalid values
func ParseString(s string, ifEmpty string, ifInvalid string, validityRegex *regexp.Regexp, notUsed ...string) (interface{}, error) {
	i, err := parseGeneric(stringParser, s, ifEmpty, ifInvalid, validityRegex)
	return i, err
}

//
// Parse regular times
//

// pTimeParser parses dates and times as pointers. Layout of time.Parse is
// given as option.
func pTimeParser(s string, options ...string) (interface{}, error) {
	if len(options) != 1 {
		panic("PTime parsing requires one option: the time.Parse layout")
	}
	t, err := time.Parse(options[0], s)
	var i interface{}
	i = &t
	return i, err
}

// ParsePTime parses a time as a pointer
func ParsePTime(s string, ifEmpty string, ifInvalid string, validityRegex *regexp.Regexp, timeFormat ...string) (interface{}, error) {
	i, err := parseGeneric(pTimeParser, s, ifEmpty, ifInvalid, nil, timeFormat[0])
	return i, err
}

//////////////////////////////////////////
// Parse times in urssaf period: 200612 //
//////////////////////////////////////////

// urssafPPeriodParser parses an urssaf period time as a pointer  and returns
// the period start. No option.
func urssafPPeriodParser(s string, options ...string) (interface{}, error) {
	p, err := urssafToPeriod(s)
	var i interface{}
	i = &(p.Start)
	return i, err
}

// ParseUrssafPPeriod parses a urssaf period time as a pointer.
func ParseUrssafPPeriod(s string, ifEmpty string, ifInvalid string, validityRegex *regexp.Regexp, notUsed ...string) (interface{}, error) {
	i, err := parseGeneric(urssafPPeriodParser, s, ifEmpty, ifInvalid, validityRegex)
	return i, err
}

/////////////////////////////////////////
// Parse times in urssaf time: 1060301 //
/////////////////////////////////////////

// urssafPDateParser parses a urssaf date. No option.
func urssafPDateParser(s string, options ...string) (interface{}, error) {
	p, err := urssafToDate(s)
	var i interface{}
	i = &p
	return i, err
}

// ParseUrssafPDate parses a urssaf date time as a pointer, and deals with
// empty and invalid inputs.
func ParseUrssafPDate(s string, ifEmpty string, ifInvalid string, validityRegex *regexp.Regexp, notUsed ...string) (interface{}, error) {
	i, err := parseGeneric(urssafPDateParser, s, ifEmpty, ifInvalid, validityRegex)
	return i, err
}

////////////////
// Parse bool //
////////////////

func boolParser(s string, options ...string) (interface{}, error) {
	b, err := strconv.ParseBool(s)
	return b, err
}

// ParseBool parses a bool
func ParseBool(s string, ifEmpty string, ifInvalid string, validityRegex *regexp.Regexp, notUsed ...string) (interface{}, error) {
	i, err := parseGeneric(boolParser, s, ifEmpty, ifInvalid, validityRegex)
	return i, err
}

// GetParserType gets parser type from parser name
//func GetParserType(parserName string) (string, error) {
//	//TODO do the same with reflect !
//	switch parserName {
//	case "ParsePTime", "ParseUrssafPPeriod", "ParseUrssafPDate":
//		return "*time.Time", nil
//	case "ParseString":
//		return "string", nil
//	case "ParsePFloat":
//		return "*float64", nil
//	case "ParsePInt":
//		return "*int", nil
//	case "ParseBool":
//		return "bool", nil
//	default:
//		return "", errors.New("Parser not found")
//	}
//}

// UrssafToDate convertit le format de date urssaf en type Date.
// Les dates urssaf sont au format YYYMMJJ tels que YYY = YYYY - 1900 (e.g: 118 signifie
// 2018)
func urssafToDate(urssaf string) (time.Time, error) {

	intUrsaff, err := strconv.Atoi(urssaf)
	if err != nil {
		return time.Time{}, engine.NewCriticError(errors.New("Valeur non autorisée pour une conversion en date: "+urssaf), "fatal")
	}
	strDate := strconv.Itoa(intUrsaff + 19000000)
	date, err := time.Parse("20060102", strDate)
	if err != nil {
		return time.Time{}, engine.NewCriticError(errors.New("Valeur non autorisée pour une conversion en date: "+urssaf), "fatal")
	}

	return date, nil
}

// UrssafToPeriod convertit le format de période urssaf en type misc.Periode. On trouve ces
// périodes formatées en 4 ou 6 caractère (YYQM ou YYYYQM).
// si YY < 50 alors YYYY = 20YY sinon YYYY = 19YY.
// si QM == 62 alors période annuelle sur YYYY.
// si M == 0 alors période trimestrielle sur le trimestre Q de YYYY.
// si 0 < M < 4 alors mois M du trimestre Q.
func urssafToPeriod(urssaf string) (misc.Periode, error) {
	period := misc.Periode{}

	if len(urssaf) == 4 {
		if urssaf[0:2] < "50" {
			urssaf = "20" + urssaf
		} else {
			urssaf = "19" + urssaf
		}
	}

	if len(urssaf) != 6 {
		return period, errors.New("Valeur non autorisée")
	}

	year, err := strconv.Atoi(urssaf[0:4])
	if err != nil {
		return misc.Periode{}, errors.New("Valeur non autorisée")
	}

	if urssaf[4:6] == "62" {
		period.Start = time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		period.End = time.Date(year+1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	} else {
		quarter, err := strconv.Atoi(urssaf[4:5])
		if err != nil {
			return period, err
		}
		monthOfQuarter, err := strconv.Atoi(urssaf[5:6])
		if err != nil {
			return period, err
		}
		if monthOfQuarter == 0 {
			period.Start = time.Date(year, time.Month((quarter-1)*3+1), 1, 0, 0, 0, 0, time.UTC)
			period.End = time.Date(year, time.Month((quarter-1)*3+4), 1, 0, 0, 0, 0, time.UTC)
		} else {
			period.Start = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter), 1, 0, 0, 0, 0, time.UTC)
			period.End = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter+1), 1, 0, 0, 0, 0, time.UTC)
		}
	}
	return period, nil
}
