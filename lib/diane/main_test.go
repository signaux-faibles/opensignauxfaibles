package diane

import (
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestDiane(t *testing.T) {

	t.Run("Diane parser (JSON output)", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedDiane.json")
		var testData = filepath.Join("testData", "dianeTestData.txt")
		marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)
	})

	t.Run("doit détecter s'il manque des colonnes", func(t *testing.T) {
		var parser = Parser
		err := parser.initCsvReader(strings.NewReader("D1;ANNEES;ARRETE_BILAN;DENOM;CP;REGION;SECTEUR;POIDS_FRNG;TX_MARGE;DELAI_FRS;POIDS_DFISC_SOC;POIDS_FIN_CT;POIDS_FRAIS_FIN")) // typo: ANNEES au lieu de ANNEE
		if assert.Error(t, err, "initCsvReader() devrait échouer") {
			assert.Contains(t, err.Error(), "Colonne ANNEE non trouvée")
		}
	})
}
