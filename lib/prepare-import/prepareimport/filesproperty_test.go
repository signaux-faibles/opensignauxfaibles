package prepareimport

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{}, dummyBatchKey)
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, FilesProperty{}, filesProperty)
	})

	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"sigfaibles_effectif_siret.csv", ""}}, dummyBatchKey)
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []BatchFile{dummyBatchFile("sigfaibles_effectif_siret.csv")}, filesProperty[effectif])
		}
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"sigfaibles_debits.csv", ""}}, dummyBatchKey)
		expected := FilesProperty{debit: {dummyBatchFile("sigfaibles_debits.csv")}}
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expected, filesProperty)
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"sigfaibles_debits.csv", ""}, SimpleDataFile{"sigfaibles_debits2.csv", ""}}, dummyBatchKey)
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []BatchFile{dummyBatchFile("sigfaibles_debits.csv"), dummyBatchFile("sigfaibles_debits2.csv")}, filesProperty[debit])
		}
	})

	t.Run("Should support multiple types of csv files", func(t *testing.T) {
		type File struct {
			Type     ValidFileType
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
		expectedFiles := FilesProperty{}
		inputFiles := []DataFile{}
		for _, file := range files {
			expectedFiles[file.Type] = append(expectedFiles[file.Type], dummyBatchFile(file.Filename))
			inputFiles = append(inputFiles, SimpleDataFile{file.Filename, ""})
		}
		resFilesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles(inputFiles, dummyBatchKey)
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expectedFiles, resFilesProperty)
	})

	t.Run("Should not include unsupported files", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"coco.csv", ""}}, dummyBatchKey)
		assert.Len(t, unsupportedFiles, 1)
		assert.Equal(t, FilesProperty{}, filesProperty)
	})

	t.Run("Should report unsupported files", func(t *testing.T) {
		_, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"coco.csv", ""}}, dummyBatchKey)
		assert.Equal(t, []string{dummyBatchKey.Path() + "coco.csv"}, unsupportedFiles)
	})

	t.Run("Should skip subdirectories", func(t *testing.T) {
		subBatch := newSafeBatchKey("1803_01")
		parentBatch := subBatch.GetParentBatch()
		parentDir := CreateTempFiles(t, newSafeBatchKey(parentBatch), []string{})
		subBatchDir := filepath.Join(parentDir, parentBatch, subBatch.String())
		_ = os.Mkdir(subBatchDir, 0777)
		parentFilesProperty, unsupportedFiles := PopulateFilesProperty(parentDir, newSafeBatchKey(parentBatch))
		assert.Equal(t, []string{}, unsupportedFiles)
		assert.Equal(t, FilesProperty{}, parentFilesProperty)
	})

	t.Run("Should add a 'gzip:' prefix to compressed files", func(t *testing.T) {
		bytes, err := GzipString(SomeText(254781489))
		if err != nil {
			t.Errorf("erreur pendant la compression d'un texte aléatoire : %s", err)
		}
		dir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"sigfaibles_debits.csv.gz": bytes,
		})
		resFilesProperty, _ := PopulateFilesProperty(dir, dummyBatchKey)
		assert.Len(t, resFilesProperty["debit"], 1)
		actualFilePath := resFilesProperty["debit"][0].Path() // cf batchFile.MarshalJSON()
		assert.Equal(t, "gzip:/1802/sigfaibles_debits.csv.gz", actualFilePath)
	})
}
