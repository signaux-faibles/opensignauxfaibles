package apconso

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApconso(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApconso.json")
	var testData = filepath.Join("testData", "apconsoTestData.csv")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)

	t.Run("doit détecter s'il manque des colonnes", func(t *testing.T) {
		output := marshal.RunParserInline(t, Parser, []string{"ID_DA,ETAB_SIRET,MOIS,HEURE,MONTANTS,EFFECTIFS"}) // typo: HEURE au lieu de HEURES
		assert.Equal(t, []marshal.Tuple(nil), output.Tuples, "le parseur doit retourner aucun tuple")
		assert.Contains(t, marshal.GetFatalError(output), "entête non conforme")
	})
}
