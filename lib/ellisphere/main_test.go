package ellisphere

import (
	"flag"
	"path/filepath"
	"testing"

	"opensignauxfaibles/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestEllisphere(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEllisphere.json")
	var testData = filepath.Join("testData", "ellisphereTestData.excel")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)
}
