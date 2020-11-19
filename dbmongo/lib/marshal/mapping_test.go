package marshal

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

func TestGetSiret(t *testing.T) {

	stdTime1, _ := time.Parse("2006-02-01", "2015-01-01")
	stdTime2, _ := time.Parse("2006-02-01", "2016-01-01")
	var stdMapping = map[string][]SiretDate{
		"abc": {
			{"01234567891011", stdTime1},
			{"87654321091011", stdTime2},
		},
	}
	var batch = base.AdminBatch{}
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

		var cache = Cache{"comptes": tc.mapping}

		time, _ := time.Parse("2006-02-01", tc.date)
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

	var batch = base.AdminBatch{}

	stdTime1, _ := time.Parse("2006-02-01", "2899-01-01")
	stdTime2, _ := time.Parse("2006-02-01", "2015-01-01")
	stdExpected1 := Comptes{
		"abc": []SiretDate{{"01234567891011", stdTime1}},
	}

	stdExpected2 := Comptes{
		"abc": []SiretDate{{"01234567891011", stdTime2}},
	}

	stdExpected3 := Comptes{
		"abc": []SiretDate{
			{"01234567891011", stdTime2},
			{"87654321091011", stdTime1},
		},
	}

	stdFilterCache := Cache{
		"filter": SirenFilter{"012345678": true},
	}

	testCases := []struct {
		csv         string
		cache       Cache
		expectError bool
		expected    Comptes
	}{
		// No closing date
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;""`, Cache{}, false, stdExpected1},
		// With closing date
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;"1150101"`, Cache{}, false, stdExpected2},
		// With filtered siret
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;"1150101"`, stdFilterCache, false, stdExpected2},
		// With two entries, including excluded siret
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;"1150101"
		;;"abc";;;"87654321091011";;""`, stdFilterCache, false, stdExpected2}, // i.e. no mapping stored for 87654321091011, because it's not included by Filter
		// With two entries 1
		{`0;1;2;3;4;5;6;7
		;;"abc";;;"01234567891011";;"1150101"
		;;"abc";;;"87654321091011";;""`, Cache{}, false, stdExpected3},
		// With two entries 2: different order
		{`0;1;2;3;4;5;6;7
	    ;;"abc";;;"87654321091011";;""
	    ;;"abc";;;"01234567891011";;"1150101"`, Cache{}, false, stdExpected3},
		// With invalid siret
		{`0;1;2;3;4;5;6;7
		  ;;"abc";;;"8765432109101A";;""
	    ;;"abc";;;"01234567891011";;"1150101"`, Cache{}, false, stdExpected2},
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
		"abc": []SiretDate{{"01234567891011", stdTime1}},
	}
	stdExpected2 := Comptes{
		"abc": []SiretDate{{"01234567891011", stdTime2}},
	}

	// When file is read, returnd stdExpected1
	mockOpenFile := func(s1 string, s2 string, c Comptes, ca Cache, ba *base.AdminBatch) (Comptes, error) {
		for key := range stdExpected1 {
			c[key] = stdExpected1[key]
		}
		return c, nil
	}

	testCases := []struct {
		cache       Cache
		batch       base.AdminBatch
		expectError bool
		expected    Comptes
	}{
		// Basic reading from file
		{NewCache(), base.MockBatch("admin_urssaf", []string{"a"}), false, stdExpected1},
		// Cache superseeds reading from file
		{Cache{"comptes": stdExpected2}, base.MockBatch("admin_urssaf", []string{"a"}), false, stdExpected2},
		// No cache, no file = error
		{NewCache(), base.MockBatch("otherStuff", []string{"a"}), true, nil},
	}

	for ind, tc := range testCases {

		actual, err := GetCompteSiretMapping(tc.cache, &tc.batch, mockOpenFile)
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
