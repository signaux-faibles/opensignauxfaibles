package effectif

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

func TestEffectifEnt(t *testing.T) {
	testCases := []struct {
		name     string
		csvRow   []string
		expected []EffectifEnt
	}{
		{
			name:   "one row produces multiple EffectifEnt tuples",
			csvRow: []string{"951406823818346903", "542368702", "qyt myzpmmc rkv tcpz nk fzterm", "1086e", "23", "629", "536", "209", "0", "585", "719", "733965"},
			expected: []EffectifEnt{
				{
					Siren:       "542368702",
					Periode:     parsing.MustParseTime("2006-01-02", "2010-01-01"),
					EffectifEnt: 629,
				},
				{
					Siren:       "542368702",
					Periode:     parsing.MustParseTime("2006-01-02", "2010-02-01"),
					EffectifEnt: 536,
				},
				{
					Siren:       "542368702",
					Periode:     parsing.MustParseTime("2006-01-02", "2010-03-01"),
					EffectifEnt: 209,
				},
				{
					Siren:       "542368702",
					Periode:     parsing.MustParseTime("2006-01-02", "2010-04-01"),
					EffectifEnt: 0,
				},
				{
					Siren:       "542368702",
					Periode:     parsing.MustParseTime("2006-01-02", "2010-05-01"),
					EffectifEnt: 585,
				},
			},
		},
		{
			name:   "empty values are skipped",
			csvRow: []string{"951406823818346903", "542368702", "qyt myzpmmc", "1086e", "23", "100", "", "200", "", "", "719", "733965"},
			expected: []EffectifEnt{
				{
					Siren:       "542368702",
					Periode:     parsing.MustParseTime("2006-01-02", "2010-01-01"),
					EffectifEnt: 100,
				},
				{
					Siren:       "542368702",
					Periode:     parsing.MustParseTime("2006-01-02", "2010-03-01"),
					EffectifEnt: 200,
				},
			},
		},
	}

	header := "compte;siren;rais_soc;ape_ins;dep;eff201011;eff201012;eff201013;eff201021;eff201022;base;UR_EMET"
	parser := NewEffectifEntParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instance := parser.New(parsing.CreateReader(header, ";", tc.csvRow))
			instance.Init(engine.NoFilter, nil)
			res := &engine.ParsedLineResult{}
			instance.ReadNext(res)

			assert.Empty(t, res.Errors)
			require.Equal(t, len(tc.expected), len(res.Tuples), "should produce expected number of tuples")

			for i, tuple := range res.Tuples {
				effectifEnt, ok := tuple.(EffectifEnt)
				if !ok {
					t.Fatal("tuple type is not as expected")
				}
				assert.Equal(t, tc.expected[i], effectifEnt)
			}
		})
	}
}

func TestEffectif(t *testing.T) {
	testCases := []struct {
		name     string
		csvRow   []string
		expected []Effectif
	}{
		{
			name:   "one row produces multiple Effectif tuples",
			csvRow: []string{"951406823818346903", "54236870268934", "qyt myzpmmc rkv tcpz nk fzterm", "1086e", "23", "629", "536", "209", "0", "585", "835", "733965"},
			expected: []Effectif{
				{
					Siret:        "54236870268934",
					NumeroCompte: "951406823818346903",
					Periode:      parsing.MustParseTime("2006-01-02", "2010-01-01"),
					Effectif:     629,
				},
				{
					Siret:        "54236870268934",
					NumeroCompte: "951406823818346903",
					Periode:      parsing.MustParseTime("2006-01-02", "2010-02-01"),
					Effectif:     536,
				},
				{
					Siret:        "54236870268934",
					NumeroCompte: "951406823818346903",
					Periode:      parsing.MustParseTime("2006-01-02", "2010-03-01"),
					Effectif:     209,
				},
				{
					Siret:        "54236870268934",
					NumeroCompte: "951406823818346903",
					Periode:      parsing.MustParseTime("2006-01-02", "2010-04-01"),
					Effectif:     0,
				},
				{
					Siret:        "54236870268934",
					NumeroCompte: "951406823818346903",
					Periode:      parsing.MustParseTime("2006-01-02", "2010-05-01"),
					Effectif:     585,
				},
			},
		},
		{
			name:   "empty values are skipped",
			csvRow: []string{"951406823818346903", "54236870268934", "qyt myzpmmc", "1086e", "23", "100", "", "200", "", "", "835", "733965"},
			expected: []Effectif{
				{
					Siret:        "54236870268934",
					NumeroCompte: "951406823818346903",
					Periode:      parsing.MustParseTime("2006-01-02", "2010-01-01"),
					Effectif:     100,
				},
				{
					Siret:        "54236870268934",
					NumeroCompte: "951406823818346903",
					Periode:      parsing.MustParseTime("2006-01-02", "2010-03-01"),
					Effectif:     200,
				},
			},
		},
	}

	header := "compte;siret;rais_soc;ape_ins;dep;eff201011;eff201012;eff201013;eff201021;eff201022;base;UR_EMET"
	parser := NewEffectifParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instance := parser.New(parsing.CreateReader(header, ";", tc.csvRow))
			instance.Init(engine.NoFilter, nil)
			res := &engine.ParsedLineResult{}
			instance.ReadNext(res)

			assert.Empty(t, res.Errors)
			require.Equal(t, len(tc.expected), len(res.Tuples), "should produce expected number of tuples")

			for i, tuple := range res.Tuples {
				effectif, ok := tuple.(Effectif)
				if !ok {
					t.Fatal("tuple type is not as expected")
				}
				assert.Equal(t, tc.expected[i], effectif)
			}
		})
	}
}
