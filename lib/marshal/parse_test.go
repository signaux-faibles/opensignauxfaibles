package marshal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUrssafToDate(t *testing.T) {

	t.Run("échoue si la date fournie n'est pas un nombre", func(t *testing.T) {
		_, err := UrssafToDate("11A0203")
		assert.EqualError(t, err, "Valeur non autorisée pour une conversion en date: 11A0203")
	})

	t.Run("échoue si la date obtenue n'est pas valide", func(t *testing.T) {
		_, err := UrssafToDate("0000000")
		assert.EqualError(t, err, "Valeur non autorisée pour une conversion en date: 0000000")
	})

	t.Run("reconnait 1180203 comme représentant le 3 février 2018", func(t *testing.T) {
		date, err := UrssafToDate("1180203")
		if assert.NoError(t, err) {
			assert.Equal(t, 2018, date.Year())
			assert.Equal(t, time.Month(2), date.Month())
			assert.Equal(t, 3, date.Day())
		}
	})
}

func TestUrssafToPeriod(t *testing.T) {

	t.Run("échoue si la date fournie n'est pas un nombre", func(t *testing.T) {
		_, err := UrssafToPeriod("")
		assert.EqualError(t, err, "Valeur non autorisée")
	})

	t.Run("échoue si la date ne s'étend pas sur 6 chiffre", func(t *testing.T) {
		_, err := UrssafToPeriod("2004101")
		assert.EqualError(t, err, "Valeur non autorisée")
	})

	t.Run("échoue si l'année n'est pas un nombre", func(t *testing.T) {
		_, err := UrssafToPeriod("AAAA1010")
		assert.EqualError(t, err, "Valeur non autorisée")
	})

	t.Run("reconnait 0162 comme représentant l'année 2001", func(t *testing.T) {
		date, err := UrssafToPeriod("0162")
		if assert.NoError(t, err) {
			assert.Equal(t, 2001, date.Start.Year())
			assert.Equal(t, time.Month(1), date.Start.Month())
			assert.Equal(t, 1, date.Start.Day())
			assert.Equal(t, 2002, date.End.Year())
			assert.Equal(t, time.Month(1), date.End.Month())
			assert.Equal(t, 1, date.End.Day())
		}
	})

	// si QM == 62 alors période annuelle sur YYYY.
	t.Run("reconnait 0162 comme représentant l'année 2001", func(t *testing.T) {
		date, err := UrssafToPeriod("0162")
		if assert.NoError(t, err) {
			assert.Equal(t, 2001, date.Start.Year())
			assert.Equal(t, time.Month(1), date.Start.Month())
			assert.Equal(t, 1, date.Start.Day())
			assert.Equal(t, 2002, date.End.Year())
			assert.Equal(t, time.Month(1), date.End.Month())
			assert.Equal(t, 1, date.End.Day())
		}
	})

	// si YY ≥ 50 alors YYYY = 19YY.
	t.Run("reconnait 5062 comme représentant l'année 1950", func(t *testing.T) {
		date, err := UrssafToPeriod("5062")
		if assert.NoError(t, err) {
			assert.Equal(t, 1950, date.Start.Year())
			assert.Equal(t, time.Month(1), date.Start.Month())
			assert.Equal(t, 1, date.Start.Day())
			assert.Equal(t, 1951, date.End.Year())
			assert.Equal(t, time.Month(1), date.End.Month())
			assert.Equal(t, 1, date.End.Day())
		}
	})

	// si M == 0 alors période trimestrielle sur le trimestre Q de YYYY.
	t.Run("reconnait 2110 comme représentant le 1er trimestre de 2021", func(t *testing.T) {
		date, err := UrssafToPeriod("2110")
		if assert.NoError(t, err) {
			assert.Equal(t, 2021, date.Start.Year())
			assert.Equal(t, time.Month(1), date.Start.Month())
			assert.Equal(t, 1, date.Start.Day())
			assert.Equal(t, 2021, date.End.Year())
			assert.Equal(t, time.Month(4), date.End.Month())
			assert.Equal(t, 1, date.End.Day())
		}
	})

	// si 0 < M < 4 alors mois M du trimestre Q.
	t.Run("reconnait 2041 comme représentant le 1er mois du 4ème trimestre de 2020", func(t *testing.T) {
		date, err := UrssafToPeriod("2041")
		if assert.NoError(t, err) {
			assert.Equal(t, 2020, date.Start.Year())
			assert.Equal(t, time.Month(10), date.Start.Month())
			assert.Equal(t, 1, date.Start.Day())
			assert.Equal(t, 2020, date.End.Year())
			assert.Equal(t, time.Month(11), date.End.Month())
			assert.Equal(t, 1, date.End.Day())
		}
	})
}
