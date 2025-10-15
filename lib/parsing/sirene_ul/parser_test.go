package sireneul

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/engine"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSireneUL.json")
var testData = engine.NewBatchFile("testData", "sireneULTestData.csv")

func TestSireneUl(t *testing.T) {
	engine.TestParserOutput(t, NewSireneULParser(), engine.NewEmptyCache(), testData, golden, *update)
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
		assert.Equal(t, "Fatal: Colonne siren non trouvée. Abandon.", engine.GetFatalError(output))
	})
}
