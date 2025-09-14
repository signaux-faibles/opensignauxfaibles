package prepareimport

import (
	"opensignauxfaibles/lib/base"
	"testing"

	"github.com/stretchr/testify/assert"
)

func dummyBatchFile(filename string) BatchFile {
	return newBatchFile(dummyBatchKey, filename)
}

func TestPopulateCompleteTypesProperty(t *testing.T) {
	t.Run("Should not return a debit file as a complete_type, by default", func(t *testing.T) {
		res := populateCompleteTypesProperty(FilesProperty{"debit": {dummyBatchFile("sigfaibles_debits.csv")}})
		expected := []base.ValidFileType{}
		assert.Equal(t, expected, res)
	})

	t.Run("Should not return a small debit file as a complete_type", func(t *testing.T) {
		expected := []base.ValidFileType{}
		debitBatchFile := batchFile{
			batchKey:    dummyBatchKey,
			filename:    "sigfaibles_debits.csv",
			gzippedSize: thresholdPerGzippedFileType["debit"] - 1, // just below the threshold
		}
		res := populateCompleteTypesProperty(FilesProperty{"debit": {&debitBatchFile}})
		assert.Equal(t, expected, res)
	})

	t.Run("Should return a large gzipped debit file as a complete_type", func(t *testing.T) {
		expected := []base.ValidFileType{base.Debit}
		debitBatchFile := batchFile{
			batchKey:    dummyBatchKey,
			filename:    "sigfaibles_debits.csv",
			gzippedSize: 254781489, // thresholdPerGzippedFileType["debit"]
		}
		res := populateCompleteTypesProperty(FilesProperty{"debit": {&debitBatchFile}})
		assert.Equal(t, expected, res)
	})

	t.Run("Should return apconso as a complete_type", func(t *testing.T) {
		res := populateCompleteTypesProperty(FilesProperty{"apconso": {dummyBatchFile("act_partielle_conso_depuis2014_FRANCE.csv")}})
		expected := []base.ValidFileType{base.Apconso}
		assert.Equal(t, expected, res)
	})

}

func TestPopulateParamProperty(t *testing.T) {
	t.Run("Should return a date_fin consistent with batch key", func(t *testing.T) {
		res := populateParamProperty(newSafeBatchKey("1912"), validDateFinEffectif)
		expected := makeDayDate(2019, 12, 01)
		assert.Equal(t, expected, res.DateFin)
	})
}
