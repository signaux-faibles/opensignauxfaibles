package sireneul

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
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
		csvFile := createTempFileWithContent(t, []byte(csvData))
		output := marshal.RunParser(Parser, marshal.NewCache(), csvFile.Name())
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Equal(t, 1, len(output.Events), "should return a parsing report")
		reportData, _ := parseReport(output.Events[0])
		assert.Equal(t, false, reportData["isFatal"], "should not report a fatal error")
	})
}

func parseReport(reportEvent marshal.Event) (map[string]interface{}, error) {
	var jsonDocument map[string]interface{}
	temporaryBytes, err := bson.MarshalJSON(reportEvent.Comment)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(temporaryBytes, &jsonDocument)
	return jsonDocument, err
}

func createTempFileWithContent(t *testing.T, content []byte) *os.File {
	t.Helper()
	tmpfile, err := ioutil.TempFile("", "createTempFileWithContent")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile
}
