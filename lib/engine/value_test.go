package engine

import (
	"reflect"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

type Test struct {
	value string
}

func (test Test) Key() string   { return "" }
func (test Test) Scope() string { return "" }
func (test Test) Type() string  { return "" }

func Test_mergeBatch(t *testing.T) {

	batch1 := Batch{
		"test1": map[string]marshal.Tuple{
			"hash1": Test{"test1"},
		},
	}

	batch2 := Batch{
		"test2": map[string]marshal.Tuple{
			"hash2": Test{"test2"},
		},
	}

	batch3 := Batch{
		"test1": map[string]marshal.Tuple{
			"hash2": Test{"test2"},
		},
	}

	batch4 := Batch{
		"test1": map[string]marshal.Tuple{
			"hash1": Test{"test2"},
		},
	}

	mergedBatch2 := Batch{
		"test1": {
			"hash1": Test{"test1"},
		},
		"test2": {
			"hash2": Test{"test2"},
		},
	}

	mergedBatch3 := Batch{
		"test1": {
			"hash1": Test{"test1"},
			"hash2": Test{"test2"},
		},
		"test2": {
			"hash2": Test{"test2"},
		},
	}

	mergedBatch4 := Batch{
		"test1": {
			"hash1": Test{"test2"},
			"hash2": Test{"test2"},
		},
		"test2": {
			"hash2": Test{"test2"},
		},
	}

	batch1.Merge(batch2)
	if reflect.DeepEqual(batch1, mergedBatch2) {
		t.Log("Test d'ajout d'un type: OK")
	} else {
		t.Error("Test d'ajout d'un type: Fail.")
	}

	batch1.Merge(batch3)
	if reflect.DeepEqual(batch1, mergedBatch3) {
		t.Log("Test d'ajout d'un hash: OK")
	} else {
		t.Error("Test d'ajout d'un hash: Fail.")
	}

	batch1.Merge(batch4)
	if reflect.DeepEqual(batch1, mergedBatch4) {
		t.Log("Test d'écrasement d'un hash: OK")
	} else {
		t.Error("Test d'écrasement d'un hash: Fail.")
	}

}
