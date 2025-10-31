package prepareimport

import (
	"bytes"
	"compress/gzip"
	"errors"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/filter"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	mockFilterWriter = &filter.MemoryFilterWriter{}

	// A filter can be read
	mockFilterReader = &filter.MemoryFilterReader{Filter: engine.NoFilter}

	// No filter can be read
	errFilterReader = &filter.MemoryFilterReader{Filter: nil}
)

const (
	filterFilename         = "filter.csv"
	effectifFilename       = "sigfaible_effectif_siret.csv"
	zippedEffectifFilename = "sigfaible_effectif_siret.csv.gz"
	debitsFilename         = "sigfaibles_debits.csv"
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
		wantedBatch := engine.NewSafeBatchKey("1803") // different of dummyBatchKey
		parentDir := CreateTempFiles(t, dummyBatchKey, []string{})
		_, err := PrepareImport(parentDir, wantedBatch, mockFilterReader, mockFilterWriter)
		expected := "could not find directory 1803 in provided path"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if 2 effectif files are provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{effectifFilename, "sigfaible_effectif_siret2.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, errFilterReader, mockFilterWriter)
		expected := "filter is missing: batch should include a filter or one effectif file"
		assert.Error(t, err)
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should identify single filter file", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})

		res, err := PrepareImport(dir, dummyBatchKey, mockFilterReader, mockFilterWriter)

		if assert.NoError(t, err) {
			assert.Contains(t, res.Files, engine.Filter)
			assert.Len(t, res.Files[engine.Filter], 1)

			filterFile := res.Files[engine.Filter][0]
			assert.NotNil(t, filterFile)
			assert.Equal(t, filterFile.Path(), engine.NewBatchFileFromBatch(dir, dummyBatchKey, "filter_2002.csv").Path())
		}
	})

	t.Run("Should include an id property", func(t *testing.T) {
		batch := engine.NewSafeBatchKey("1802")
		dir := CreateTempFiles(t, batch, []string{"filter_2002.csv"})
		res, err := PrepareImport(dir, batch, mockFilterReader, mockFilterWriter)

		if assert.NoError(t, err) {
			assert.Equal(t, batch, res.Key)
		}
	})

	t.Run("Should associate the correct type given file name", func(t *testing.T) {
		cases := []struct {
			filename string
			filetype engine.ParserType
		}{
			{debitsFilename, engine.Debit},
			{"StockEtablissement_utf8_geo.csv", engine.Sirene},
		}

		for _, testCase := range cases {

			dir := CreateTempFiles(t, dummyBatchKey, []string{testCase.filename, "filter_2002.csv"})

			res, err := PrepareImport(dir, dummyBatchKey, mockFilterReader, mockFilterWriter)

			expected := []engine.BatchFile{engine.NewBatchFileFromBatch(dir, dummyBatchKey, testCase.filename)}

			if assert.NoError(t, err) {
				assert.Equal(t, expected, res.Files[testCase.filetype])
			}
		}
	})

	t.Run("should return list of unsupported files", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"unsupported-file.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, mockFilterReader, mockFilterWriter)
		var e *UnsupportedFilesError
		if assert.Error(t, err) && errors.As(err, &e) {
			assert.Equal(t, []string{path.Join(dummyBatchKey.String(), "unsupported-file.csv")}, e.UnsupportedFiles)
		}
	})

	t.Run("should write/update filter file if an effectif file is present", func(t *testing.T) {
		// run prepare-import
		tmpDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			effectifFilename: ReadFileData(t, "../filter/testData/test_data.csv"),
			"sireneUL.csv":   ReadFileData(t, "../filter/testData/test_uniteLegale.csv"),
		})

		w := &filter.MemoryFilterWriter{}
		adminObject, err := PrepareImport(tmpDir, dummyBatchKey, mockFilterReader, w)

		if assert.NoError(t, err) {
			// Filter only appears in adminObject.Files if it has been explicitely provided
			assert.NotContains(t, adminObject.Files, engine.Filter)

			// Check that the filter data has been written
			assert.NotNil(t, w.Filter)
			assert.True(t, w.Filter.ShouldSkip("000000000"))
			assert.False(t, w.Filter.ShouldSkip("444444444"))
			assert.False(t, w.Filter.ShouldSkip("555555555"))
		}

	})

	t.Run("should create filter file even if effectif file is compressed", func(t *testing.T) {
		compressedEffectifData := compressFileData(t, "../filter/testData/test_data.csv")

		// run prepare-import
		batchDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			zippedEffectifFilename: compressedEffectifData.Bytes(),
			"sireneUL.csv":         ReadFileData(t, "../filter/testData/test_uniteLegale.csv"),
		})

		w := &filter.MemoryFilterWriter{}
		adminObject, err := PrepareImport(batchDir, dummyBatchKey, mockFilterReader, w)

		if assert.NoError(t, err) {
			// Filter only appears in adminObject.Files if it has been provided by
			// the user, not when generated from effectif
			assert.NotContains(t, adminObject.Files, engine.Filter)

			// check that the filter data has been written
			assert.NotNil(t, w.Filter)
			assert.True(t, w.Filter.ShouldSkip("000000000"))
			assert.False(t, w.Filter.ShouldSkip("444444444"))
			assert.False(t, w.Filter.ShouldSkip("555555555"))
		}
	})
}

func TestFilterErrors(t *testing.T) {

	testCases := []struct {
		name         string
		files        map[string][]byte
		filterReader engine.FilterReader
		expectError  bool
	}{
		{
			"Filtre valid explicitement fourni par l'utilisateur -> OK",
			map[string][]byte{filterFilename: nil},
			mockFilterReader, // valid filter provided, no error
			false,
		},
		{
			"Fichier effectif valide -> on crée le filtre",
			map[string][]byte{
				effectifFilename: ReadFileData(t, "../filter/testData/test_data.csv"),
			},
			errFilterReader,
			false,
		},
		{
			"Pas de fichier filtre ou effectif ou filtre en base -> échec",
			map[string][]byte{
				debitsFilename: nil,
			},
			errFilterReader,
			true,
		},
		{
			"Pas de fichier filtre ou effectif mais filtre en base -> OK",
			map[string][]byte{
				debitsFilename: nil,
			},
			mockFilterReader, // provides a valid filter
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := CreateTempFilesWithContent(t, dummyBatchKey, tc.files)

			_, err := PrepareImport(dir, dummyBatchKey, tc.filterReader, mockFilterWriter)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
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
