package urssaf

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/marshal"
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

func TestUrssaf(t *testing.T) {
	t.Run("Les fichiers urssaf gzippés peuvent être décompressés à la volée", func(t *testing.T) {
		type TestCase struct {
			Parser     marshal.Parser
			InputFile  string
			GoldenFile string
			Cache      marshal.Cache
		}
		urssafFiles := []TestCase{
			{ParserCCSF, "ccsfTestData.csv", "expectedCcsf.json", makeCacheWithComptesMapping()},
			{ParserCompte, "comptesTestData.csv", "expectedComptes.json", marshal.NewCache()},
			{ParserDebit, "debitTestData.csv", "expectedDebit.json", makeCacheWithComptesMapping()},
			{ParserDelai, "delaiTestData.csv", "expectedDelai.json", makeCacheWithComptesMapping()},
			{ParserEffectifEnt, "effectifEntTestData.csv", "expectedEffectifEnt.json", makeCacheWithComptesMapping()},
			{ParserEffectif, "effectifTestData.csv", "expectedEffectif.json", makeCacheWithComptesMapping()},
			{ParserProcol, "procolTestData.csv", "expectedProcol.json", makeCacheWithComptesMapping()},
		}
		for _, testCase := range urssafFiles {
			t.Run(testCase.Parser.Type(), func(t *testing.T) {
				// Compression du fichier de données
				err := exec.Command("gzip", "--keep", filepath.Join("testData", testCase.InputFile)).Run() // créée une version gzippée du fichier
				assert.NoError(t, err)
				compressedFilePath := filepath.Join("testData", testCase.InputFile+".gz")
				t.Cleanup(func() { os.Remove(compressedFilePath) })
				// Création d'un fichier Golden temporaire mentionnant le nom du fichier compressé
				initialGoldenContent, err := os.ReadFile(filepath.Join("testData", testCase.GoldenFile))
				assert.NoError(t, err)
				goldenContent := bytes.ReplaceAll(initialGoldenContent, []byte(testCase.InputFile), []byte(testCase.InputFile+".gz"))
				tmpGoldenFile := marshal.CreateTempFileWithContent(t, goldenContent)
				marshal.TestParserOutput(t, testCase.Parser, testCase.Cache, compressedFilePath, tmpGoldenFile.Name(), false)
			})
		}
	})
}

func TestComptes(t *testing.T) {

	t.Run("Le fichier de test Comptes est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedComptes.json")
		var testData = filepath.Join("testData", "comptesTestData.csv")
		marshal.TestParserOutput(t, ParserCompte, marshal.NewCache(), testData, golden, *update)
	})
}

func TestDebit(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebit.json")
	var testData = filepath.Join("testData", "debitTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserDebit, cache, testData, golden, *update)

	t.Run("doit rapporter une erreur fatale s'il manque une colonne", func(t *testing.T) {
		output := marshal.RunParserInlineEx(t, cache, ParserDebit, []string{"dummy"})
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, marshal.GetFatalError(output), "Colonne num_cpte non trouvée")
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

	t.Run("doit rapporter une erreur fatale s'il manque une colonne", func(t *testing.T) {
		output := marshal.RunParserInlineEx(t, cache, ParserDelai, []string{"dummy"})
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, marshal.GetFatalError(output), "Colonne Numero_compte_externe non trouvée")
	})
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsf.json")
	var testData = filepath.Join("testData", "ccsfTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserCCSF, cache, testData, golden, *update)

	t.Run("doit rapporter une erreur fatale s'il manque une colonne", func(t *testing.T) {
		output := marshal.RunParserInlineEx(t, cache, ParserCCSF, []string{"dummy"})
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, marshal.GetFatalError(output), "Colonne Compte non trouvée")
	})
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisation.json")
	var testData = filepath.Join("testData", "cotisationTestData.csv")
	var cache = makeCacheWithComptesMapping()
	marshal.TestParserOutput(t, ParserCotisation, cache, testData, golden, *update)

	t.Run("doit rapporter une erreur fatale s'il manque une colonne", func(t *testing.T) {
		output := marshal.RunParserInlineEx(t, cache, ParserCotisation, []string{"dummy"})
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, marshal.GetFatalError(output), "Colonne Compte non trouvée")
	})

	t.Run("toute ligne de cotisation d'un établissement hors périmètre doit être sautée silencieusement", func(t *testing.T) {
		allowedSiren := "111111111" // SIREN correspondant à un des 3 comptes mentionnés dans le fichier testData
		cache := makeCacheWithComptesMapping()
		cache.Set("filter", marshal.SirenFilter{allowedSiren: true})
		// test
		output := marshal.RunParser(ParserCotisation, cache, testData)
		reportData := output.Events[0].Report
		assert.Equal(t, false, reportData.IsFatal, "aucune erreur fatale ne doit être rapportée")
		assert.Equal(t, []string{}, reportData.HeadRejected, "aucune erreur de parsing ne doit être rapportée")
		assert.Equal(t, int64(1.0), reportData.LinesValid, "seule la ligne de cotisation liée à un établissement du périmètre doit être incluse")
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
		report := output.Events[0].Report
		assert.Equal(t, false, report.IsFatal, "aucune erreur fatale ne doit être rapportée")
		assert.Equal(t, []string{}, report.HeadRejected, "aucune erreur de parsing ne doit être rapportée")
		assert.Equal(t, int64(1.0), report.LinesSkipped, "seule la ligne de cotisation liée à un établissement hors mapping doit être sautée")
	})
}

func TestProcol(t *testing.T) {
	var cache = makeCacheWithComptesMapping()

	t.Run("Le fichier de test Procol est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedProcol.json")
		var testData = filepath.Join("testData", "procolTestData.csv")
		marshal.TestParserOutput(t, ParserProcol, cache, testData, golden, *update)
	})

	t.Run("doit rapporter une erreur fatale s'il manque une colonne", func(t *testing.T) {
		output := marshal.RunParserInlineEx(t, cache, ParserProcol, []string{"dummy"})
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, marshal.GetFatalError(output), "non trouvée")
	})

	t.Run("est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := marshal.RunParserInline(t, ParserProcol, []string{"dT_eFfeT;lIb_aCtx_stDx;sIret"})
		assert.Len(t, marshal.GetFatalErrors(output.Events[0]), 0)
	})
}

func TestEffectif(t *testing.T) {
	var testData = filepath.Join("testData", "effectifTestData.csv") // Données pour 3 établissements

	t.Run("Le fichier de test Effectif est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedEffectif.json")
		cache := marshal.NewCache()
		marshal.TestParserOutput(t, ParserEffectif, cache, testData, golden, *update)
	})

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

	t.Run("Effectif ne peut pas être importé s'il manque une colonne", func(t *testing.T) {
		output := marshal.RunParserInline(t, ParserEffectif, []string{"siret"}) // "compte" column is missing
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, marshal.GetFatalError(output), "Colonne compte non trouvée")
	})

	t.Run("Effectif est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := marshal.RunParserInline(t, ParserEffectif, []string{"CoMpTe;SiReT"})
		assert.Len(t, marshal.GetFatalErrors(output.Events[0]), 0)
	})
}

func TestEffectifEnt(t *testing.T) {
	t.Run("Le fichier de test EffectifEnt est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedEffectifEnt.json")
		var testData = filepath.Join("testData", "effectifEntTestData.csv")
		cache := marshal.NewCache()
		marshal.TestParserOutput(t, ParserEffectifEnt, cache, testData, golden, *update)
	})

	t.Run("EffectifEnt ne peut pas être importé s'il manque une colonne", func(t *testing.T) {
		output := marshal.RunParserInline(t, ParserEffectifEnt, []string{"siret"})
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, marshal.GetFatalError(output), "Colonne siren non trouvée")
	})

	t.Run("EffectifEnt est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := marshal.RunParserInline(t, ParserEffectifEnt, []string{"SiReN"})
		assert.Len(t, marshal.GetFatalErrors(output.Events[0]), 0)
	})
}
