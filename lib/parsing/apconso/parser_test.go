package apconso

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/engine"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApconso(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApconso.json")
	var testData = engine.NewBatchFile("testData", "apconsoTestData.csv")
	engine.TestParserOutput(t, NewApconsoParser(), testData, golden, *update)

	t.Run("doit détecter s'il manque des colonnes", func(t *testing.T) {
		output := engine.RunParserInline(t, NewApconsoParser(), []string{"ID_DA,ETAB_SIRET,MOIS,HEURE,MONTANTS,EFFECTIFS"}) // typo: HEURE au lieu de HEURES
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "le parseur doit retourner aucun tuple")
		assert.Contains(t, engine.GetFatalError(output), "column HEURES not found")
	})
}
