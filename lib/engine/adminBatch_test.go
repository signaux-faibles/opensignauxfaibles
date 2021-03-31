package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/urssaf"
	"github.com/stretchr/testify/assert"
)

func Test_IsBatchID(t *testing.T) {
	if !base.IsBatchID("1801") {
		t.Error("1801 devrait être un ID de batch")
	}

	if base.IsBatchID("") {
		t.Error("'' ne devrait pas être considéré comme un ID de batch")
	}

	if base.IsBatchID("190193039") {
		t.Error("'190193039' ne devrait pas être considéré comme un ID de batch")
	}
	if !base.IsBatchID("1901_93039") {
		t.Error("'190193039'  devrait être considéré comme un ID de batch")
	}

	if base.IsBatchID("abcd") {
		t.Error("'abcd' ne devrait pas être considéré comme un ID de batch")
	} else {
		t.Log("'abcd' est bien rejeté: ")
	}
}

func Test_CheckBatchPaths(t *testing.T) {
	testCases := []struct {
		Filepath      string
		ErrorExpected bool
	}{
		{"./test_data/empty_file", false},
		{"./test_data/missing_file", true},
	}
	for _, tc := range testCases {
		mockbatch := base.MockBatch("debit", []string{tc.Filepath})
		err := CheckBatchPaths(&mockbatch)
		if (err == nil && tc.ErrorExpected) ||
			(err != nil && !tc.ErrorExpected) {
			// t.Log(err.Error()) // delete_me
			t.Error("Validity of path " + tc.Filepath + " is wrongly checked")
		}
	}
}

func Test_ImportBatch(t *testing.T) {
	data := make(chan *Value)
	go func() {
		for range data {
		}
	}()
	batch := base.AdminBatch{}
	err := ImportBatch(batch, []marshal.Parser{}, false, data)
	if err == nil {
		t.Error("ImportBatch devrait nous empêcher d'importer sans filtre")
	}
}

func Test_ImportBatchWithUnreadableFilter(t *testing.T) {
	data := make(chan *Value)
	go func() {
		for range data {
		}
	}()
	batch := base.MockBatch("filter", []string{"this_file_does_not_exist"})
	err := ImportBatch(batch, []marshal.Parser{}, false, data)
	if err == nil {
		t.Error("ImportBatch devrait échouer en tentant d'ouvrir un fichier filtre illisible")
	}
}

func Test_CheckBatch(t *testing.T) {

	t.Run("CheckBatch devrait réussir à parser un fichier gzip", func(t *testing.T) {
		// TODO: renommer en "CheckBatch devrait réussir à parser un fichier compressé spécifié avec le préfixe 'gzip:'"

		// Compression du fichier de données
		procolFilePath := filepath.Join("..", "urssaf", "testData", "procolTestData.csv")
		compressedFilePath := filepath.Join("..", "urssaf", "testData", "procolTestData.csv.gz")
		cmd := exec.Command("gzip", "--suffix", ".gz", "--keep", procolFilePath) // créée une version gzippée du fichier, TODO: utiliser une autre extension que .gz
		err := cmd.Run()
		assert.NoError(t, err)
		cmd.Wait()
		t.Cleanup(func() { os.Remove(compressedFilePath) })

		// Exécution de CheckBatch sur un AdminBatch mentionnant un fichier compressé
		batch := base.AdminBatch{
			Files: base.BatchFiles{
				"procol": {"../../lib/urssaf/testData/procolTestData.csv.gz"}, // TODO: utiliser extension alternative et prefixe "gzip:"
			},
		}
		InitVoidEventQueue() // permettre à CheckBatch d'envoyer des messages au canal d'événements, sans stocker dans la db
		reports, err := CheckBatch(batch, []marshal.Parser{urssaf.ParserProcol})
		if assert.NoError(t, err) {
			expectedReports := []string{"../../lib/urssaf/testData/procolTestData.csv.gz: intégration terminée, 3 lignes traitées, 0 erreurs fatales, 0 lignes rejetées, 0 lignes filtrées, 3 lignes valides"}
			assert.Equal(t, expectedReports, reports)
		}
	})
}
