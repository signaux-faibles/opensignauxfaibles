package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUrssafToDate(t *testing.T) {

	t.Run("échoue si la date fournie n'est pas un nombre", func(t *testing.T) {
		_, err := UrssafToDate("11A0203")
		assert.EqualError(t, err, "valeur non autorisée pour une conversion en date: 11A0203")
	})

	t.Run("échoue si la date obtenue n'est pas valide", func(t *testing.T) {
		_, err := UrssafToDate("0000000")
		assert.EqualError(t, err, "valeur non autorisée pour une conversion en date: 0000000")
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
	testCases := []struct {
		name      string
		input     string
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{
			name:    "échoue si la date fournie n'est pas un nombre",
			input:   "",
			wantErr: true,
		},
		{
			name:    "échoue si la date ne s'étend pas sur 6 chiffre",
			input:   "2004101",
			wantErr: true,
		},
		{
			name:    "échoue si l'année n'est pas un nombre",
			input:   "AAAA1010",
			wantErr: true,
		},
		{
			name:      "0162 représente l'année 2001",
			input:     "0162",
			wantStart: makeDate(2001, 1, 1),
			wantEnd:   makeDate(2002, 1, 1),
		},
		{
			name:      "5062 représente l'année 1950",
			input:     "5062",
			wantStart: makeDate(1950, 1, 1),
			wantEnd:   makeDate(1951, 1, 1),
		},
		{
			name:      "2110 représente le 1er trimestre 2021",
			input:     "2110",
			wantStart: makeDate(2021, 1, 1),
			wantEnd:   makeDate(2021, 4, 1),
		},
		{
			name:      "2041 représente le 1er mois du 4e trimestre 2020",
			input:     "2041",
			wantStart: makeDate(2020, 10, 1),
			wantEnd:   makeDate(2020, 11, 1),
		},
		{
			name:      "1862 -> année 2018",
			input:     "1862",
			wantStart: makeDate(2018, 1, 1),
			wantEnd:   makeDate(2019, 1, 1),
		},
		{
			name:      "1820 -> 2e trimestre 2018",
			input:     "1820",
			wantStart: makeDate(2018, 4, 1),
			wantEnd:   makeDate(2018, 7, 1),
		},
		{
			name:      "6331 -> juillet 1963",
			input:     "6331",
			wantStart: makeDate(1963, 7, 1),
			wantEnd:   makeDate(1963, 8, 1),
		},
		{
			name:    "56331 -> erreur de format",
			input:   "56331",
			wantErr: true,
		},
		{
			name:    "56a1 -> erreur de format",
			input:   "56a1",
			wantErr: true,
		},
		{
			name:    "5a31 -> erreur de format",
			input:   "5a31",
			wantErr: true,
		},
		{
			name:    "564a -> erreur de format",
			input:   "564a",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotStart, gotEnd, err := UrssafToPeriod(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				assert.True(t, gotStart.IsZero())
				assert.True(t, gotEnd.IsZero())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantStart, gotStart)
				assert.Equal(t, tc.wantEnd, gotEnd)
			}
		})
	}
}

func makeDate(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
