package urssaf

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestDebit(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebit.json")
	var testData = filepath.Join("testData", "debitTestData.csv")
	cache := base.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, ParserDebit, cache, "debit", testData, golden, *update)
}

func TestDelai(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDelai.json")
	var testData = filepath.Join("testData", "delaiTestData.csv")
	cache := base.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, ParserDelai, cache, "delai", testData, golden, *update)
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsf.json")
	var testData = filepath.Join("testData", "ccsfTestData.csv")
	cache := base.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, ParserCCSF, cache, "ccsf", testData, golden, *update)
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisation.json")
	var testData = filepath.Join("testData", "cotisationTestData.csv")
	cache := base.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, ParserCotisation, cache, "cotisation", testData, golden, *update)
}

func TestProcol(t *testing.T) {
	var golden = filepath.Join("testData", "expectedProcol.json")
	var testData = filepath.Join("testData", "procolTestData.csv")
	cache := base.NewCache()
	cache.Set("comptes", marshal.MockComptesMapping(
		map[string]string{
			"111982477292496174": "000000000000000",
			"636043216536562844": "111111111111111",
			"450359886246036238": "222222222222222",
		},
	))

	marshal.TestParserTupleOutput(t, ParserProcol, cache, "procol", testData, golden, *update)
}

func TestEffectif(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEffectif.json")
	var testData = filepath.Join("testData", "effectifTestData.csv")
	cache := base.NewCache()
	marshal.TestParserTupleOutput(t, ParserEffectif, cache, "effectif", testData, golden, *update)
}
