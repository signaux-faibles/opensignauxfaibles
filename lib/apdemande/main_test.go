package apdemande

import (
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApdemande(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApdemande.json")
	var testData = filepath.Join("testData", "apdemandeTestData.csv")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)

	t.Run("should fail if one column misses", func(t *testing.T) {
		csvData := strings.Join([]string{"ID_DA,ETAB_SIRET"}, "\n") // EFF_ENT is missing (among others)
		csvFile := marshal.CreateTempFileWithContent(t, []byte(csvData))
		output := marshal.RunParser(Parser, marshal.NewCache(), csvFile.Name())
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, 1, len(output.Events), "should return a parsing report")
		reportData, _ := output.Events[0].ParseReport()
		assert.Equal(t, true, reportData["isFatal"], "should report a fatal error")
		assert.Contains(t, reportData["headFatal"], "Fatal: Colonne EFF_ENT non trouvée. Abandon.")
	})

	t.Run("should fail if a composite column misses", func(t *testing.T) {
		csvData := strings.Join([]string{"ID_DA,ETAB_SIRET,EFF_ENT,EFF_ETAB,DATE_STATUT,HTA,EFF_AUTO,MOTIF_RECOURS_SE,S_HEURE_CONSOM_TOT,S_HEURE_CONSOM_TOT,DATE_FIN"}, "\n") // DATE_DEB is missing
		csvFile := marshal.CreateTempFileWithContent(t, []byte(csvData))
		output := marshal.RunParser(Parser, marshal.NewCache(), csvFile.Name())
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, 1, len(output.Events), "should return a parsing report")
		reportData, _ := output.Events[0].ParseReport()
		assert.Equal(t, true, reportData["isFatal"], "should report a fatal error")
		assert.Contains(t, reportData["headFatal"], "Fatal: Colonne DATE_DEB non trouvée. Abandon.")
	})
}
