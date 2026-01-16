package apdemande

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

func TestApdemande(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected APDemande
	}{
		{
			csvRow: []string{"f288626887", "15756082097503", "559.0", "0.0", "2012-02-11", "2017-08-09", "2017-09-09", "42078.52", "666113.56", "99.0", "1", "2", "48512.0", "319.0", "689663.3", "3"},
			expected: APDemande{
				ID:                 "f288626887",
				Siret:              "15756082097503",
				EffectifEntreprise: parsing.IntPtr(559),
				Effectif:           parsing.IntPtr(0),
				DateStatut:         parsing.MustParseTime("2006-01-02", "2012-02-11"),
				PeriodeDebut:       parsing.MustParseTime("2006-01-02", "2017-08-09"),
				PeriodeFin:         parsing.MustParseTime("2006-01-02", "2017-09-09"),
				HTA:                parsing.Float64Ptr(42078.52),
				MTA:                parsing.Float64Ptr(666113.56),
				EffectifAutorise:   parsing.IntPtr(99),
				MotifRecoursSE:     parsing.IntPtr(1),
				HeureConsommee:     parsing.Float64Ptr(48512.0),
				MontantConsomme:    parsing.Float64Ptr(689663.3),
				EffectifConsomme:   parsing.IntPtr(319),
				Perimetre:          parsing.IntPtr(2),
			},
		},
		{
			csvRow: []string{"j354624024", "50930194570891", "20.0", "82.0", "2012-03-11", "2005-08-29", "2005-09-29", "178.0", "953.96", "9.0", "3", "2", "81.0", "8", "820.94", "1"},
			expected: APDemande{
				ID:                 "j354624024",
				Siret:              "50930194570891",
				EffectifEntreprise: parsing.IntPtr(20),
				Effectif:           parsing.IntPtr(82),
				DateStatut:         parsing.MustParseTime("2006-01-02", "2012-03-11"),
				PeriodeDebut:       parsing.MustParseTime("2006-01-02", "2005-08-29"),
				PeriodeFin:         parsing.MustParseTime("2006-01-02", "2005-09-29"),
				HTA:                parsing.Float64Ptr(178.0),
				MTA:                parsing.Float64Ptr(953.96),
				EffectifAutorise:   parsing.IntPtr(9),
				MotifRecoursSE:     parsing.IntPtr(3),
				HeureConsommee:     parsing.Float64Ptr(81.0),
				MontantConsomme:    parsing.Float64Ptr(820.94),
				EffectifConsomme:   parsing.IntPtr(8),
				Perimetre:          parsing.IntPtr(2),
			},
		},
		{
			csvRow: []string{"m80v947695", "99508628226971", "949.0", "535.0", "2012-04-11", "2014-04-29", "2014-05-29", "345.0", "72697.76", "83.0", "4", "1", "6165.0", "64.0", "1456.27", "2"},
			expected: APDemande{
				ID:                 "m80v947695",
				Siret:              "99508628226971",
				EffectifEntreprise: parsing.IntPtr(949),
				Effectif:           parsing.IntPtr(535),
				DateStatut:         parsing.MustParseTime("2006-01-02", "2012-04-11"),
				PeriodeDebut:       parsing.MustParseTime("2006-01-02", "2014-04-29"),
				PeriodeFin:         parsing.MustParseTime("2006-01-02", "2014-05-29"),
				HTA:                parsing.Float64Ptr(345.0),
				MTA:                parsing.Float64Ptr(72697.76),
				EffectifAutorise:   parsing.IntPtr(83),
				MotifRecoursSE:     parsing.IntPtr(4),
				HeureConsommee:     parsing.Float64Ptr(6165.0),
				MontantConsomme:    parsing.Float64Ptr(1456.27),
				EffectifConsomme:   parsing.IntPtr(64),
				Perimetre:          parsing.IntPtr(1),
			},
		},
	}

	parser := NewApdemandeParser()
	header := "ID_DA,ETAB_SIRET,EFF_ENT,EFF_ETAB,DATE_STATUT,DATE_DEB,DATE_FIN,HTA,MTA,EFF_AUTO,MOTIF_RECOURS_SE,PERIMETRE_AP,S_HEURE_CONSOM_TOT,S_EFF_CONSOM_TOT,S_MONTANT_CONSOM_TOT,RECOURS_ANTERIEUR"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ",", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		apdemande, ok := res.Tuples[0].(APDemande)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, apdemande)
	}
}

func TestApdemandeMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		output := engine.RunParserInline(t, NewApdemandeParser(), []string{"ID_DA,ETAB_SIRET"}) // EFF_ENT is missing (among others)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "column EFF_ENT not found")
	})

	t.Run("should fail if a composite column misses", func(t *testing.T) {
		headerRow := []string{"ID_DA,ETAB_SIRET,EFF_ENT,EFF_ETAB,DATE_STATUT,HTA,EFF_AUTO,MOTIF_RECOURS_SE,S_HEURE_CONSOM_TOT,S_HEURE_CONSOM_TOT,DATE_FIN"} // DATE_DEB is missing
		output := engine.RunParserInline(t, NewApdemandeParser(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "column DATE_DEB not found")
	})
}
