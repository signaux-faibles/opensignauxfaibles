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
	marshal.TestParserTupleOutput(t, ParserDebit, cache, "debit", testData, golden, *update)
}

func TestDebitCorrompu(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebitCorrompu.json")
	var testData = filepath.Join("testData", "debitCorrompuTestData.csv")
	marshal.TestParserOutput(t, ParserDebit, cache, "debit", testData, golden, *update)
}

func TestDebitEncoding(t *testing.T) {
	var utf8Input = filepath.Join("testData", "debitUtf8TestData.csv")
	var isoInput = filepath.Join("testData", "debitIsoTestData.csv")
	var outputFromUtf8 = marshal.RunParser(ParserDebit, cache, "debit", utf8Input)
	var outputFromIso = marshal.RunParser(ParserDebit, cache, "debit", isoInput)
	assert.Equal(t, 1, len(outputFromUtf8.Tuples))
	assert.Equal(t, 1, len(outputFromIso.Tuples))
	assert.Equal(t, outputFromUtf8.Tuples[0], outputFromIso.Tuples[0])
}

func TestDelai(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDelai.json")
	var testData = filepath.Join("testData", "delaiTestData.csv")
	marshal.TestParserTupleOutput(t, ParserDelai, cache, "delai", testData, golden, *update)
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsf.json")
	var testData = filepath.Join("testData", "ccsfTestData.csv")
	marshal.TestParserTupleOutput(t, ParserCCSF, cache, "ccsf", testData, golden, *update)
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisation.json")
	var testData = filepath.Join("testData", "cotisationTestData.csv")
	marshal.TestParserTupleOutput(t, ParserCotisation, cache, "cotisation", testData, golden, *update)
}

func TestProcol(t *testing.T) {
	var golden = filepath.Join("testData", "expectedProcol.json")
	var testData = filepath.Join("testData", "procolTestData.csv")
	marshal.TestParserTupleOutput(t, ParserProcol, cache, "procol", testData, golden, *update)
}

func TestEffectif(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEffectif.json")
	var testData = filepath.Join("testData", "effectifTestData.csv")
	cache := marshal.NewCache()
	marshal.TestParserTupleOutput(t, ParserEffectif, cache, "effectif", testData, golden, *update)
}
