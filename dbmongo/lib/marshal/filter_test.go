package marshal

import (
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

func TestIsFiltered(t *testing.T) {

	testCases := []struct {
		siret       string
		filter      SirenFilter
		expectError bool
		expected    bool
	}{
		{"012345678", SirenFilter{"012345678": true}, false, false},
		{"01234567891011", SirenFilter{"012345678": true}, false, false},
		{"0123", SirenFilter{"012345678": true}, true, true},
		{"0123456789", SirenFilter{"012345678": true}, true, true},
		{"876543210", SirenFilter{"012345678": true}, false, true},
		{"87654321091011", SirenFilter{"012345678": true}, false, true},
		{"012345678", nil, false, false},
		{"0123", nil, true, true},
	}

	for ind, tc := range testCases {
		actual, err := IsFiltered(tc.siret, tc.filter)
		if err != nil && !tc.expectError {
			t.Fatalf("Unexpected error during cache request in test %d: %v", ind, err)
		}
		if err == nil && tc.expectError {
			t.Fatalf("Expected error missing during cache request in test %d: %v", ind, err)
		}
		expected := tc.expected
		if actual != expected {
			t.Fatalf("IsFiltered failed on test %d", ind)
		}
	}
}

func TestGetSirenFilter(t *testing.T) {

	testCases := []struct {
		experimentName string
		cacheKey       string
		cacheValue     interface{}
		batch          base.AdminBatch
		expectedFilter SirenFilter
	}{
		{"existing cache",
			"filter", SirenFilter{"012345678": true}, base.AdminBatch{}, SirenFilter{"012345678": true}},
		{"No cache, no filter in batch 1",
			"", "", base.AdminBatch{}, nil},
		{"No cache, no filter in batch 2",
			"", "", base.MockBatch("filter", nil), nil},
		{"No cache, (mock)read from file",
			"", "", base.MockBatch("filter", []string{"at least one"}), SirenFilter{"012345678": true}},
		{"Cache has precedence over file",
			"filter", SirenFilter{"876543210": true}, base.MockBatch("filter", []string{"at least one"}), SirenFilter{"876543210": true}},
	}

	for ind, tc := range testCases {
		cache := NewCache()
		cache.Set(tc.cacheKey, tc.cacheValue)
		actual, err := getSirenFilter(cache, &tc.batch, mockReadFilter)
		if err != nil {
			t.Fatalf("getSirenFilter had an unexpected error on test %d: %v", ind, err)
		}

		// filter returns as expected
		if !reflect.DeepEqual(actual, tc.expectedFilter) {
			t.Log(actual)
			t.Log(tc.expectedFilter)
			t.Fatalf("Test %d failed", ind)
		}

		// cache has filter
		if tc.expectedFilter != nil {
			cacheFilter, err := cache.Get("filter")
			cacheFilterTyped, ok := cacheFilter.(SirenFilter)
			if err != nil || !ok || !reflect.DeepEqual(cacheFilterTyped, tc.expectedFilter) {
				t.Fatalf("Test %d failed: filter is missing from cache", ind)
			}
		}
	}
}

func mockReadFilter(string, []string) (SirenFilter, error) {
	return SirenFilter{"012345678": true}, nil
}

func TestReadFilter(t *testing.T) {
	var testFilter = make(SirenFilter)
	err := readFilter(strings.NewReader("012345678\n876543210"), testFilter)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !reflect.DeepEqual(testFilter, SirenFilter{"012345678": true, "876543210": true}) {
		t.Fatalf("Filter not read as expected, failure")
	}

	testFilter = make(SirenFilter)
	err = readFilter(strings.NewReader("0123456789\n876543210"), testFilter)
	if err == nil {
		t.Fatalf("readFilter should fail on incorrect siren")
	}
}
