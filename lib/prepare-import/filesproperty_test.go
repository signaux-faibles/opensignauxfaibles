package prepareimport

import (
	"opensignauxfaibles/lib/engine"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		batchFiles, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]engine.BatchFile{})
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, engine.BatchFiles{}, batchFiles)
	})

	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filename := "sigfaibles_effectif_siret.csv"
		batchFiles, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]engine.BatchFile{engine.NewBatchFile(filename)})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t,
				[]engine.BatchFile{engine.NewBatchFile("sigfaibles_effectif_siret.csv")},
				batchFiles[engine.Effectif])
		}
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		batchFiles, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]engine.BatchFile{engine.NewBatchFile("sigfaibles_debits.csv")})
		expected := engine.BatchFiles{engine.Debit: {engine.NewBatchFile("sigfaibles_debits.csv")}}
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expected, batchFiles)
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		batchFiles, unsupportedFiles :=
			PopulateFilesPropertyFromDataFiles([]engine.BatchFile{
				engine.NewBatchFile("sigfaibles_debits.csv"),
				engine.NewBatchFile("sigfaibles_debits2.csv"),
			})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []engine.BatchFile{engine.NewBatchFile("sigfaibles_debits.csv"),
				engine.NewBatchFile("sigfaibles_debits2.csv")}, batchFiles[engine.Debit])
		}
	})

	t.Run("Should support multiple types of csv files", func(t *testing.T) {
		type File struct {
			Type     engine.ParserType
			Filename string
		}
		files := []File{
			{"effectif", "effectif_dom.csv"},              // --> EFFECTIF
			{"filter", "filter_siren_2002.csv"},           // --> FILTER
			{"sirene_ul", "sireneUL.csv"},                 // --> SIRENE_UL
			{"sirene", "StockEtablissement_utf8_geo.csv"}, // --> SIRENE
		}
		expectedFiles := engine.BatchFiles{}
		inputFiles := []engine.BatchFile{}
		for _, file := range files {
			expectedFiles[file.Type] = append(expectedFiles[file.Type], engine.NewBatchFile(file.Filename))
			inputFiles = append(inputFiles, engine.NewBatchFile(file.Filename))
		}
		resFilesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles(inputFiles)
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expectedFiles, resFilesProperty)
	})

	t.Run("Should not include unsupported files", func(t *testing.T) {
		batchFiles, unsupportedFiles :=
			PopulateFilesPropertyFromDataFiles([]engine.BatchFile{engine.NewBatchFile("coco.csv")})
		assert.Len(t, unsupportedFiles, 1)
		assert.Equal(t, engine.BatchFiles{}, batchFiles)
	})

	t.Run("Should report unsupported files", func(t *testing.T) {
		batchFiles := []engine.BatchFile{engine.NewBatchFileFromBatch(".", dummyBatchKey, "coco.csv")}
		_, unsupportedFiles := PopulateFilesPropertyFromDataFiles(batchFiles)
		assert.Equal(t, []string{path.Join(dummyBatchKey.String(), "coco.csv")}, unsupportedFiles)
	})
}
