package effectif

import (
	"flag"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update the expected test values in golden file")

func TestEffectifEnt(t *testing.T) {
	t.Run("Le fichier de test EffectifEnt est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedEffectifEntShort.json")
		var testData = base.NewBatchFile("testData", "effectifEntTestDataShort.csv")
		cache := engine.NewEmptyCache()

		engine.TestParserOutput(t, NewEffectifEntParser(), cache, testData, golden, *update)
	})

	t.Run("Le fichier de test EffectifEnt est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedEffectifEnt.json")
		var testData = base.NewBatchFile("testData", "effectifEntTestData.csv")
		cache := engine.NewEmptyCache()

		engine.TestParserOutput(t, NewEffectifEntParser(), cache, testData, golden, *update)
	})

	t.Run("EffectifEnt ne peut pas être importé s'il manque une colonne", func(t *testing.T) {
		output := engine.RunParserInline(t, NewEffectifEntParser(), []string{"siret"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Colonne siren non trouvée")
	})

	t.Run("EffectifEnt est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := engine.RunParserInline(t, NewEffectifEntParser(), []string{"SiReN"})
		assert.Len(t, engine.GetFatalErrors(output.Reports[0]), 0)
	})
}
