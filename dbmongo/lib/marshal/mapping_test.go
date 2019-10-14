package marshal

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
)

func TestGetSiret(t *testing.T) {

	stdTime1, _ := time.Parse("2006-02-01", "2015-01-01")
	stdTime2, _ := time.Parse("2006-02-01", "2016-01-01")
	var stdMapping = map[string][]SiretDate{
		"abc": []SiretDate{
			SiretDate{"01234567891011", stdTime1},
			SiretDate{"87654321091011", stdTime2},
		},
	}
	var batch = engine.AdminBatch{}
	testCases := []struct {
		compte      string
		date        string
		mapping     Comptes
		expectError bool
		expected    string
	}{
		{"abc", "2015-01-01", map[string][]SiretDate{}, true, ""},
		{"abc", "2014-06-01", stdMapping, false, "01234567891011"},
		{"abc", "2014-12-31", stdMapping, false, "01234567891011"},
		{"abc", "2015-01-01", stdMapping, false, "87654321091011"},
		{"abc", "2016-01-01", stdMapping, true, ""},
	}

	for ind, tc := range testCases {

		var cache = engine.Cache{"comptes": tc.mapping}

		time, err := time.Parse("2006-02-01", tc.date)
		actual, err := GetSiret(tc.compte, &time, cache, &batch)
		if err != nil && !tc.expectError {
			t.Fatalf("Unexpected error during cache request in test %d: %v", ind, err)
		}
		if err == nil && tc.expectError {
			t.Fatalf("Expected error missing during cache request in test %d: %v", ind, err)
		}
		expected := tc.expected
		if actual != expected {
			t.Log(actual)
			t.Log(expected)
			t.Fatalf("GetSiret failed on test %d", ind)
		}
	}
}

func TestReadSiretMapping(t *testing.T) {

	var batch = engine.AdminBatch{}

	stdTime1, _ := time.Parse("2006-02-01", "2899-01-01")
	stdTime2, _ := time.Parse("2006-02-01", "2015-01-01")
	stdExpected1 := Comptes{
		"abc": []SiretDate{SiretDate{"01234567891011", stdTime1}},
	}
	stdExpected2 := Comptes{
		"abc": []SiretDate{SiretDate{"01234567891011", stdTime2}},
	}
	stdExpected3 := Comptes{
		"abc": []SiretDate{
			SiretDate{"01234567891011", stdTime2},
			SiretDate{"87654321091011", stdTime1},
		},
	}
	stdFilterCache := engine.Cache{
		"filter": map[string]bool{"01234567891011": true},
	}

	testCases := []struct {
		csv         string
		cache       engine.Cache
		expectError bool
		expected    Comptes
	}{
		// No closing date
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;""`, engine.Cache{}, false, stdExpected1},
		// With closing date
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;"1150101"`, engine.Cache{}, false, stdExpected2},
		// With filtered siret
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;"1150101"`, stdFilterCache, false, map[string][]SiretDate{}},
		// With two entries 1
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;"1150101"
		;;"abc";;;"87654321091011";;""`, engine.Cache{}, false, stdExpected3},
		// With two entries 2: different order
		{`0;1;2;3;4;5;6;7
	    ;;"abc";;;"87654321091011";;""
	    ;;"abc";;;"01234567891011";;"1150101"`, engine.Cache{}, false, stdExpected3},
	}

	for ind, tc := range testCases {
		actual, err := readSiretMapping(strings.NewReader(tc.csv), tc.cache, &batch)
		if err != nil && !tc.expectError {
			t.Fatalf("unexpected error during file reading in test %d: %v", ind, err)
		}
		if err == nil && tc.expectError {
			t.Fatalf("expected error missing during  file reading in test %d: %v", ind, err)
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Logf("Actual %v", actual)
			t.Logf("Expected %v", tc.expected)
			t.Fatalf("ReadSiretMapping failed test %d", ind)
		}
	}
}

func TestGetCompteSiretMapping(t *testing.T) {
	stdTime1, _ := time.Parse("2006-02-01", "2899-01-01")
	stdTime2, _ := time.Parse("2006-02-01", "2016-01-01")
	stdExpected1 := Comptes{
		"abc": []SiretDate{SiretDate{"01234567891011", stdTime1}},
	}
	stdExpected2 := Comptes{
		"abc": []SiretDate{SiretDate{"01234567891011", stdTime2}},
	}
	mockOpenFile := func(s1 string, s2 string, c Comptes, ca engine.Cache, ba *engine.AdminBatch) (Comptes, error) {
		for key := range stdExpected1 {
			c[key] = stdExpected1[key]
		}
		return c, nil
	}
	testCases := []struct {
		cache       engine.Cache
		batch       engine.AdminBatch
		expectError bool
		expected    Comptes
	}{
		// Basic reading from file
		{nil, engine.MockBatch("admin_urssaf", []string{"a"}), false, stdExpected1},
		// Cache superseeds reading from file
		{engine.Cache{"comptes": stdExpected2}, engine.MockBatch("admin_urssaf", []string{"a"}), false, stdExpected2},
		// No cache, no file = error
		{nil, engine.MockBatch("otherStuff", []string{"a"}), true, nil},
	}

	for ind, tc := range testCases {

		actual, err := getCompteSiretMapping(tc.cache, &tc.batch, mockOpenFile)
		if err != nil && !tc.expectError {
			t.Fatalf("Unexpected error during mapping request in test %d: %v", ind, err)
		}
		if err == nil && tc.expectError {
			t.Fatalf("Expected error missing during mapping request in test %d: %v", ind, err)
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Logf("Actual %v", actual)
			t.Logf("Expected %v", tc.expected)
			t.Fatalf("getSiretMapping failed test %d", ind)
		}
	}
}
