package prepareimport

import (
	"errors"
	"opensignauxfaibles/lib/engine"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	debitsFilename = "sigfaibles_debits.csv"
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
		_, err := PrepareImport(parentDir, wantedBatch)
		expected := "could not find directory"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expected)
	})

	t.Run("Should identify single filter file", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})

		res, err := PrepareImport(dir, dummyBatchKey)

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
		res, err := PrepareImport(dir, batch)

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

			res, err := PrepareImport(dir, dummyBatchKey)

			expected := []engine.BatchFile{engine.NewBatchFileFromBatch(dir, dummyBatchKey, testCase.filename)}

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
}

func makeDayDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
