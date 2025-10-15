package prepareimport

import (
	"opensignauxfaibles/lib/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateParamProperty(t *testing.T) {
	t.Run("Should return a date_fin consistent with batch key", func(t *testing.T) {
		res := populateParamProperty(engine.NewSafeBatchKey("1912"))
		expected := makeDayDate(2019, 12, 01)
		assert.Equal(t, expected, res.DateFin)
	})
}
