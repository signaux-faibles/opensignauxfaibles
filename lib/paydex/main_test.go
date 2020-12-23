package paydex

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

// unit tests
func TestParsePaydexLine(t *testing.T) {
	t.Run("should parse a valid row", func(t *testing.T) {
		row := []string{"000000001", "2", "2 jours", "15/12/2018"}
		expected := Paydex{
			Siren:      "000000001",
			DateValeur: time.Date(2018, 12, 15, 00, 00, 00, 0, time.UTC),
			NbJours:    2,
		}
		colIndex := marshal.ColMapping{"SIREN": 0, "NB_JOURS": 1, "DATE_VALEUR": 3}
		actual, err := parsePaydexLine(colIndex, row)
		assert.Equal(t, expected, *actual)
		assert.Equal(t, nil, err)
	})

	t.Run("should report parse error on invalid date", func(t *testing.T) {
		row := []string{"000000001", "2", "2 jours", "12/15/2018"} // "15" is in the "month" slot
		colIndex := marshal.ColMapping{"SIREN": 0, "NB_JOURS": 1, "DATE_VALEUR": 3}
		actual, err := parsePaydexLine(colIndex, row)
		assert.EqualError(t, err, "invalid date: 12/15/2018")
		assert.Nil(t, actual)
	})
}

// integration tests
func TestPaydex(t *testing.T) {
	t.Run("should fail if one of the required columns is missing", func(t *testing.T) {
		output := marshal.RunParserInline(t, ParserPaydex, []string{"SIREN;NB_JOURS_LIB;DATE_VALEUR"}) // NB_JOURS is missing
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, 1, len(output.Events), "should return a parsing report")
		reportData, _ := output.Events[0].ParseReport()
		assert.Equal(t, true, reportData["isFatal"], "should report a fatal error")
	})

	t.Run("should adapt to different order of columns", func(t *testing.T) {
		csvRows := []string{
			"SIREN;DATE_VALEUR;NB_JOURS",
			"000000001;15/12/2018;2",
		}
		output := marshal.RunParserInline(t, ParserPaydex, csvRows)
		expected := Paydex{
			Siren:      "000000001",
			DateValeur: time.Date(2018, 12, 15, 00, 00, 00, 0, time.UTC),
			NbJours:    2,
		}
		assert.EqualValues(t, []marshal.Tuple{&expected}, output.Tuples)
	})

	t.Run("should generate the right tuples and events from test file", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedPaydex.json")
		var testData = filepath.Join("testData", "paydexTestData.csv")
		marshal.TestParserOutput(t, ParserPaydex, marshal.NewCache(), testData, golden, *update)
	})
}
