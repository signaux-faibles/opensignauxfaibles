package marshal

import (
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"

	"github.com/cnf/structhash"
	"github.com/globalsign/mgo/bson"
)

func TestCheckMarshallingMap(t *testing.T) {
	headerRow := []string{"a", "b", "c", "d", "e"}
	cases := []struct {
		marshallingMap map[string]int
		expected       bool
	}{
		{map[string]int{"a": 0, "d": 3}, true},
		{map[string]int{"a": 0, "f": 4}, false},
		{map[string]int{"a": 1, "c": 2}, false},
	}

	for ind, tc := range cases {
		actualErr := CheckMarshallingMap(headerRow, tc.marshallingMap)
		if (actualErr == nil) != tc.expected {
			t.Errorf("Test fails on case %d", ind)
		}
	}
}

func TestMD5(t *testing.T) {
	var test Object
	test.key = "0123456789"
	test.scope = "abc"
	test.datatype = "ced"
	test.Data = make(bson.M)
	test.Data["test"] = 3

	test2 := struct{ test int }{test: 3}
	test3 := map[string]int{"test": 3}

	t.Log(engine.GetMD5(test))
	t.Log(structhash.Md5(test.Data, 1))
	t.Log(structhash.Md5(test2, 1))
	t.Log(structhash.Md5(test3, 1))
}
