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
		wantedBatch := base.NewSafeBatchKey("1803") // different of dummyBatchKey
		parentDir := CreateTempFiles(t, dummyBatchKey, []string{})
		_, err := PrepareImport(parentDir, wantedBatch)
		expected := "could not find directory 1803 in provided path"
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

	t.Run("Should identify single filter file", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})

		res, err := PrepareImport(dir, dummyBatchKey)

		if assert.NoError(t, err) {
			assert.Contains(t, res.Files, base.Filter)
			assert.Len(t, res.Files[base.Filter], 1)

			filterFile := res.Files[base.Filter][0]
			assert.NotNil(t, filterFile)
			assert.Equal(t, filterFile.Path(), base.NewBatchFileFromBatch(dir, dummyBatchKey, "filter_2002.csv").Path())
		}
	})

	t.Run("Should include an id property", func(t *testing.T) {
		batch := base.NewSafeBatchKey("1802")
		dir := CreateTempFiles(t, batch, []string{"filter_2002.csv"})
		res, err := PrepareImport(dir, batch)

		if assert.NoError(t, err) {
			assert.Equal(t, batch, res.Key)
		}
	})

	t.Run("Should associate the correct type given file name", func(t *testing.T) {
		cases := []struct {
			filename string
			filetype base.ParserType
		}{
			{"sigfaible_debits.csv", base.Debit},
			{"StockEtablissement_utf8_geo.csv", base.Sirene},
		}

		for _, testCase := range cases {

			dir := CreateTempFiles(t, dummyBatchKey, []string{testCase.filename, "filter_2002.csv"})

			res, err := PrepareImport(dir, dummyBatchKey)

			expected := []base.BatchFile{base.NewBatchFileFromBatch(dir, dummyBatchKey, testCase.filename)}

			if assert.NoError(t, err) {
				assert.Equal(t, expected, res.Files[testCase.filetype])
			}
		}
	})

	t.Run("should return list of unsupported files", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"unsupported-file.csv"})
		_, err := PrepareImport(dir, dummyBatchKey)
		var e *UnsupportedFilesError
		if assert.Error(t, err) && errors.As(err, &e) {
			assert.Equal(t, []string{path.Join(dummyBatchKey.String(), "unsupported-file.csv")}, e.UnsupportedFiles)
		}
	})

	t.Run("should create filter file if an effectif file is present", func(t *testing.T) {
		// setup expectations
		filterFileName := "filter_siren.csv"

		// run prepare-import
		tmpDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"sigfaible_effectif_siret.csv": ReadFileData(t, "./createfilter/test_data.csv"),
			"sireneUL.csv":                 ReadFileData(t, "./createfilter/test_uniteLegale.csv"),
		})

		adminObject, err := PrepareImport(tmpDir, dummyBatchKey)

		// check that the filter is listed in the "files" property
		if assert.NoError(t, err) {
			assert.Contains(t, adminObject.Files, base.Filter)
			assert.Len(t, adminObject.Files[base.Filter], 1)
			assert.Equal(t, adminObject.Files[base.Filter][0].Filename(), filterFileName)

			// check that the filter file exists
			filterFilePath := adminObject.Files[base.Filter][0].Path()
			assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
		}

	})

	t.Run("should create filter file even if effectif file is compressed", func(t *testing.T) {
		compressedEffectifData := compressFileData(t, "./createfilter/test_data.csv")

		// setup expectations
		filterFileName := "filter_siren.csv"

		// run prepare-import
		batchDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"sigfaible_effectif_siret.csv.gz": compressedEffectifData.Bytes(),
			"sireneUL.csv":                    ReadFileData(t, "./createfilter/test_uniteLegale.csv"),
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
