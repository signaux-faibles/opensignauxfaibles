package engine

import (
	"flag"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
)

var _ = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

func TestToQueries(t *testing.T) {
	var chunks = Chunks{
		OK: 1,
		SplitKeys: []splitKey{
			{"01234567800001"},
			{"01234567800002"},
			{"11234567800001"},
		},
	}
	var res = chunks.ToQueries(bson.M{}, "_id")
	if len(res) != 3 {
		t.Errorf("ToQueries devrait produire 3 requêtes, et non %d", len(res))
	}
}

func TestMakeMapReduceJobFromJsFunctions(t *testing.T) {
	t.Run("doit encapsuler les fonctions et paramètres transmis", func(t *testing.T) {
		jsFunctions := map[string]string{
			"map":           "function map(){ /* javascript code */ }",
			"reduce":        "function reduce(){ /* javascript code */ }",
			"finalize":      "function finalize(){ /* javascript code */ }",
			"otherFunction": "function otherFunction(){ /* javascript code */ }",
		}
		jsParams := bson.M{
			"globalParam": "someValue",
		}
		expectedFunctions := map[string]bson.JavaScript{
			"map":           {Code: jsFunctions["map"]},
			"reduce":        {Code: jsFunctions["reduce"]},
			"finalize":      {Code: jsFunctions["finalize"]},
			"otherFunction": {Code: jsFunctions["otherFunction"]},
		}
		expectedJob := mgo.MapReduce{
			Map:      jsFunctions["map"],
			Reduce:   jsFunctions["reduce"],
			Finalize: jsFunctions["finalize"],
			Scope: bson.M{
				"f":           expectedFunctions,
				"globalParam": jsParams["globalParam"],
			},
		}
		mapReduceJob, err := makeMapReduceJobFromJsFunctions(jsFunctions, jsParams)
		if assert.NoError(t, err) {
			assert.EqualValues(t, expectedJob, *mapReduceJob)
		}
	})
}
