package sirene

import (
	"flag"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSirene.json")
var testData = filepath.Join("testData", "sireneTestData.csv")

func TestSirene(t *testing.T) {
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)

	t.Run("should fail if a required column is missing", func(t *testing.T) {
		output := marshal.RunParserInline(t, Parser, []string{"siren"}) // many columns are missing
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, 1, len(output.Events), "should return a parsing report")
		reportData, _ := output.Events[0].ParseReport()
		assert.Equal(t, true, reportData["isFatal"], "should report a fatal error")
		assert.Regexp(t, regexp.MustCompile("Colonne [^ ]+ non trouvée"), reportData["headFatal"])
	})
}
