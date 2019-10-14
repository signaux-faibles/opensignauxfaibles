package marshal

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
)

// testHelperError tests if an error is as expected
func testHelperError(
	err error,
	ind int,
	expectedError bool,
	expectedErrorType string,
	t *testing.T,
) {
	if err != nil && !expectedError {
		t.Errorf("Case %d failed: unexpected error : %s", ind, err.Error())
	}
	if err == nil && expectedError {
		t.Errorf("Case %d failed: shoud have thrown an error", ind)
	}
	if err != nil {
		switch c := err.(type) {
		case *engine.CriticError:
			if c.Criticity() != expectedErrorType {
				t.Error("Wrong criticity of error")
			}
		default:
			t.Error("wrong type of error")
		}
	}
}

func TestParseGeneric(t *testing.T) {
	validityRegex := regexp.MustCompile("^[A-Z]{2}$")
	stdParser := func(s string, options ...string) (interface{}, error) {
		return s, nil
	}
	errParser := func(s string, options ...string) (interface{}, error) {
		return "", errors.New("error")
	}

	cases := []struct {
		input             string
		parser            parser
		ifEmpty           string
		ifInvalid         string
		expectedValue     interface{}
		expectedError     bool
		expectedErrorType string
	}{
		// Valid
		{"ML", stdParser, "ignore", "ignore", "ML", false, ""},
		{"ML", stdParser, "ignore", "filter", "ML", false, ""},
		{"ML", stdParser, "ignore", "error", "ML", false, ""},
		{"ML", stdParser, "ignore", "fatal", "ML", false, ""},
		{"ML", stdParser, "filter", "ignore", "ML", false, ""},
		{"ML", stdParser, "filter", "filter", "ML", false, ""},
		{"ML", stdParser, "filter", "error", "ML", false, ""},
		{"ML", stdParser, "filter", "fatal", "ML", false, ""},
		{"ML", stdParser, "error", "ignore", "ML", false, ""},
		{"ML", stdParser, "error", "filter", "ML", false, ""},
		{"ML", stdParser, "error", "error", "ML", false, ""},
		{"ML", stdParser, "error", "fatal", "ML", false, ""},
		{"ML", stdParser, "fatal", "ignore", "ML", false, ""},
		{"ML", stdParser, "fatal", "filter", "ML", false, ""},
		{"ML", stdParser, "fatal", "error", "ML", false, ""},
		{"ML", stdParser, "fatal", "fatal", "ML", false, ""},
		// empty
		{"", stdParser, "ignore", "ignore", nil, false, ""},
		{"", stdParser, "ignore", "filter", nil, false, ""},
		{"", stdParser, "ignore", "error", nil, false, ""},
		{"", stdParser, "ignore", "fatal", nil, false, ""},
		{"", stdParser, "filter", "ignore", nil, true, "filter"},
		{"", stdParser, "filter", "filter", nil, true, "filter"},
		{"", stdParser, "filter", "error", nil, true, "filter"},
		{"", stdParser, "filter", "fatal", nil, true, "filter"},
		{"", stdParser, "error", "ignore", nil, true, "error"},
		{"", stdParser, "error", "filter", nil, true, "error"},
		{"", stdParser, "error", "error", nil, true, "error"},
		{"", stdParser, "error", "fatal", nil, true, "error"},
		{"", stdParser, "fatal", "ignore", nil, true, "fatal"},
		{"", stdParser, "fatal", "filter", nil, true, "fatal"},
		{"", stdParser, "fatal", "error", nil, true, "fatal"},
		{"", stdParser, "fatal", "fatal", nil, true, "fatal"},
		// invalid regexp
		{"12", stdParser, "ignore", "ignore", nil, false, ""},
		{"12", stdParser, "filter", "ignore", nil, false, ""},
		{"12", stdParser, "error", "ignore", nil, false, ""},
		{"12", stdParser, "fatal", "ignore", nil, false, ""},
		{"12", stdParser, "ignore", "filter", nil, true, "filter"},
		{"12", stdParser, "filter", "filter", nil, true, "filter"},
		{"12", stdParser, "error", "filter", nil, true, "filter"},
		{"12", stdParser, "fatal", "filter", nil, true, "filter"},
		{"12", stdParser, "ignore", "error", nil, true, "error"},
		{"12", stdParser, "filter", "error", nil, true, "error"},
		{"12", stdParser, "error", "error", nil, true, "error"},
		{"12", stdParser, "fatal", "error", nil, true, "error"},
		{"12", stdParser, "ignore", "fatal", nil, true, "fatal"},
		{"12", stdParser, "filter", "fatal", nil, true, "fatal"},
		{"12", stdParser, "error", "fatal", nil, true, "fatal"},
		{"12", stdParser, "fatal", "fatal", nil, true, "fatal"},
		// invalid parsing
		{"AB", errParser, "ignore", "ignore", nil, false, ""},
		{"AB", errParser, "filter", "ignore", nil, false, ""},
		{"AB", errParser, "error", "ignore", nil, false, ""},
		{"AB", errParser, "fatal", "ignore", nil, false, ""},
		{"AB", errParser, "ignore", "filter", nil, true, "filter"},
		{"AB", errParser, "filter", "filter", nil, true, "filter"},
		{"AB", errParser, "error", "filter", nil, true, "filter"},
		{"AB", errParser, "fatal", "filter", nil, true, "filter"},
		{"AB", errParser, "ignore", "error", nil, true, "error"},
		{"AB", errParser, "filter", "error", nil, true, "error"},
		{"AB", errParser, "error", "error", nil, true, "error"},
		{"AB", errParser, "fatal", "error", nil, true, "error"},
		{"AB", errParser, "ignore", "fatal", nil, true, "fatal"},
		{"AB", errParser, "filter", "fatal", nil, true, "fatal"},
		{"AB", errParser, "error", "fatal", nil, true, "fatal"},
		{"AB", errParser, "fatal", "fatal", nil, true, "fatal"},
	}

	for ind, tc := range cases {
		i, err := parseGeneric(
			tc.parser,
			tc.input,
			tc.ifEmpty,
			tc.ifInvalid,
			validityRegex,
		)
		if (i == nil && tc.expectedValue != nil) ||
			(i != nil && tc.expectedValue == nil) {
			t.Errorf("Case %d failed: actual or expected nil", ind)
		} else if i != nil && i.(string) != tc.expectedValue {
			t.Errorf("Case %d failed: %s is not equal to %s", ind, i.(string), tc.expectedValue)
		}

		testHelperError(err, ind, tc.expectedError, tc.expectedErrorType, t)

	}
}

func TestParsePFloat(t *testing.T) {
	res, err := ParsePFloat("54.234", "fatal", "fatal", nil)
	if err != nil {
		t.Errorf("Unexpected parsing error %v", err)
	}
	typedRes := res.(*float64)
	if *typedRes != 54.234 {
		t.Errorf("Unexpected parsing typedResults %f", *typedRes)
	}
	res, _ = ParsePFloat("foo", "fatal", "fatal", nil)
	typedRes = res.(*float64)
	if typedRes != nil {
		t.Errorf("Invalid value does not return nil value")
	}
}

func TestParsePInt(t *testing.T) {
	res, err := ParsePInt("54", "fatal", "fatal", nil)
	if err != nil {
		t.Errorf("Unexpected parsing error %v", err)
	}

	typedRes := res.(*int)
	if *typedRes != 54 {
		t.Errorf("Unexpected parsing typedResults %d", *typedRes)
	}
	res, _ = ParsePInt("foo", "fatal", "fatal", nil)
	if res != nil {
		t.Errorf("Invalid value does not return nil value")
	}
}

func TestParseString(t *testing.T) {
	res, err := ParseString("abc", "fatal", "fatal", nil)
	if err != nil {
		t.Errorf("Unexpected parsing error %v", err)
	}
	typedRes := res.(string)
	if typedRes != "abc" {
		t.Errorf("Unexpected parsing typedResults %s", typedRes)
	}
}

func TestParsePTime(t *testing.T) {
	res, err := ParsePTime("17-04-1990", "fatal", "fatal", nil, "02-01-2006")
	if err != nil {
		t.Errorf("Unexpected parsing error %v", err)
	}
	typedRes := res.(*time.Time)
	aux, _ := time.Parse("02-01-2006", "17-04-1990")
	if *typedRes != aux {
		t.Errorf("Unexpected parsing typedResults %s", typedRes)
	}
	res, _ = ParsePTime("foo", "fatal", "fatal", nil, "")
	if res != nil {
		t.Errorf("Invalid value does not return nil value")
	}
}

func TestParseBool(t *testing.T) {
	res, err := ParseBool("true", "fatal", "fatal", nil)
	if err != nil {
		t.Errorf("Unexpected parsing error %v", err)
	}
	typedRes := res.(bool)
	if typedRes != true {
		t.Errorf("Unexpected parsing typedResults")
	}
	res, _ = ParseBool("foo", "fatal", "fatal", nil)
	if res != nil {
		t.Errorf("Invalid value does not return nil value")
	}
}
