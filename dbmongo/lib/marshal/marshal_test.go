package marshal

import (
	"fmt"
	"opensignauxfaibles/dbmongo/lib/engine"
	"testing"

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

	fmt.Println(engine.GetMD5(test))
}
