package parsing

import (
	"opensignauxfaibles/lib/engine"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestTuple struct {
	a string `input:"a"`
	b string `input:"b"`
	c string
}

type TestRowParser struct{}

func (TestRowParser) ParseRow(row []string, result *engine.ParsedLineResult, idx ColIndex) {}

func TestInitWithMissingCSVHeader(t *testing.T) {
	testCases := []struct {
		name          string
		csvContent    string
		comma         rune
		caseSensitive bool
		expectError   bool
	}{
		{
			name:          "Headers are complete",
			csvContent:    "a,b,c\n1,2,3",
			comma:         ',',
			caseSensitive: true,
			expectError:   false,
		},
		{
			name:          "b is missing",
			csvContent:    "a\n1",
			comma:         ',',
			caseSensitive: true,
			expectError:   true,
		},
		{
			name:          "Optional c is missing",
			csvContent:    "a,b\n1,2",
			comma:         ',',
			caseSensitive: true,
			expectError:   false,
		},
		{
			name:          "Changing the comma works",
			csvContent:    "a;b\n1;2",
			comma:         ';',
			caseSensitive: true,
			expectError:   false,
		},
		{
			name:          "Changing the comma works (2)",
			csvContent:    "a;b\n1;2",
			comma:         ',',
			caseSensitive: true,
			expectError:   true,
		},
		{
			name:          "wrong casing with case sensitivity",
			csvContent:    "A,b,c\n1,2,3",
			comma:         ',',
			caseSensitive: true,
			expectError:   true,
		},
		{
			name:          "wrong casing with case INsensitivity",
			csvContent:    "A,b,c\n1,2,3",
			comma:         ',',
			caseSensitive: false,
			expectError:   false,
		},
		{
			name:          "trims newlines before headers on second line",
			csvContent:    "\n\n\n\na,b,c\n1,2,3",
			comma:         ',',
			caseSensitive: true,
			expectError:   false,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			content := strings.NewReader(tc.csvContent)
			parser := &CsvParserInst{
				Reader:        content,
				RowParser:     TestRowParser{},
				Comma:         tc.comma,
				CaseSensitive: tc.caseSensitive,
				LazyQuotes:    false,
				DestTuple:     TestTuple{},
			}

			err := parser.Init(engine.NoFilter, nil)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMalformedCSV(t *testing.T) {
	malformedCSV := "a,b\n3\n1,2"
	content := strings.NewReader(malformedCSV)
	parser := &CsvParserInst{
		Reader:        content,
		RowParser:     TestRowParser{},
		Comma:         ',',
		CaseSensitive: true,
		LazyQuotes:    false,
		DestTuple:     TestTuple{},
	}

	err := parser.Init(engine.NoFilter, nil)
	require.NoError(t, err)

	res := &engine.ParsedLineResult{}
	err = parser.ReadNext(res)
	require.NoError(t, err, "should not interrupt parsing because of a malformed row")
	assert.Len(t, res.Errors, 1, "should issue a result error when encounting a malformed row")

	err = parser.ReadNext(res)
	require.NoError(t, err, "should succesfully continue to parse after a malformed row")
	assert.Len(t, res.Errors, 1, "should succesfully continue to parse after a malformed row")
}
