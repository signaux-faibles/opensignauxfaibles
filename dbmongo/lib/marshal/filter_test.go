package marshal

import (
	"fmt"
	"opensignauxfaibles/dbmongo/lib/engine"
	"reflect"
	"strings"
	"testing"
)

func TestIsFiltered(t *testing.T) {

	test_cases := []struct {
		siret       string
		filter      map[string]bool
		expectError bool
		expected    bool
	}{
		{"012345678", map[string]bool{"012345678": true}, false, false},
		{"01234567891011", map[string]bool{"012345678": true}, false, false},
		{"0123", map[string]bool{"012345678": true}, true, false},
		{"0123456789", map[string]bool{"012345678": true}, true, false},
		{"876543210", map[string]bool{"012345678": true}, false, true},
		{"87654321091011", map[string]bool{"012345678": true}, false, true},
		{"012345678", nil, false, false},
		{"0123", nil, true, false},
	}

	var batch = engine.AdminBatch{}
	for ind, tc := range test_cases {
		var cache = engine.Cache{"filter": tc.filter}
		actual, err := IsFiltered(tc.siret, cache, &batch)
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

	test_cases := []struct {
		experimentName string
		cacheKey       string
		cacheValue     interface{}
		batch          engine.AdminBatch
		expectedFilter map[string]bool
	}{
		{"existing cache",
			"filter", map[string]bool{"012345678": true}, engine.AdminBatch{}, map[string]bool{"012345678": true}},
		{"No cache, no filter in batch 1",
			"", "", engine.AdminBatch{}, nil},
		{"No cache, no filter in batch 2",
			"", "", engine.MockBatch("filter", nil), nil},
		{"No cache, (mock)read from file",
			"", "", engine.MockBatch("filter", []string{"at least one"}), map[string]bool{"012345678": true}}, // no filter
		{"Cache has precedence over file",
			"filter", map[string]bool{"876543210": true}, engine.MockBatch("filter", []string{"at least one"}), map[string]bool{"876543210": true}},
	}

	for ind, tc := range test_cases {
		cache := engine.NewCache()
		cache.Set(tc.cacheKey, tc.cacheValue)
		actual, err := getSirenFilter(cache, &tc.batch, mockReadFilter)
		if err != nil {
			t.Fatalf("getSirenFilter had an unexpected error on test %d: %v", ind, err)
		}

		// filter returns as expected
		if !reflect.DeepEqual(actual, tc.expectedFilter) {
			fmt.Println(actual)
			fmt.Println(tc.expectedFilter)
			t.Fatalf("Test %d failed", ind)
		}

		// cache has filter
		if tc.expectedFilter != nil {
			cacheFilter, err := cache.Get("filter")
			cacheFilterTyped, ok := cacheFilter.(map[string]bool)
			if err != nil || !ok || !reflect.DeepEqual(cacheFilterTyped, tc.expectedFilter) {
				t.Fatalf("Test %d failed: filter is missing from cache", ind)
			}
		}
	}
}

func mockReadFilter(string, []string) (map[string]bool, error) {
	return map[string]bool{"012345678": true}, nil
}

func TestReadFilter(t *testing.T) {
	var test_filter = make(map[string]bool)
	err := readFilter(strings.NewReader("012345678\n876543210"), test_filter)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !reflect.DeepEqual(test_filter, map[string]bool{"012345678": true, "876543210": true}) {
		t.Fatalf("Filter not read as expected, failure")
	}

	test_filter = make(map[string]bool)
	err = readFilter(strings.NewReader("0123456789\n876543210"), test_filter)
	if err == nil {
		t.Fatalf("readFilter should fail on incorrect siren")
	}
}
