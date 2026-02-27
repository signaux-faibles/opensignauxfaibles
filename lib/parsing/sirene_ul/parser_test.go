package sireneul

import (
	"flag"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSireneUL.json")
var testData = engine.NewBatchFile("testData", "sireneULTestData.csv")

func TestSireneUl(t *testing.T) {
	engine.TestParserOutput(t, NewSireneULParser(), testData, golden, *update)
}

func TestSireneULParser(t *testing.T) {
	testCases := []struct {
		fields   map[string]string
		expected SireneUL
	}{
		{
			fields: map[string]string{
				"siren":                                     "123456789",
				"prenom1UniteLegale":                       "Jean",
				"prenom2UniteLegale":                       "Pierre",
				"nomUniteLegale":                           "Dupont",
				"denominationUniteLegale":                  "SARL Dupont Services",
				"categorieJuridiqueUniteLegale":            "5710",
				"activitePrincipaleUniteLegale":            "47.11A",
				"nomenclatureActivitePrincipaleUniteLegale": "NAFRev2",
				"dateCreationUniteLegale":                  "2010-05-15",
				"etatAdministratifUniteLegale":             "A",
			},
			expected: SireneUL{
				Siren:               "123456789",
				RaisonSociale:       "SARL Dupont Services",
				Prenom1UniteLegale:  "Jean",
				Prenom2UniteLegale:  "Pierre",
				Prenom3UniteLegale:  "",
				Prenom4UniteLegale:  "",
				NomUniteLegale:      "Dupont",
				NomUsageUniteLegale: "",
				CategorieJuridique:  "5710",
				APE:                 "47.11A",
				Creation:            parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2010-05-15")),
				EstActif:            true,
			},
		},
		{
			fields: map[string]string{
				"siren":                                     "987654321",
				"prenom1UniteLegale":                       "Marie",
				"prenom2UniteLegale":                       "Louise",
				"prenom3UniteLegale":                       "Anne",
				"nomUniteLegale":                           "Martin",
				"nomUsageUniteLegale":                      "Durand",
				"denominationUniteLegale":                  "SAS Martin Tech",
				"categorieJuridiqueUniteLegale":            "5499",
				"activitePrincipaleUniteLegale":            "62.01Z",
				"nomenclatureActivitePrincipaleUniteLegale": "NAFRev2",
				"dateCreationUniteLegale":                  "2015-03-20",
				"etatAdministratifUniteLegale":             "F",
			},
			expected: SireneUL{
				Siren:               "987654321",
				RaisonSociale:       "SAS Martin Tech",
				Prenom1UniteLegale:  "Marie",
				Prenom2UniteLegale:  "Louise",
				Prenom3UniteLegale:  "Anne",
				Prenom4UniteLegale:  "",
				NomUniteLegale:      "Martin",
				NomUsageUniteLegale: "Durand",
				CategorieJuridique:  "5499",
				APE:                 "62.01Z",
				Creation:            parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2015-03-20")),
				EstActif:            false,
			},
		},
		{
			fields: map[string]string{
				"siren":                                     "111222333",
				"denominationUniteLegale":                  "Entreprise Collective",
				"categorieJuridiqueUniteLegale":            "9220",
				"activitePrincipaleUniteLegale":            "85.59A",
				"nomenclatureActivitePrincipaleUniteLegale": "NAFRev2",
				"dateCreationUniteLegale":                  "2018-11-10",
				"etatAdministratifUniteLegale":             "A",
			},
			expected: SireneUL{
				Siren:               "111222333",
				RaisonSociale:       "Entreprise Collective",
				Prenom1UniteLegale:  "",
				Prenom2UniteLegale:  "",
				Prenom3UniteLegale:  "",
				Prenom4UniteLegale:  "",
				NomUniteLegale:      "",
				NomUsageUniteLegale: "",
				CategorieJuridique:  "9220",
				APE:                 "85.59A",
				Creation:            parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2018-11-10")),
				EstActif:            true,
			},
		},
	}

	parser := NewSireneULParser()
	header := strings.Join(sireneULColumns, ",")
	for _, tc := range testCases {
		row, err := makeRow(tc.fields)
		require.NoError(t, err)
		instance := parser.New(strings.NewReader(header + "\n" + row))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		sireneul, ok := res.Tuples[0].(SireneUL)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, sireneul)
	}
}

func TestSireneULMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		parser := NewSireneULParser()
		instance := parser.New(parsing.CreateReader("siren", ",", []string{}))
		err := instance.Init(engine.NoFilter, nil)

		assert.Error(t, err, "should report a fatal error")
		assert.Regexp(t, regexp.MustCompile("column [^ ]+ not found"), err.Error())
	})

	t.Run("should fail if denominationUniteLegale column is missing", func(t *testing.T) {
		parser := NewSireneULParser()
		instance := parser.New(parsing.CreateReader("siren,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale", ",", []string{}))
		err := instance.Init(engine.NoFilter, nil)

		assert.Error(t, err, "should report a fatal error")
		assert.Contains(t, err.Error(), "column denominationUniteLegale not found")
	})
}

func TestSireneULNAFFiltering(t *testing.T) {
	testCases := []struct {
		name                 string
		nomenclatureActivite string
		shouldFilter         bool
	}{
		{
			"NAFRev2 should be included",
			"NAFRev2",
			false,
		},
		{
			"NAFRev1 should be excluded",
			"NAFRev1",
			true,
		},
		{
			"NAF1993 should be excluded",
			"NAF1993",
			true,
		},
		{
			"NAP should be excluded",
			"NAP",
			true,
		},
		{
			"missing nomenclature should be excluded",
			"",
			true,
		},
	}

	header := strings.Join(sireneULColumns, ",")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			row, err := makeRow(map[string]string{
				"siren":                                     "123456789",
				"etatAdministratifUniteLegale":              "A",
				"denominationUniteLegale":                   "SARL Dupont Services",
				"categorieJuridiqueUniteLegale":             "5710",
				"activitePrincipaleUniteLegale":             "47.11A",
				"nomenclatureActivitePrincipaleUniteLegale": tc.nomenclatureActivite,
			})
			require.NoError(t, err)
			parser := NewSireneULParser()
			instance := parser.New(strings.NewReader(header + "\n" + row))
			instance.Init(engine.NoFilter, nil)
			res := &engine.ParsedLineResult{}
			instance.ReadNext(res)

			if tc.shouldFilter {
				assert.NotNil(t, res.FilterError)
			} else {
				assert.Nil(t, res.FilterError)
			}
		})
	}
}

func TestSireneUlHeader(t *testing.T) {
	t.Run("can parse file that just contains a header", func(t *testing.T) {
		csvRows := []string{strings.Join(sireneULColumns, ",")}
		output := engine.RunParserInline(t, NewSireneULParser(), csvRows)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, "", engine.GetFatalError(output), "should not report a fatal error")
	})

	t.Run("reports a fatal error in case of unexpected csv header", func(t *testing.T) {
		csvRows := []string{"sirenXYZ,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale,sigleUniteLegale,sexeUniteLegale,prenom1UniteLegale,prenom2UniteLegale,prenom3UniteLegale,prenom4UniteLegale,prenomUsuelUniteLegale,pseudonymeUniteLegale,identifiantAssociationUniteLegale,trancheEffectifsUniteLegale,anneeEffectifsUniteLegale,dateDernierTraitementUniteLegale,nombrePeriodesUniteLegale,categorieEntreprise,anneeCategorieEntreprise,dateDebut,etatAdministratifUniteLegale,nomUniteLegale,nomUsageUniteLegale,denominationUniteLegale,denominationUsuelle1UniteLegale,denominationUsuelle2UniteLegale,denominationUsuelle3UniteLegale,categorieJuridiqueUniteLegale,activitePrincipaleUniteLegale,nomenclatureActivitePrincipaleUniteLegale,nicSiegeUniteLegale,economieSocialeSolidaireUniteLegale,caractereEmployeurUniteLegale"}
		output := engine.RunParserInline(t, NewSireneULParser(), csvRows)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, "Fatal: column siren not found, aborting", engine.GetFatalError(output))
	})
}
