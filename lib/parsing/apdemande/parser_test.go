package apdemande

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApdemande(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApdemande.json")
	var testData = base.NewBatchFile("testData", "apdemandeTestData.csv")
	engine.TestParserOutput(t, NewParserApdemande(), engine.NewEmptyCache(), testData, golden, *update)

	t.Run("should fail if one column misses", func(t *testing.T) {
		output := engine.RunParserInline(t, NewParserApdemande(), []string{"ID_DA,ETAB_SIRET"}) // EFF_ENT is missing (among others)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Fatal: Colonne EFF_ENT non trouvée. Abandon.")
	})

	t.Run("should fail if a composite column misses", func(t *testing.T) {
		headerRow := []string{"ID_DA,ETAB_SIRET,EFF_ENT,EFF_ETAB,DATE_STATUT,HTA,EFF_AUTO,MOTIF_RECOURS_SE,S_HEURE_CONSOM_TOT,S_HEURE_CONSOM_TOT,DATE_FIN"} // DATE_DEB is missing
		output := engine.RunParserInline(t, NewParserApdemande(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Fatal: Colonne DATE_DEB non trouvée. Abandon.")
	})
}
