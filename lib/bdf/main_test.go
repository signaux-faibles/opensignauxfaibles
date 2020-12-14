package bdf

import (
	"encoding/json"
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var testData = filepath.Join("testData", "bdfTestData.csv") // ce fichier définit 3 entreprises: 000111222, 000111223 et 000111224

func TestBdfOutput(t *testing.T) {
	var golden = filepath.Join("testData", "expectedBdfOutput.json")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)
}

func TestBdfParser(t *testing.T) {
	t.Run("Should return only the companies listed in the filter file", func(t *testing.T) {
		var cache = marshal.NewCache()
		cache.Set("filter", marshal.SirenFilter{"000111222": true, "000111224": true})
		var output = marshal.RunParser(Parser, cache, testData)
		assert.Equal(t, 2, len(output.Tuples))
		lastEvent := output.Events[len(output.Events)-1]
		lastReportJSON, err := json.MarshalIndent(lastEvent, "", "  ")
		assert.NoError(t, err)
		assert.Contains(t, string(lastReportJSON), "3 lignes traitées, 0 erreurs fatales, 0 lignes rejetées, 1 lignes filtrées, 2 lignes valides")
	})
}
