package engine

import (
	"testing"

	"github.com/globalsign/mgo/bson"
)

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
