package engine

import (
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
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
