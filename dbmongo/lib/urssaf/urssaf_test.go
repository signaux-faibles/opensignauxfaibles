package urssaf

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestDebit(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebitMD5.csv")
	var testData = filepath.Join("testData", "debitTestData.csv")
	cache := engine.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, parseDebit, cache, "debit", testData, golden, *update)
}

func TestDelai(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDelaiMD5.csv")
	var testData = filepath.Join("testData", "delaiTestData.csv")
	cache := engine.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, parseDelai, cache, "delai", testData, golden, *update)
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsfMD5.csv")
	var testData = filepath.Join("testData", "ccsfTestData.csv")
	cache := engine.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, parseCCSF, cache, "ccsf", testData, golden, *update)
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisationMD5.csv")
	var testData = filepath.Join("testData", "cotisationTestData.csv")
	cache := engine.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, parseCotisation, cache, "cotisation", testData, golden, *update)
}

func TestProcol(t *testing.T) {
	var golden = filepath.Join("testData", "expectedProcolMD5.csv")
	var testData = filepath.Join("testData", "procolTestData.csv")
	cache := engine.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, parseProcol, cache, "procol", testData, golden, *update)
}

func TestEffectif(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEffectifMD5.csv")
	var testData = filepath.Join("testData", "effectifTestData.csv")
	cache := engine.NewCache()
	marshal.TestParserTupleOutput(t, parseEffectif, cache, "effectif", testData, golden, *update)
}
