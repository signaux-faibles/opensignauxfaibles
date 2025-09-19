package prepareimport

import (
	"bytes"
	"compress/gzip"
	"errors"
	"opensignauxfaibles/lib/base"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadFilenames(t *testing.T) {
	t.Run("Should return filenames in a directory", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"tmpfile"})
		filenames, err := ReadFilenames(path.Join(dir, dummyBatchKey.String()))
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, []string{"tmpfile"}, filenames)
	})
}

func TestPrepareImport(t *testing.T) {
	t.Run("Should warn if the batch was not found in the specified directory", func(t *testing.T) {
		wantedBatch := newSafeBatchKey("1803") // different of dummyBatchKey
		parentDir := CreateTempFiles(t, dummyBatchKey, []string{})
		_, err := PrepareImport(parentDir, wantedBatch)
		expected := "could not find directory 1803 in provided path"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if the sub-batch was not found in the specified directory", func(t *testing.T) {
		subBatch := newSafeBatchKey("1803_01")
		parentBatch := newSafeBatchKey("1803")
		parentDir := CreateTempFiles(t, parentBatch, []string{})
		_, err := PrepareImport(parentDir, subBatch)
		expected := "could not find directory 1803/1803_01 in provided path"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if no filter is provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"sigfaibles_debits.csv"})
		_, err := PrepareImport(dir, dummyBatchKey)
		expected := "filter is missing: batch should include a filter or one effectif file"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if 2 effectif files are provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"sigfaible_effectif_siret.csv", "sigfaible_effectif_siret2.csv"})
		_, err := PrepareImport(dir, dummyBatchKey)
		expected := "filter is missing: batch should include a filter or one effectif file"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if neither effectif and date_fin_effectif are provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})
		_, err := PrepareImport(dir, dummyBatchKey)
		expected := "date_fin_effectif is missing or invalid: "
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should return a json with one file", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})
		res, err := PrepareImport(dir, dummyBatchKey)
		//expected := FilesProperty{filter: {dummyBatchFile("filter_2002.csv")}}
		//expected := make(map[string][]string)
		expected := base.BatchFiles{base.Filter: {base.NewDummyBatchFileFromBatch("filter_2002.csv", dummyBatchKey)}}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res.Files)
		}
	})

	t.Run("Should return an _id property", func(t *testing.T) {
		batch := newSafeBatchKey("1802")
		dir := CreateTempFiles(t, batch, []string{"filter_2002.csv"})
		res, err := PrepareImport(dir, batch)
		if assert.NoError(t, err) {
			assert.Equal(t, batch, res.Key)
		}
	})

	cases := []struct {
		//id       string
		filename string
		filetype base.ParserType
	}{
		{"sigfaible_debits.csv", base.Debit},
		{"StockEtablissement_utf8_geo.csv", base.Sirene},
	}

	for _, testCase := range cases {
		t.Run("Uploaded file originally named "+testCase.filename+" should be of type "+string(testCase.filetype), func(t *testing.T) {

			dir := CreateTempFiles(t, dummyBatchKey, []string{testCase.filename, "filter_2002.csv"})

			res, err := PrepareImport(dir, dummyBatchKey)
			expected := []base.BatchFile{base.NewBatchFileFromBatch(dir, dummyBatchKey, testCase.filename)}
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res.Files[testCase.filetype])
			}
		})
	}

	t.Run("should return list of unsupported files", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"unsupported-file.csv"})
		_, err := PrepareImport(dir, dummyBatchKey)
		var e *UnsupportedFilesError
		if assert.Error(t, err) && errors.As(err, &e) {
			assert.Equal(t, []string{path.Join(dummyBatchKey.String(), "unsupported-file.csv")}, e.UnsupportedFiles)
		}
	})

	t.Run("should create filter file and fill date_fin_effectif if an effectif file is present", func(t *testing.T) {
		// setup expectations
		filterFileName := "filter_siren_" + dummyBatchKey.String() + ".csv"
		sireneULFileName := "sireneUL.csv"
		expected := base.BatchFiles{
			base.Effectif: {base.BatchFile(base.NewDummyBatchFile("sigfaible_effectif_siret.csv"))},
			base.Filter:   {base.BatchFile(base.NewDummyBatchFile(filterFileName))},
			base.SireneUl: {base.BatchFile(base.NewDummyBatchFile(sireneULFileName))},
		}
		// run prepare-import
		batchDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"sigfaible_effectif_siret.csv": ReadFileData(t, "../createfilter/test_data.csv"),
			"sireneUL.csv":                 ReadFileData(t, "../createfilter/test_uniteLegale.csv"),
		})
		adminObject, err := PrepareImport(batchDir, dummyBatchKey)
		// check that the filter is listed in the "files" property
		if assert.NoError(t, err) {
			assert.Equal(t, expected, adminObject.Files)
		}

		// check that the filter file exists
		filterFilePath := path.Join(batchDir, dummyBatchKey.String(), filterFileName)
		assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
	})

	t.Run("should create filter file even if effectif file is compressed", func(t *testing.T) {
		compressedEffectifData := compressFileData(t, "../createfilter/test_data.csv")

		// setup expectations
		filterFileName := "filter_siren_" + dummyBatchKey.String() + ".csv"

		// run prepare-import
		batchDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"sigfaible_effectif_siret.csv.gz": compressedEffectifData.Bytes(),
			"sireneUL.csv":                    ReadFileData(t, "../createfilter/test_uniteLegale.csv"),
		})

		expectedFiles := base.BatchFiles{
			base.Effectif: {base.NewBatchFileFromBatch(batchDir, dummyBatchKey, "sigfaible_effectif_siret.csv.gz")},
			base.Filter:   {base.NewBatchFileFromBatch(batchDir, dummyBatchKey, filterFileName)},
			base.SireneUl: {base.NewBatchFileFromBatch(batchDir, dummyBatchKey, "sireneUL.csv")},
		}

		adminObject, err := PrepareImport(batchDir, dummyBatchKey)
		// check that the filter is listed in the "files" property
		if assert.NoError(t, err) {
			assert.Equal(t, expectedFiles, adminObject.Files)
		}
		// check that the filter file exists
		filterFilePath := path.Join(batchDir, dummyBatchKey.String(), filterFileName)
		assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
	})
}

func makeDayDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func ReadFileData(t *testing.T, filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func compressFileData(t *testing.T, filePath string) (compressedData bytes.Buffer) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	zw := gzip.NewWriter(&compressedData)
	if _, err = zw.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return compressedData
}
