package marshal

import (
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

func TestIncludes(t *testing.T) {

	testCases := []struct {
		siret    string
		filter   SirenFilter
		expected bool
	}{
		{"012345678", SirenFilter{"012345678": true}, true},       // siren inclus dans filtre
		{"01234567891011", SirenFilter{"012345678": true}, true},  // siret inclus dans filtre
		{"0123", SirenFilter{"012345678": true}, false},           // siren trop court
		{"0123456789", SirenFilter{"012345678": true}, true},      // numéro invalide mais ayant comme prefixe un siret filtré
		{"876543210", SirenFilter{"012345678": true}, false},      // siren non inclus dans filtre
		{"87654321091011", SirenFilter{"012345678": true}, false}, // siret non inclus dans filtre
		{"012345678", nil, false},                                 // pas de filtre
		{"0123", nil, false},                                      // pas de filtre + numéro invalide
	}

	for ind, tc := range testCases {
		included := tc.filter.includes(tc.siret)
		expected := tc.expected
		if included != expected {
			t.Fatalf("Includes failed on test %d", ind)
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

func TestNilFilter(t *testing.T) {
	filter := GetSirenFilterFromCache(Cache{}) // => nil
	assert.Equal(t, true, filter == nil)
	assert.Equal(t, false, filter.includes("012345678"))
	assert.Equal(t, false, filter.Skips("012345678"))
	assert.Equal(t, false, filter.Skips("912345678"))
	cache := Cache{}
	cache.Set("filter", SirenFilter{"012345678": true})
	filter = GetSirenFilterFromCache(cache) // => not nil
	assert.Equal(t, false, filter == nil)
	assert.Equal(t, true, filter.includes("012345678"))
	assert.Equal(t, false, filter.Skips("012345678"))
	assert.Equal(t, true, filter.Skips("912345678"))
}
