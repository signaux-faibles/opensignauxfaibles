package sirenehisto

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/engine"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSireneHisto.json")
var testData = engine.NewBatchFile("testData", "sireneHistoTestData.csv")

func TestSireneUl(t *testing.T) {
	engine.TestParserOutput(t, NewSireneHistoParser(), engine.NewEmptyCache(), testData, golden, *update)
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
