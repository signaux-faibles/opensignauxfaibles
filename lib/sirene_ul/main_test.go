package sireneul

import (
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSireneUL.json")
var testData = filepath.Join("testData", "sireneULTestData.csv")

func TestSireneUl(t *testing.T) {
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)
}

func TestSireneUlHeader(t *testing.T) {
	t.Run("can parse file that just contains a header", func(t *testing.T) {
		csvData := strings.Join([]string{
			"siren,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale,sigleUniteLegale,sexeUniteLegale,prenom1UniteLegale,prenom2UniteLegale,prenom3UniteLegale,prenom4UniteLegale,prenomUsuelUniteLegale,pseudonymeUniteLegale,identifiantAssociationUniteLegale,trancheEffectifsUniteLegale,anneeEffectifsUniteLegale,dateDernierTraitementUniteLegale,nombrePeriodesUniteLegale,categorieEntreprise,anneeCategorieEntreprise,dateDebut,etatAdministratifUniteLegale,nomUniteLegale,nomUsageUniteLegale,denominationUniteLegale,denominationUsuelle1UniteLegale,denominationUsuelle2UniteLegale,denominationUsuelle3UniteLegale,categorieJuridiqueUniteLegale,activitePrincipaleUniteLegale,nomenclatureActivitePrincipaleUniteLegale,nicSiegeUniteLegale,economieSocialeSolidaireUniteLegale,caractereEmployeurUniteLegale",
		}, "\n")
		csvFile := marshal.CreateTempFileWithContent(t, []byte(csvData))
		output := marshal.RunParser(Parser, marshal.NewCache(), csvFile.Name())
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, 1, len(output.Events), "should return a parsing report")
		reportData, _ := output.Events[0].ParseReport()
		assert.Equal(t, false, reportData["isFatal"], "should not report a fatal error")
	})

	t.Run("reports a fatal error in case of unexpected csv header", func(t *testing.T) {
		csvData := strings.Join([]string{
			"sirenXYZ,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale,sigleUniteLegale,sexeUniteLegale,prenom1UniteLegale,prenom2UniteLegale,prenom3UniteLegale,prenom4UniteLegale,prenomUsuelUniteLegale,pseudonymeUniteLegale,identifiantAssociationUniteLegale,trancheEffectifsUniteLegale,anneeEffectifsUniteLegale,dateDernierTraitementUniteLegale,nombrePeriodesUniteLegale,categorieEntreprise,anneeCategorieEntreprise,dateDebut,etatAdministratifUniteLegale,nomUniteLegale,nomUsageUniteLegale,denominationUniteLegale,denominationUsuelle1UniteLegale,denominationUsuelle2UniteLegale,denominationUsuelle3UniteLegale,categorieJuridiqueUniteLegale,activitePrincipaleUniteLegale,nomenclatureActivitePrincipaleUniteLegale,nicSiegeUniteLegale,economieSocialeSolidaireUniteLegale,caractereEmployeurUniteLegale",
		}, "\n")
		csvFile := marshal.CreateTempFileWithContent(t, []byte(csvData))
		output := marshal.RunParser(Parser, marshal.NewCache(), csvFile.Name())
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, 1, len(output.Events), "should return a parsing report")
		reportData, _ := output.Events[0].ParseReport()
		assert.Equal(t, true, reportData["isFatal"], "should report a fatal error")
	})
}
