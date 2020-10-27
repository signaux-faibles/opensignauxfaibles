package urssaf

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func makeCacheWithComptesMapping() marshal.Cache {
	cache := marshal.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))
	return cache
}

var cache = makeCacheWithComptesMapping()

func TestDebit(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebit.json")
	var testData = filepath.Join("testData", "debitTestData.csv")
	marshal.TestParserOutput(t, ParserDebit, cache, testData, golden, *update)
}

func TestDebitCorrompu(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebitCorrompu.json")
	var testData = filepath.Join("testData", "debitCorrompuTestData.csv")
	marshal.TestParserOutput(t, ParserDebit, cache, testData, golden, *update)
}

func TestDelai(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDelai.json")
	var testData = filepath.Join("testData", "delaiTestData.csv")
	marshal.TestParserOutput(t, ParserDelai, cache, testData, golden, *update)
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsf.json")
	var testData = filepath.Join("testData", "ccsfTestData.csv")
	marshal.TestParserOutput(t, ParserCCSF, cache, testData, golden, *update)
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisation.json")
	var testData = filepath.Join("testData", "cotisationTestData.csv")
	marshal.TestParserOutput(t, ParserCotisation, cache, testData, golden, *update)
}

func TestProcol(t *testing.T) {
	var golden = filepath.Join("testData", "expectedProcol.json")
	var testData = filepath.Join("testData", "procolTestData.csv")
	marshal.TestParserOutput(t, ParserProcol, cache, testData, golden, *update)
}

func TestEffectif(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEffectif.json")
	var testData = filepath.Join("testData", "effectifTestData.csv") // Données pour 3 établissements
	cache := marshal.NewCache()
	marshal.TestParserOutput(t, ParserEffectif, cache, testData, golden, *update)

	t.Run("Effectif n'est importé que si inclus dans le filtre", func(t *testing.T) {
		cache := marshal.NewCache()
		cache.Set("filter", marshal.SirenFilter{"149285238": true}) // SIREN correspondant à un des 3 SIRETs mentionnés dans le fichier
		output := marshal.RunParser(ParserEffectif, cache, testData)
		assert.Equal(t, 111, len(output.Tuples))
	})
}

func TestEffectifEnt(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEffectifEnt.json")
	var testData = filepath.Join("testData", "effectifEntTestData.csv")
	cache := marshal.NewCache()
	marshal.TestParserOutput(t, ParserEffectifEnt, cache, testData, golden, *update)
}
