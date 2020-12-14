package urssaf

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func makeCacheWithComptesMapping() marshal.Cache {
	cache := marshal.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "00000000000000",
			"636043216536562844": "11111111111111",
			"450359886246036238": "22222222222222",
		},
	))
	return cache
}

func TestDebit(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebit.json")
	var testData = filepath.Join("testData", "debitTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserDebit, cache, testData, golden, *update)

	t.Run("Debit n'est importé que si inclus dans le filtre", func(t *testing.T) {
		cache.Set("filter", marshal.SirenFilter{"111111111": true}) // SIREN correspondant à un des 3 comptes retournés par makeCacheWithComptesMapping
		output := marshal.RunParser(ParserDebit, cache, testData)
		// test: tous les tuples retournés concernent le compte associé au SIREN spécifié ci-dessus
		for _, tuple := range output.Tuples {
			debit, _ := tuple.(Debit)
			assert.Equal(t, "636043216536562844", debit.NumeroCompte)
		}
	})

	t.Run("Debit n'est importé que si inclus dans le filtre", func(t *testing.T) {
		cache.Set("filter", marshal.SirenFilter{"111111111": true}) // SIREN correspondant à un des 3 comptes retournés par makeCacheWithComptesMapping
		output := marshal.RunParser(ParserDebit, cache, testData)
		// test: tous les tuples retournés concernent le compte associé au SIREN spécifié ci-dessus
		for _, tuple := range output.Tuples {
			debit, _ := tuple.(Debit)
			assert.Equal(t, "636043216536562844", debit.NumeroCompte)
		}
	})
}

func TestDebitCorrompu(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebitCorrompu.json")
	var testData = filepath.Join("testData", "debitCorrompuTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserDebit, cache, testData, golden, *update)
}

func TestDelai(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDelai.json")
	var testData = filepath.Join("testData", "delaiTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserDelai, cache, testData, golden, *update)
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsf.json")
	var testData = filepath.Join("testData", "ccsfTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserCCSF, cache, testData, golden, *update)
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisation.json")
	var testData = filepath.Join("testData", "cotisationTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserCotisation, cache, testData, golden, *update)

	t.Run("toute ligne de cotisation d'un établissement hors périmètre doit être sautée silencieusement", func(t *testing.T) {
		allowedSiren := "111111111" // SIREN correspondant à un des 3 comptes mentionnés dans le fichier testData
		cache := makeCacheWithComptesMapping()
		cache.Set("filter", marshal.SirenFilter{allowedSiren: true})
		// test
		output := marshal.RunParser(ParserCotisation, cache, testData)
		reportData, _ := output.Events[0].ParseReport()
		assert.Equal(t, false, reportData["isFatal"], "aucune erreur fatale ne doit être rapportée")
		assert.Equal(t, []interface{}{}, reportData["headRejected"], "aucune erreur de parsing ne doit être rapportée")
		assert.Equal(t, 1.0, reportData["linesValid"], "seule la ligne de cotisation liée à un établissement du périmètre doit être incluse")
	})

	t.Run("toute ligne de cotisation d'un établissement non inclus dans les comptes urssaf doit être sautée silencieusement", func(t *testing.T) {
		cache := marshal.NewCache()
		cache.Set("comptes", marshal.MockComptesMapping(
			map[string]string{
				"111982477292496174": "00000000000000",
				// "636043216536562844": "11111111111111", // on retire volontairement ce mapping qui va être demandé par le parseur de cotisations
				"450359886246036238": "22222222222222",
			},
		))
		// test
		output := marshal.RunParser(ParserCotisation, cache, testData)
		reportData, _ := output.Events[0].ParseReport()
		assert.Equal(t, false, reportData["isFatal"], "aucune erreur fatale ne doit être rapportée")
		assert.Equal(t, []interface{}{}, reportData["headRejected"], "aucune erreur de parsing ne doit être rapportée")
		assert.Equal(t, 1.0, reportData["linesSkipped"], "seule la ligne de cotisation liée à un établissement hors mapping doit être sautée")
	})
}

func TestProcol(t *testing.T) {
	var golden = filepath.Join("testData", "expectedProcol.json")
	var testData = filepath.Join("testData", "procolTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserProcol, cache, testData, golden, *update)
}

func TestEffectif(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEffectif.json")
	var testData = filepath.Join("testData", "effectifTestData.csv") // Données pour 3 établissements
	cache := marshal.NewCache()
	marshal.TestParserOutput(t, ParserEffectif, cache, testData, golden, *update)

	t.Run("Effectif n'est importé que si inclus dans le filtre", func(t *testing.T) {
		allowedSiren := "149285238" // SIREN correspondant à un des 3 SIRETs mentionnés dans le fichier
		cache := marshal.NewCache()
		cache.Set("filter", marshal.SirenFilter{allowedSiren: true})
		output := marshal.RunParser(ParserEffectif, cache, testData)
		// test: vérifier que tous les tuples retournés concernent ce SIREN
		for _, tuple := range output.Tuples {
			effectif, _ := tuple.(Effectif)
			assert.Equal(t, allowedSiren, effectif.Siret[0:9])
		}
	})
}

func TestEffectifEnt(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEffectifEnt.json")
	var testData = filepath.Join("testData", "effectifEntTestData.csv")
	cache := marshal.NewCache()
	marshal.TestParserOutput(t, ParserEffectifEnt, cache, testData, golden, *update)
}
