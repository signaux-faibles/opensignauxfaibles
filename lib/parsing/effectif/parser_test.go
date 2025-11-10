package effectif

import (
	"flag"
	"opensignauxfaibles/lib/engine"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update the expected test values in golden file")

func TestEffectifEnt(t *testing.T) {
	t.Run("Le fichier de test EffectifEnt est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedEffectifEntShort.json")
		var testData = engine.NewBatchFile("testData", "effectifEntTestDataShort.csv")
		cache := engine.NewEmptyCache()

		engine.TestParserOutput(t, NewEffectifEntParser(), cache, testData, golden, *update)
	})

	t.Run("Le fichier de test EffectifEnt est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedEffectifEnt.json")
		var testData = engine.NewBatchFile("testData", "effectifEntTestData.csv")
		cache := engine.NewEmptyCache()

		engine.TestParserOutput(t, NewEffectifEntParser(), cache, testData, golden, *update)
	})

	t.Run("EffectifEnt ne peut pas être importé s'il manque une colonne", func(t *testing.T) {
		output := engine.RunParserInline(t, NewEffectifEntParser(), []string{"siret"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "column siren not found")
	})

	t.Run("EffectifEnt est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := engine.RunParserInline(t, NewEffectifEntParser(), []string{"SiReN"})
		assert.Len(t, output.Reports[0].HeadFatal, 0)
	})
}

func TestEffectif(t *testing.T) {
	var testData = engine.NewBatchFile("testData", "effectifTestData.csv") // Données pour 3 établissements)

	t.Run("Le fichier de test Effectif est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedEffectif.json")
		cache := engine.NewEmptyCache()
		engine.TestParserOutput(t, NewEffectifParser(), cache, testData, golden, *update)
	})

	t.Run("Effectif ne peut pas être importé s'il manque une colonne", func(t *testing.T) {
		output := engine.RunParserInline(t, NewEffectifParser(), []string{"siret"}) // "compte" column is missing
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "column compte not found")
	})

	t.Run("Effectif est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := engine.RunParserInline(t, NewEffectifParser(), []string{"CoMpTe;SiReT"})
		assert.Len(t, output.Reports[0].HeadFatal, 0)
	})
}
