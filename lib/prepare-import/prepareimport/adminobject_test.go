package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func dummyBatchFile(filename string) BatchFile {
	return newBatchFile(dummyBatchKey, filename)
}

func TestPopulateParamProperty(t *testing.T) {
	t.Run("Should return a date_fin consistent with batch key", func(t *testing.T) {
		res := populateParamProperty(newSafeBatchKey("1912"), validDateFinEffectif)
		expected := makeDayDate(2019, 12, 01)
		assert.Equal(t, expected, res.DateFin)
	})
}
