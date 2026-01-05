package sirenehisto

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

var golden = filepath.Join("testData", "expectedSireneHisto.json")
var testData = engine.NewBatchFile("testData", "sireneHistoTestData.csv")

func TestSireneUl(t *testing.T) {
	engine.TestParserOutput(t, NewSireneHistoParser(), testData, golden, *update)
}

func TestSireneHistoParser(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected SireneHisto
	}{
		{
			csvRow: []string{
				"123456789", "00012", "12345678900012", "2020-12-31", "2010-05-15", "A", "true",
				"", "", "", "", "", "", "", "", "", "", "",
			},
			expected: SireneHisto{
				Siret:                 "12345678900012",
				DateDebut:             parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2010-05-15")),
				DateFin:               parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2020-12-31")),
				EstActif:              true,
				ChangementStatutActif: true,
			},
		},
		{
			csvRow: []string{
				"987654321", "00025", "98765432100025", "", "2015-03-20", "F", "false",
				"", "", "", "", "", "", "", "", "", "", "",
			},
			expected: SireneHisto{
				Siret:                 "98765432100025",
				DateDebut:             parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2015-03-20")),
				DateFin:               nil,
				EstActif:              false,
				ChangementStatutActif: false,
			},
		},
		{
			csvRow: []string{
				"111222333", "00030", "11122233300030", "2019-06-30", "", "A", "true",
				"", "", "", "", "", "", "", "", "", "", "",
			},
			expected: SireneHisto{
				Siret:                 "11122233300030",
				DateDebut:             nil,
				DateFin:               parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2019-06-30")),
				EstActif:              true,
				ChangementStatutActif: true,
			},
		},
	}

	parser := NewSireneHistoParser()
	header := "siren,nic,siret,dateFin,dateDebut,etatAdministratifEtablissement,changementEtatAdministratifEtablissement,enseigne1Etablissement,enseigne2Etablissement,enseigne3Etablissement,changementEnseigneEtablissement,denominationUsuelleEtablissement,changementDenominationUsuelleEtablissement,activitePrincipaleEtablissement,nomenclatureActivitePrincipaleEtablissement,changementActivitePrincipaleEtablissement,caractereEmployeurEtablissement,changementCaractereEmployeurEtablissement"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ",", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		sireneHisto, ok := res.Tuples[0].(SireneHisto)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, sireneHisto)
	}
}

func TestSireneHistoMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		output := engine.RunParserInline(t, NewSireneHistoParser(), []string{"siren,nic"})
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, regexp.MustCompile("column [^ ]+ not found"), engine.GetFatalError(output))
	})

	t.Run("should fail if siret column is missing", func(t *testing.T) {
		headerRow := []string{"siren,nic,dateFin,dateDebut,etatAdministratifEtablissement"}
		output := engine.RunParserInline(t, NewSireneHistoParser(), headerRow)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Contains(t, engine.GetFatalError(output), "column siret not found")
	})
}

func TestSireneHistoFilterCases(t *testing.T) {
	parser := NewSireneHistoParser()
	header := "siren,nic,siret,dateFin,dateDebut,etatAdministratifEtablissement,changementEtatAdministratifEtablissement,enseigne1Etablissement,enseigne2Etablissement,enseigne3Etablissement,changementEnseigneEtablissement,denominationUsuelleEtablissement,changementDenominationUsuelleEtablissement,activitePrincipaleEtablissement,nomenclatureActivitePrincipaleEtablissement,changementActivitePrincipaleEtablissement,caractereEmployeurEtablissement,changementCaractereEmployeurEtablissement"

	t.Run("should filter when both dates are missing", func(t *testing.T) {
		csvRow := []string{
			"123456789", "00012", "12345678900012", "", "", "A", "true",
			"", "", "", "", "", "", "", "", "", "", "",
		}
		instance := parser.New(parsing.CreateReader(header, ",", csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.NotNil(t, res.FilterError)
		assert.Contains(t, res.FilterError.Error(), "both dateDebut and dateFin are missing")
		assert.Empty(t, res.Tuples)
	})

	t.Run("should filter when etatAdministratif is invalid", func(t *testing.T) {
		csvRow := []string{
			"123456789", "00012", "12345678900012", "2020-12-31", "2010-05-15", "X", "true",
			"", "", "", "", "", "", "", "", "", "", "",
		}
		instance := parser.New(parsing.CreateReader(header, ",", csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.NotNil(t, res.FilterError)
		assert.Contains(t, res.FilterError.Error(), "état administratif malformé")
		assert.Empty(t, res.Tuples)
	})

	t.Run("should filter when etatAdministratif is empty (null)", func(t *testing.T) {
		csvRow := []string{
			"123456789", "00012", "12345678900012", "2020-12-31", "2010-05-15", "", "true",
			"", "", "", "", "", "", "", "", "", "", "",
		}
		instance := parser.New(parsing.CreateReader(header, ",", csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.NotNil(t, res.FilterError)
		assert.Contains(t, res.FilterError.Error(), "état administratif malformé")
		assert.Empty(t, res.Tuples)
	})
}

func TestSireneUlHeader(t *testing.T) {
	t.Run("can parse file that just contains a header", func(t *testing.T) {
		csvRows := []string{"siren,nic,siret,dateFin,dateDebut,etatAdministratifEtablissement,changementEtatAdministratifEtablissement,enseigne1Etablissement,enseigne2Etablissement,enseigne3Etablissement,changementEnseigneEtablissement,denominationUsuelleEtablissement,changementDenominationUsuelleEtablissement,activitePrincipaleEtablissement,nomenclatureActivitePrincipaleEtablissement,changementActivitePrincipaleEtablissement,caractereEmployeurEtablissement,changementCaractereEmployeurEtablissement"}
		output := engine.RunParserInline(t, NewSireneHistoParser(), csvRows)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, "", engine.GetFatalError(output), "should not report a fatal error")
	})

	t.Run("reports a fatal error in case of unexpected csv header", func(t *testing.T) {
		csvRows := []string{"siren,nic,siretXYZ,dateFin,dateDebut,etatAdministratifEtablissement,changementEtatAdministratifEtablissement,enseigne1Etablissement,enseigne2Etablissement,enseigne3Etablissement,changementEnseigneEtablissement,denominationUsuelleEtablissement,changementDenominationUsuelleEtablissement,activitePrincipaleEtablissement,nomenclatureActivitePrincipaleEtablissement,changementActivitePrincipaleEtablissement,caractereEmployeurEtablissement,changementCaractereEmployeurEtablissement"}
		output := engine.RunParserInline(t, NewSireneHistoParser(), csvRows)
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, "Fatal: column siret not found, aborting", engine.GetFatalError(output))
	})
}
