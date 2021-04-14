package marshal

import (
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
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
			assert.Equal(t, makeDate(2018, 2, 3), date)
		}
	})

	t.Run("(tests récupérés depuis lib/misc/main_test.go)", func(t *testing.T) {

		a, e := UrssafToDate("1180101")
		if a == time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC) && e == nil {
			t.Log("UrssafToDate: 1180101 -> 1er janvier 2018: OK")
		} else {
			t.Error("UrssafToDate: 1180101 -> 1er janvier 2018: Fail")
		}

		a, e = UrssafToDate("11a0101")
		z := time.Time{}
		if a == z && e != nil {
			t.Log("UrssafToDate: 11a0101 -> erreur: OK")
		} else {
			t.Error("UrssafToDate: 1180101 -> erreur: Fail")
		}

		a, e = UrssafToDate("1180151")
		if a == z && e != nil {
			t.Log("UrssafToDate: 1180151 -> erreur: OK")
		} else {
			t.Error("UrssafToDate: 1180151 -> erreur: Fail")
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

	// si QM == 62 alors période annuelle sur YYYY.
	t.Run("reconnait 0162 comme représentant l'année 2001", func(t *testing.T) {
		date, err := UrssafToPeriod("0162")
		if assert.NoError(t, err) {
			assert.Equal(t, makeDate(2001, 1, 1), date.Start)
			assert.Equal(t, makeDate(2002, 1, 1), date.End)
		}
	})

	// si YY ≥ 50 alors YYYY = 19YY.
	t.Run("reconnait 5062 comme représentant l'année 1950", func(t *testing.T) {
		date, err := UrssafToPeriod("5062")
		if assert.NoError(t, err) {
			assert.Equal(t, makeDate(1950, 1, 1), date.Start)
			assert.Equal(t, makeDate(1951, 1, 1), date.End)
		}
	})

	// si M == 0 alors période trimestrielle sur le trimestre Q de YYYY.
	t.Run("reconnait 2110 comme représentant le 1er trimestre de 2021", func(t *testing.T) {
		date, err := UrssafToPeriod("2110")
		if assert.NoError(t, err) {
			assert.Equal(t, makeDate(2021, 1, 1), date.Start)
			assert.Equal(t, makeDate(2021, 4, 1), date.End)
		}
	})

	// si 0 < M < 4 alors mois M du trimestre Q.
	t.Run("reconnait 2041 comme représentant le 1er mois du 4ème trimestre de 2020", func(t *testing.T) {
		date, err := UrssafToPeriod("2041")
		if assert.NoError(t, err) {
			assert.Equal(t, makeDate(2020, 10, 1), date.Start)
			assert.Equal(t, makeDate(2020, 11, 1), date.End)
		}
	})

	t.Run("(tests récupérés depuis lib/misc/main_test.go)", func(t *testing.T) {
		a, e := UrssafToPeriod("1862")
		b := misc.Periode{
			Start: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		}
		if a == b && e == nil {
			t.Log("UrssafToPeriod: 1862 -> l'année 2018: OK")
		} else {
			t.Error("UrssafToPeriod: 1862 -> l'année 2018: Fail")
		}

		b = misc.Periode{
			Start: time.Date(2018, time.Month(4), 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2018, time.Month(7), 1, 0, 0, 0, 0, time.UTC),
		}
		a, e = UrssafToPeriod("1820")
		if a == b && e == nil {
			t.Log("UrssafToPeriod: 1820 -> 2° trimestre 2018: OK")
		} else {
			t.Error("UrssafToPeriod: 1820 -> 2° trimestre 2018: Fail")
		}

		b = misc.Periode{
			Start: time.Date(1963, time.Month(7), 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(1963, time.Month(8), 1, 0, 0, 0, 0, time.UTC),
		}
		a, e = UrssafToPeriod("6331")
		if a == b && e == nil {
			t.Log("UrssafToPeriod: 6331 -> Juillet 1963: OK")
		} else {
			t.Error("UrssafToPeriod: 6331 -> Juillet 1963: Fail")
		}

		b = misc.Periode{
			Start: time.Time{},
			End:   time.Time{},
		}
		a, e = UrssafToPeriod("56331")
		if a == b && e != nil {
			t.Log("UrssafToPeriod: 56331 -> erreur: OK")
		} else {
			t.Error("UrssafToPeriod: 56331 -> erreur: Fail")
		}

		a, e = UrssafToPeriod("56a1")
		if a == b && e != nil {
			t.Log("UrssafToPeriod: 56a1 -> erreur: OK")
		} else {
			t.Error("UrssafToPeriod: 56a1 -> erreur: Fail")
		}

		a, e = UrssafToPeriod("5a31")
		if a == b && e != nil {
			t.Log("UrssafToPeriod: 5a31 -> erreur: OK")
		} else {
			t.Error("UrssafToPeriod: 5a31 -> erreur: Fail")
		}

		a, e = UrssafToPeriod("564a")
		if a == b && e != nil {
			t.Log("UrssafToPeriod: 564a -> erreur: OK")
		} else {
			t.Error("UrssafToPeriod: 56564aa1 -> erreur: Fail")
		}
	})
}

func makeDate(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
