package prepareimport

import (
	"opensignauxfaibles/lib/base"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		batchFiles, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]base.BatchFile{})
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, base.BatchFiles{}, batchFiles)
	})

	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filename := "sigfaibles_effectif_siret.csv"
		batchFiles, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]base.BatchFile{base.NewDummyBatchFile(filename)})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t,
				[]base.BatchFile{base.NewDummyBatchFile("sigfaibles_effectif_siret.csv")},
				batchFiles[base.Effectif])
		}
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		batchFiles, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]base.BatchFile{base.NewDummyBatchFile("sigfaibles_debits.csv")})
		expected := base.BatchFiles{base.Debit: {base.NewDummyBatchFile("sigfaibles_debits.csv")}}
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expected, batchFiles)
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		batchFiles, unsupportedFiles :=
			PopulateFilesPropertyFromDataFiles([]base.BatchFile{
				base.NewDummyBatchFile("sigfaibles_debits.csv"),
				base.NewDummyBatchFile("sigfaibles_debits2.csv"),
			})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []base.BatchFile{base.NewDummyBatchFile("sigfaibles_debits.csv"),
				base.NewDummyBatchFile("sigfaibles_debits2.csv")}, batchFiles[base.Debit])
		}
	})

	t.Run("Should support multiple types of csv files", func(t *testing.T) {
		type File struct {
			Type     base.ParserType
			Filename string
		}
		files := []File{
			{"diane", "diane_req_2002.csv"},               // --> DIANE
			{"diane", "diane_req_dom_2002.csv"},           // --> DIANE
			{"effectif", "effectif_dom.csv"},              // --> EFFECTIF
			{"filter", "filter_siren_2002.csv"},           // --> FILTER
			{"sirene_ul", "sireneUL.csv"},                 // --> SIRENE_UL
			{"sirene", "StockEtablissement_utf8_geo.csv"}, // --> SIRENE
		}
		expectedFiles := base.BatchFiles{}
		inputFiles := []base.BatchFile{}
		for _, file := range files {
			expectedFiles[file.Type] = append(expectedFiles[file.Type], base.NewDummyBatchFile(file.Filename))
			inputFiles = append(inputFiles, base.NewDummyBatchFile(file.Filename))
		}
		resFilesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles(inputFiles)
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expectedFiles, resFilesProperty)
	})

	t.Run("Should not include unsupported files", func(t *testing.T) {
		batchFiles, unsupportedFiles :=
			PopulateFilesPropertyFromDataFiles([]base.BatchFile{base.NewDummyBatchFile("coco.csv")})
		assert.Len(t, unsupportedFiles, 1)
		assert.Equal(t, base.BatchFiles{}, batchFiles)
	})

	t.Run("Should report unsupported files", func(t *testing.T) {
		batchFiles := []base.BatchFile{base.NewDummyBatchFileFromBatch("coco.csv", dummyBatchKey)}
		_, unsupportedFiles := PopulateFilesPropertyFromDataFiles(batchFiles)
		assert.Equal(t, []string{path.Join(dummyBatchKey.String(), "coco.csv")}, unsupportedFiles)
	})
}
