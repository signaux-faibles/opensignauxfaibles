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

	t.Run("Should return a date_debut 10 years before date_fin", func(t *testing.T) {
		res := populateParamProperty(engine.NewSafeBatchKey("2603"))
		expectedDateFin := makeDayDate(2026, 3, 1)
		expectedDateDebut := makeDayDate(2016, 3, 1)
		assert.Equal(t, expectedDateFin, res.DateFin, "DateFin should match batch key")
		assert.Equal(t, expectedDateDebut, res.DateDebut, "DateDebut should be 10 years before DateFin")
	})

	t.Run("Should maintain 10 year window for different batch keys", func(t *testing.T) {
		testCases := []struct {
			batchKey         string
			expectedDateFin  string
			expectedDateDebut string
		}{
			{"2601", "2026-01-01", "2016-01-01"},
			{"2512", "2025-12-01", "2015-12-01"},
			{"2406", "2024-06-01", "2014-06-01"},
		}

		for _, tc := range testCases {
			res := populateParamProperty(engine.NewSafeBatchKey(tc.batchKey))
			assert.Equal(t, tc.expectedDateFin, res.DateFin.Format("2006-01-02"), "DateFin for batch %s", tc.batchKey)
			assert.Equal(t, tc.expectedDateDebut, res.DateDebut.Format("2006-01-02"), "DateDebut for batch %s", tc.batchKey)
		}
	})
}
