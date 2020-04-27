package reporder

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestReporder(t *testing.T) {
	var golden = filepath.Join("testData", "expectedReporder.json")
	var testData = filepath.Join("testData", "reporderTestData.csv")
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "reporder", testData, golden, *update)
}
