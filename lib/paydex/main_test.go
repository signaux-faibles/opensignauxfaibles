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

func TestPaydex(t *testing.T) {
	t.Run("can parse a line", func(t *testing.T) {
		row := []string{"000000001", "2", "2 jours", "15/12/2018"}
		expected := Paydex{
			Siren:      "000000001",
			DateValeur: time.Date(2018, 12, 01, 00, 00, 00, 0, time.UTC),
			NbJours:    2,
		}
		actual := parsePaydexLine(row)
		assert.Equal(t, expected, actual)
	})

	t.Run("generate the right tuples and events from test file", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedPaydex.json")
		var testData = filepath.Join("testData", "paydexTestData.csv")
		marshal.TestParserOutput(t, ParserPaydex, marshal.NewCache(), testData, golden, *update)
	})
}
