package urssaf

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/engine"
)

var update = flag.Bool("update", false, "update the expected test values in golden file")

func makeCacheWithComptesMapping() engine.Cache {
	cache := engine.NewEmptyCache()
	cache.Set("comptes", MockComptesMapping(
		map[string]string{
			"111982477292496174": "00000000000000",
			"636043216536562844": "11111111111111",
			"450359886246036238": "22222222222222",
		},
	))
	return cache
}

func TestUrssaf(t *testing.T) {
	t.Run("URSSAF gzipped files can be decompressed on the fly", func(t *testing.T) {
		type TestCase struct {
			Parser     engine.Parser
			InputFile  string
			GoldenFile string
			Cache      engine.Cache
		}
		urssafFiles := []TestCase{
			{NewCCSFParser(), "ccsfTestData.csv", "expectedCcsf.json", makeCacheWithComptesMapping()},
			{NewDebitParser(), "debitTestData.csv", "expectedDebit.json", makeCacheWithComptesMapping()},
			{NewDelaiParser(), "delaiTestData.csv", "expectedDelai.json", makeCacheWithComptesMapping()},
			{NewProcolParser(), "procolTestData.csv", "expectedProcol.json", makeCacheWithComptesMapping()},
		}
		for _, testCase := range urssafFiles {
			t.Run(string(testCase.Parser.Type()), func(t *testing.T) {
				// Compression du fichier de données
				err := exec.Command("gzip", "--keep", filepath.Join("testData", testCase.InputFile)).Run() // créée une version gzippée du fichier
				assert.NoError(t, err)
				compressedFilePath := engine.NewBatchFile("testData", testCase.InputFile+".gz")
				t.Cleanup(func() { os.Remove(compressedFilePath.Path()) })

				// Création d'un fichier Golden temporaire mentionnant le nom du fichier compressé
				initialGoldenContent, err := os.ReadFile(filepath.Join("testData", testCase.GoldenFile))
				assert.NoError(t, err)
				goldenContent := bytes.ReplaceAll(initialGoldenContent, []byte(testCase.InputFile), []byte(testCase.InputFile+".gz"))
				tmpGoldenFile := engine.CreateTempFileWithContent(t, goldenContent)

				engine.TestParserOutput(t, testCase.Parser, testCase.Cache, compressedFilePath, tmpGoldenFile.Name(), false)
			})
		}
	})
}

func TestDebit(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebit.json")
	var testData = engine.NewBatchFile("testData", "debitTestData.csv")
	var cache = makeCacheWithComptesMapping()

	engine.TestParserOutput(t, NewDebitParser(), cache, testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInlineEx(t, cache, NewDebitParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Colonne num_cpte non trouvée")
	})

	// t.Run("Debit n'est importé que si inclus dans le filtre", func(t *testing.T) {
	// 	cache.Set("filter", engine.SirenFilter{"111111111": true}) // SIREN correspondant à un des 3 comptes retournés par makeCacheWithComptesMapping
	// 	output := engine.RunParser(NewDebitParser(), cache, testData)
	// 	// test: tous les tuples retournés concernent le compte associé au SIREN spécifié ci-dessus
	// 	for _, tuple := range output.Tuples {
	// 		debit, _ := tuple.(Debit)
	// 		assert.Equal(t, "636043216536562844", debit.NumeroCompte)
	// 	}
	// })

	// t.Run("Debit n'est importé que si inclus dans le filtre", func(t *testing.T) {
	// 	cache.Set("filter", engine.SirenFilter{"111111111": true}) // SIREN correspondant à un des 3 comptes retournés par makeCacheWithComptesMapping
	// 	output := engine.RunParser(NewDebitParser(), cache, testData)
	// 	// test: tous les tuples retournés concernent le compte associé au SIREN spécifié ci-dessus
	// 	for _, tuple := range output.Tuples {
	// 		debit, _ := tuple.(Debit)
	// 		assert.Equal(t, "636043216536562844", debit.NumeroCompte)
	// 	}
	// })
}

func TestDebitCorrompu(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebitCorrompu.json")
	var testData = engine.NewBatchFile("testData", "debitCorrompuTestData.csv")
	var cache = makeCacheWithComptesMapping()
	engine.TestParserOutput(t, NewDebitParser(), cache, testData, golden, *update)
}

func TestDelai(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDelai.json")
	var testData = engine.NewBatchFile("testData", "delaiTestData.csv")
	var cache = makeCacheWithComptesMapping()
	engine.TestParserOutput(t, NewDelaiParser(), cache, testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInlineEx(t, cache, NewDelaiParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Colonne Numero_compte_externe non trouvée")
	})
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsf.json")
	var testData = engine.NewBatchFile("testData", "ccsfTestData.csv")
	var cache = makeCacheWithComptesMapping()
	engine.TestParserOutput(t, NewCCSFParser(), cache, testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInlineEx(t, cache, NewCCSFParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Colonne Compte non trouvée")
	})
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisation.json")
	var testData = engine.NewBatchFile("testData", "cotisationTestData.csv")
	var cache = makeCacheWithComptesMapping()
	engine.TestParserOutput(t, NewCotisationParser(), cache, testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInlineEx(t, cache, NewCotisationParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "Colonne Compte non trouvée")
	})

	// t.Run("toute ligne de cotisation d'un établissement hors périmètre doit être sautée silencieusement", func(t *testing.T) {
	// 	allowedSiren := "111111111" // SIREN correspondant à un des 3 comptes mentionnés dans le fichier testData
	// 	cache := makeCacheWithComptesMapping()
	// 	cache.Set("filter", engine.SirenFilter{allowedSiren: true})
	// 	// test
	// 	output := engine.RunParser(NewCotisationParser(), cache, testData)
	// 	reportData := output.Reports[0]
	// 	assert.Equal(t, false, reportData.IsFatal, "aucune erreur fatale ne doit être rapportée")
	// 	assert.Equal(t, []string{}, reportData.HeadRejected, "aucune erreur de parsing ne doit être rapportée")
	// 	assert.Equal(t, int64(1.0), reportData.LinesValid, "seule la ligne de cotisation liée à un établissement du périmètre doit être incluse")
	// })

	t.Run("toute ligne de cotisation d'un établissement non inclus dans les comptes urssaf doit être sautée silencieusement", func(t *testing.T) {
		cache := engine.NewEmptyCache()
		cache.Set("comptes", MockComptesMapping(
			map[string]string{
				"111982477292496174": "00000000000000",
				// "636043216536562844": "11111111111111", // on retire volontairement ce mapping qui va être demandé par le parseur de cotisations
				"450359886246036238": "22222222222222",
			},
		))
		// test
		output := engine.RunParser(NewCotisationParser(), cache, testData)
		report := output.Reports[0]
		assert.Equal(t, false, report.IsFatal, "aucune erreur fatale ne doit être rapportée")
		assert.Equal(t, []string{}, report.HeadRejected, "aucune erreur de parsing ne doit être rapportée")
		assert.Equal(t, int64(1.0), report.LinesSkipped, "seule la ligne de cotisation liée à un établissement hors mapping doit être sautée")
	})
}

func TestProcol(t *testing.T) {
	var cache = makeCacheWithComptesMapping()

	t.Run("Le fichier de test Procol est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedProcol.json")
		var testData = engine.NewBatchFile("testData", "procolTestData.csv")
		engine.TestParserOutput(t, NewProcolParser(), cache, testData, golden, *update)
	})

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInlineEx(t, cache, NewProcolParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "non trouvée")
	})

	t.Run("est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := engine.RunParserInlineEx(t, cache, NewProcolParser(), []string{"dT_eFfeT;lIb_aCtx_stDx;sIret"})
		assert.Len(t, output.Reports[0].HeadFatal, 0)
	})
}
