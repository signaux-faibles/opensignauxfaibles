package apdemande

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApdemande(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApdemande.json")
	var testData = filepath.Join("testData", "apdemandeTestData.csv")
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "apdemande", testData, golden, *update)
}
