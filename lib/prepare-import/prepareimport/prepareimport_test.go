package prepareimport

import (
	"bytes"
	"compress/gzip"
	"errors"
	"os"
	"path"
	"path/filepath"
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
		_, err := PrepareImport(parentDir, wantedBatch, "")
		expected := "could not find directory 1803 in provided path"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if the sub-batch was not found in the specified directory", func(t *testing.T) {
		subBatch := newSafeBatchKey("1803_01")
		parentBatch := newSafeBatchKey("1803")
		parentDir := CreateTempFiles(t, parentBatch, []string{})
		_, err := PrepareImport(parentDir, subBatch, "")
		expected := "could not find directory 1803/1803_01 in provided path"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if no filter is provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"sigfaibles_debits.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		expected := "filter is missing: batch should include a filter or one effectif file"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if 2 effectif files are provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"sigfaible_effectif_siret.csv", "sigfaible_effectif_siret2.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		expected := "filter is missing: batch should include a filter or one effectif file"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if neither effectif and date_fin_effectif are provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, "")
		expected := "date_fin_effectif is missing or invalid: "
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should return a json with one file", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})
		res, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		//expected := FilesProperty{filter: {dummyBatchFile("filter_2002.csv")}}
		//expected := make(map[string][]string)
		expected := map[ValidFileType][]string{filter: {"/1802/filter_2002.csv"}}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res.Files)
		}
	})

	t.Run("Should detect and duplicate the filter file of the parent batch, given we are generating a sub-batch", func(t *testing.T) {
		// Set expectations
		subBatch := newSafeBatchKey("1803_01")
		parentBatch := newSafeBatchKey(subBatch.GetParentBatch())
		filterFile := newBatchFile(subBatch, "filter_siren_1803.csv")
		parentFilterFile := newBatchFile(parentBatch, "filter_siren_1803.csv")
		//expectedFilesProp := FilesProperty{filter: {filterFile}}
		expectedFilesProp := map[ValidFileType][]string{filter: {filterFile.Path()}}
		// Setup test environment
		parentDir := CreateTempFiles(t, parentBatch, []string{parentFilterFile.Name()})
		subBatchDir := filepath.Join(parentDir, parentBatch.String(), subBatch.String())
		_ = os.Mkdir(subBatchDir, 0777)
		// Run the test
		res, err := PrepareImport(parentDir, subBatch, "2018-03-01")
		if assert.NoError(t, err) {
			assert.Equal(t, expectedFilesProp, res.Files)
			duplicatedFilePath := path.Join(parentDir, parentBatch.GetParentBatch(), filterFile.Path())
			assert.True(t, fileExists(duplicatedFilePath))
		}
	})

	t.Run("Should infer date_fin_effectif from the effectif file, given no filter was found", func(t *testing.T) {
		// Set expectations
		effectifFile := dummyBatchFile("sigfaible_effectif_siret.csv")
		sireneULFile := dummyBatchFile("sireneUL.csv")
		providedDateFinEffetif := ""
		expectedDateFinEffectif := NewDateFinEffectif(time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)).Date()
		// Setup test environment
		parentDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			effectifFile.Name(): ReadFileData(t, "../createfilter/test_data.csv"),
			sireneULFile.Name(): ReadFileData(t, "../createfilter/test_uniteLegale.csv"),
		})
		// Run the test
		res, err := PrepareImport(parentDir, dummyBatchKey, providedDateFinEffetif)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedDateFinEffectif, res.Param.DateFinEffectif)
		}
	})

	t.Run("Should infer date_fin_effectif from the effectif file, given a filter was found", func(t *testing.T) {
		// Set expectations
		effectifFile := dummyBatchFile("sigfaible_effectif_siret.csv")
		providedDateFinEffetif := ""
		expectedDateFinEffectif := NewDateFinEffectif(time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)).Date()
		// Setup test environment
		parentDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			effectifFile.Name(): ReadFileData(t, "../createfilter/test_data.csv"),
			"filter_2002.csv":   {},
		})
		// Run the test
		res, err := PrepareImport(parentDir, dummyBatchKey, providedDateFinEffetif)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedDateFinEffectif, res.Param.DateFinEffectif)
		}
	})

	t.Run("Should return an _id property", func(t *testing.T) {
		batch := newSafeBatchKey("1802")
		dir := CreateTempFiles(t, batch, []string{"filter_2002.csv"})
		res, err := PrepareImport(dir, batch, dummyDateFinEffectif)
		if assert.NoError(t, err) {
			assert.Equal(t, IDProperty{batch, "batch"}, res.ID)
		}
	})

	cases := []struct {
		//id       string
		filename string
		goupPath string
		filetype ValidFileType
	}{
		{"sigfaible_debits.csv", "urssaf", debit},
		{"StockEtablissement_utf8_geo.csv", "random", sirene},
	}

	for _, testCase := range cases {
		t.Run("Uploaded file originally named "+testCase.filename+" should be of type "+string(testCase.filetype), func(t *testing.T) {

			dir := CreateTempFiles(t, dummyBatchKey, []string{testCase.filename, "filter_2002.csv"})

			res, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
			expected := []string{dummyBatchFile(testCase.filename).Path()}
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res.Files[testCase.filetype])
			}
		})
	}

	t.Run("should return list of unsupported files", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"unsupported-file.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		var e *UnsupportedFilesError
		if assert.Error(t, err) && errors.As(err, &e) {
			assert.Equal(t, []string{dummyBatchKey.Path() + "unsupported-file.csv"}, e.UnsupportedFiles)
		}
	})

	t.Run("should create filter file and fill date_fin_effectif if an effectif file is present", func(t *testing.T) {
		// setup expectations
		filterFileName := "filter_siren_" + dummyBatchKey.String() + ".csv"
		sireneULFileName := "sireneUL.csv"
		expected := map[ValidFileType][]string{
			"effectif":  {dummyBatchFile("sigfaible_effectif_siret.csv").Path()},
			"filter":    {dummyBatchFile(filterFileName).Path()},
			"sirene_ul": {dummyBatchFile(sireneULFileName).Path()},
		}
		expectedDateFinEffectif := makeDayDate(2020, 1, 1)
		// run prepare-import
		batchDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"sigfaible_effectif_siret.csv": ReadFileData(t, "../createfilter/test_data.csv"),
			"sireneUL.csv":                 ReadFileData(t, "../createfilter/test_uniteLegale.csv"),
		})
		adminObject, err := PrepareImport(batchDir, dummyBatchKey, "")
		// check that the filter is listed in the "files" property
		if assert.NoError(t, err) {
			assert.Equal(t, expected, adminObject.Files)
		}
		// check that the filter file exists
		filterFilePath := path.Join(batchDir, dummyBatchKey.Path(), filterFileName)
		assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
		// check that date_fin_effectif was detected from the effectif file
		actualDateFinEffectif := adminObject.Param.DateFinEffectif
		assert.Equal(t, expectedDateFinEffectif, actualDateFinEffectif)
	})

	t.Run("should create filter file even if effectif file is compressed", func(t *testing.T) {
		compressedEffectifData := compressFileData(t, "../createfilter/test_data.csv")
		// setup expectations
		filterFileName := "filter_siren_" + dummyBatchKey.String() + ".csv"
		expectedEffectifFile := &batchFile{
			batchKey:    dummyBatchKey,
			filename:    "sigfaible_effectif_siret.csv.gz",
			gzippedSize: uint64(compressedEffectifData.Len()),
		}
		expectedFiles := map[ValidFileType][]string{
			"effectif":  {expectedEffectifFile.Path()},
			"filter":    {dummyBatchFile(filterFileName).Path()},
			"sirene_ul": {dummyBatchFile("sireneUL.csv").Path()},
		}
		expectedDateFinEffectif := makeDayDate(2020, 1, 1)
		// run prepare-import
		batchDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"sigfaible_effectif_siret.csv.gz": compressedEffectifData.Bytes(),
			"sireneUL.csv":                    ReadFileData(t, "../createfilter/test_uniteLegale.csv"),
		})
		adminObject, err := PrepareImport(batchDir, dummyBatchKey, "")
		// check that the filter is listed in the "files" property
		if assert.NoError(t, err) {
			assert.Equal(t, expectedFiles, adminObject.Files)
		}
		// check that the filter file exists
		filterFilePath := path.Join(batchDir, dummyBatchKey.Path(), filterFileName)
		assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
		// check that date_fin_effectif was detected from the effectif file
		actualDateFinEffectif := adminObject.Param.DateFinEffectif
		assert.Equal(t, expectedDateFinEffectif, actualDateFinEffectif)
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
