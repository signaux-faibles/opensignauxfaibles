package sirene

import (
	"flag"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSirene.json")
var testData = base.NewBatchFile("testData", "sireneTestData.csv")

func TestSirene(t *testing.T) {
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)

	t.Run("should fail if a required column is missing", func(t *testing.T) {
		output := marshal.RunParserInline(t, Parser, []string{"siren"}) // many columns are missing
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, regexp.MustCompile("Colonne [^ ]+ non trouvée"), marshal.GetFatalError(output))
	})
}
