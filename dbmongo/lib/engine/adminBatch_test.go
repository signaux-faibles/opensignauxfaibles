package engine

import (
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

func Test_NextBatchID(t *testing.T) {
	batchID, err := NextBatchID("abcd")

	if err == nil {
		t.Error("Erreur attendue absente. '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test valeur erronée ok")
	}
	batchID, err = NextBatchID("1801")
	if err != nil || batchID != "1802" {
		t.Error("'1802' attendu, '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test valeur courante ok")
	}
	batchID, err = NextBatchID("1812")
	if err != nil || batchID != "1901" {
		t.Error("'1901' attendu, '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test passage nouvelle année ok")
	}
}

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
	Db.ChanData = make(chan *Value)
	go func() {
		for range Db.ChanData {
		}
	}()
	batch := base.AdminBatch{}
	err := ImportBatch(batch, []marshal.Parser{}, false)
	if err == nil {
		t.Error("ImportBatch devrait nous empêcher d'importer sans filtre")
	}
}
