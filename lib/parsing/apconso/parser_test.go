package apconso

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

func TestApconso(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected APConso
	}{
		{
			csvRow: []string{"d363954866", "61433760225578", "311.54", "214.0", "0.0", "2004-12-01"},
			expected: APConso{
				ID:             "d363954866",
				Siret:          "61433760225578",
				HeureConsommee: parsing.Float64Ptr(311.54),
				Montant:        parsing.Float64Ptr(214),
				Effectif:       parsing.IntPtr(0),
				Periode:        parsing.MustParseTime("2006-01-02", "2004-12-01"),
			},
		},
		{
			csvRow: []string{"u90q913412", "65272073959845", "44.11", "659.0", "2.0", "2019-01-01"},
			expected: APConso{
				ID:             "u90q913412",
				Siret:          "65272073959845",
				HeureConsommee: parsing.Float64Ptr(44.11),
				Montant:        parsing.Float64Ptr(659),
				Effectif:       parsing.IntPtr(2),
				Periode:        parsing.MustParseTime("2006-01-02", "2019-01-01"),
			},
		},
		{
			csvRow: []string{"p95a076375", "53312397597505", "780.0", "62.0", "16.0", "2020-12-30"},
			expected: APConso{
				ID:             "p95a076375",
				Siret:          "53312397597505",
				HeureConsommee: parsing.Float64Ptr(780),
				Montant:        parsing.Float64Ptr(62),
				Effectif:       parsing.IntPtr(16),
				Periode:        parsing.MustParseTime("2006-01-02", "2020-12-30"),
			},
		},
	}

	parser := NewApconsoParser()
	header := "ID_DA,ETAB_SIRET,HEURES,MONTANTS,EFFECTIFS,MOIS"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ",", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		apconso, ok := res.Tuples[0].(APConso)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, apconso, tc.expected)
	}
}
