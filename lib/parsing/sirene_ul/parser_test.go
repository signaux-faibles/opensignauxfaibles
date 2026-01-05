package sireneul

import (
	"flag"
	"path/filepath"
	"regexp"
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
		csvRow   []string
		expected SireneUL
	}{
		{
			csvRow: []string{
				"123456789", "O", "false", "2010-05-15", "ABC", "M", "Jean", "Pierre", "", "",
				"Jean", "", "W123456789", "12", "2023", "2024-01-15T10:30:00", "3",
				"PME", "2023", "2010-05-15", "A", "Dupont", "", "SARL Dupont Services", "",
				"", "", "5710", "47.11A", "NAFRev2", "00012", "O", "O",
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
				CodeStatutJuridique: "5710",
				APE:                 "47.11A",
				Creation:            parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2010-05-15")),
				EstActif:            true,
			},
		},
		{
			csvRow: []string{
				"987654321", "O", "false", "2015-03-20", "XYZ", "F", "Marie", "Louise", "Anne", "",
				"Marie", "", "", "22", "2023", "2024-01-15T10:30:00", "2",
				"ETI", "2023", "2015-03-20", "F", "Martin", "Durand", "SAS Martin Tech", "",
				"", "", "5499", "62.01Z", "NAFRev2", "00025", "N", "O",
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
				CodeStatutJuridique: "5499",
				APE:                 "62.01Z",
				Creation:            parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2015-03-20")),
				EstActif:            false,
			},
		},
		{
			csvRow: []string{
				"111222333", "O", "false", "2018-11-10", "", "", "", "", "", "",
				"", "", "", "03", "2023", "2024-01-15T10:30:00", "1",
				"GE", "2023", "2018-11-10", "A", "", "", "Entreprise Collective", "",
				"", "", "9220", "85.59A", "NAFRev2", "00030", "O", "N",
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
				CodeStatutJuridique: "9220",
				APE:                 "85.59A",
				Creation:            parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2018-11-10")),
				EstActif:            true,
			},
		},
	}

	parser := NewSireneULParser()
	header := "siren,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale,sigleUniteLegale,sexeUniteLegale,prenom1UniteLegale,prenom2UniteLegale,prenom3UniteLegale,prenom4UniteLegale,prenomUsuelUniteLegale,pseudonymeUniteLegale,identifiantAssociationUniteLegale,trancheEffectifsUniteLegale,anneeEffectifsUniteLegale,dateDernierTraitementUniteLegale,nombrePeriodesUniteLegale,categorieEntreprise,anneeCategorieEntreprise,dateDebut,etatAdministratifUniteLegale,nomUniteLegale,nomUsageUniteLegale,denominationUniteLegale,denominationUsuelle1UniteLegale,denominationUsuelle2UniteLegale,denominationUsuelle3UniteLegale,categorieJuridiqueUniteLegale,activitePrincipaleUniteLegale,nomenclatureActivitePrincipaleUniteLegale,nicSiegeUniteLegale,economieSocialeSolidaireUniteLegale,caractereEmployeurUniteLegale"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ",", tc.csvRow))
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
		output := engine.RunParserInline(t, NewSireneULParser(), []string{"siren"}) // many columns are missing
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, regexp.MustCompile("column [^ ]+ not found"), engine.GetFatalError(output))
	})

	t.Run("should fail if denominationUniteLegale column is missing", func(t *testing.T) {
		headerRow := []string{"siren,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale"} // denominationUniteLegale is missing
		output := engine.RunParserInline(t, NewSireneULParser(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "column denominationUniteLegale not found")
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

	header := "siren,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale,sigleUniteLegale,sexeUniteLegale,prenom1UniteLegale,prenom2UniteLegale,prenom3UniteLegale,prenom4UniteLegale,prenomUsuelUniteLegale,pseudonymeUniteLegale,identifiantAssociationUniteLegale,trancheEffectifsUniteLegale,anneeEffectifsUniteLegale,dateDernierTraitementUniteLegale,nombrePeriodesUniteLegale,categorieEntreprise,anneeCategorieEntreprise,dateDebut,etatAdministratifUniteLegale,nomUniteLegale,nomUsageUniteLegale,denominationUniteLegale,denominationUsuelle1UniteLegale,denominationUsuelle2UniteLegale,denominationUsuelle3UniteLegale,categorieJuridiqueUniteLegale,activitePrincipaleUniteLegale,nomenclatureActivitePrincipaleUniteLegale,nicSiegeUniteLegale,economieSocialeSolidaireUniteLegale,caractereEmployeurUniteLegale"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			csvRow := []string{
				"123456789", "O", "false", "2010-05-15", "ABC", "M", "Jean", "Pierre", "", "",
				"Jean", "", "W123456789", "12", "2023", "2024-01-15T10:30:00", "3",
				"PME", "2023", "2010-05-15", "A", "Dupont", "", "SARL Dupont Services", "",
				"", "", "5710", "47.11A", tc.nomenclatureActivite, "00012", "O", "O",
			}

			parser := NewSireneULParser()
			instance := parser.New(parsing.CreateReader(header, ",", csvRow))
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
		csvRows := []string{"siren,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale,sigleUniteLegale,sexeUniteLegale,prenom1UniteLegale,prenom2UniteLegale,prenom3UniteLegale,prenom4UniteLegale,prenomUsuelUniteLegale,pseudonymeUniteLegale,identifiantAssociationUniteLegale,trancheEffectifsUniteLegale,anneeEffectifsUniteLegale,dateDernierTraitementUniteLegale,nombrePeriodesUniteLegale,categorieEntreprise,anneeCategorieEntreprise,dateDebut,etatAdministratifUniteLegale,nomUniteLegale,nomUsageUniteLegale,denominationUniteLegale,denominationUsuelle1UniteLegale,denominationUsuelle2UniteLegale,denominationUsuelle3UniteLegale,categorieJuridiqueUniteLegale,activitePrincipaleUniteLegale,nomenclatureActivitePrincipaleUniteLegale,nicSiegeUniteLegale,economieSocialeSolidaireUniteLegale,caractereEmployeurUniteLegale"}
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
