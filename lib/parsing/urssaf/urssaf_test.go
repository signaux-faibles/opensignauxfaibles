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

var notFoundRegexp = "column [A-Za-z]+ not found"

var update = flag.Bool("update", false, "update the expected test values in golden file")

func TestUrssaf(t *testing.T) {
	t.Run("URSSAF gzipped files can be decompressed on the fly", func(t *testing.T) {
		type TestCase struct {
			Parser     engine.Parser
			InputFile  string
			GoldenFile string
		}
		urssafFiles := []TestCase{
			{NewCCSFParser(), "ccsfTestData.csv", "expectedCcsf.json"},
			{NewDebitParser(), "debitTestData.csv", "expectedDebit.json"},
			{NewDelaiParser(), "delaiTestData.csv", "expectedDelai.json"},
			{NewProcolParser(), "procolTestData.csv", "expectedProcol.json"},
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

				engine.TestParserOutput(t, testCase.Parser, compressedFilePath, tmpGoldenFile.Name(), false)
			})
		}
	})
}

func TestDebit(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebit.json")
	var testData = engine.NewBatchFile("testData", "debitTestData.csv")

	engine.TestParserOutput(t, NewDebitParser(), testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewDebitParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})
}

func TestDebitCorrompu(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDebitCorrompu.json")
	var testData = engine.NewBatchFile("testData", "debitCorrompuTestData.csv")
	engine.TestParserOutput(t, NewDebitParser(), testData, golden, *update)
}

func TestDelai(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDelai.json")
	var testData = engine.NewBatchFile("testData", "delaiTestData.csv")
	engine.TestParserOutput(t, NewDelaiParser(), testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewDelaiParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})
}

func TestCcsf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCcsf.json")
	var testData = engine.NewBatchFile("testData", "ccsfTestData.csv")
	engine.TestParserOutput(t, NewCCSFParser(), testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewCCSFParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})
}

func TestCotisation(t *testing.T) {
	var golden = filepath.Join("testData", "expectedCotisation.json")
	var testData = engine.NewBatchFile("testData", "cotisationTestData.csv")
	engine.TestParserOutput(t, NewCotisationParser(), testData, golden, *update)

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewCotisationParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, notFoundRegexp, engine.GetFatalError(output))
	})
}

func TestProcol(t *testing.T) {
	t.Run("Le fichier de test Procol est parsé comme d'habitude", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedProcol.json")
		var testData = engine.NewBatchFile("testData", "procolTestData.csv")
		engine.TestParserOutput(t, NewProcolParser(), testData, golden, *update)
	})

	t.Run("should report fatal error when column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewProcolParser(), []string{"dummy"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "not found")
	})

	t.Run("est insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		output := engine.RunParserInline(t, NewProcolParser(), []string{"dT_eFfeT;lIb_aCtx_stDx;sIret"})
		assert.Len(t, output.Reports[0].HeadFatal, 0)
	})
}
