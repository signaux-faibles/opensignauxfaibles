package bdf

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var testData = filepath.Join("testData", "bdfTestData.csv")

func TestBdfOutput(t *testing.T) {
	var golden = filepath.Join("testData", "expectedBdfOutput.json")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), "bdf", testData, golden, *update)
}

func TestBdfOutputWithFilter(t *testing.T) {
	var cache = marshal.NewCache()
	cache.Set("filter", map[string]bool{"000111222": true, "000111224": true})
	var golden = filepath.Join("testData", "expectedBdfOutputWithFilter.json")
	marshal.TestParserOutput(t, Parser, cache, "bdf", testData, golden, *update)
}
