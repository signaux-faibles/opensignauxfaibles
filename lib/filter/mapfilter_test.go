package filter

import (
	"flag"
	"testing"
)

var _ = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

func TestIncludes(t *testing.T) {

	testCases := []struct {
		name     string
		siret    string
		filter   MapFilter
		expected bool
	}{
		{
			"siren inclus dans filtre",
			"012345678",
			MapFilter{"012345678": true},
			true,
		},
		{
			"siret inclus dans filtre",
			"01234567891011",
			MapFilter{"012345678": true},
			true,
		},
		{
			"siren trop court",
			"0123",
			MapFilter{"012345678": true},
			false,
		},
		{
			"numéro invalide mais ayant comme prefixe un siren filtré",
			"0123456789",
			MapFilter{"012345678": true},
			false,
		},
		{
			"siren non inclus dans filtre",
			"876543210",
			MapFilter{"012345678": true},
			false,
		},
		{
			"siret non inclus dans filtre",
			"87654321091011",
			MapFilter{"012345678": true},
			false,
		},
		{
			"pas de filtre",
			"012345678",
			nil,
			true,
		},
		{
			"pas de filtre + numéro invalide",
			"0123",
			nil,
			false,
		},
	}

	for ind, tc := range testCases {
		included := !tc.filter.ShouldSkip(tc.siret)
		expected := tc.expected
		if included != expected {
			t.Fatalf("Includes failed on test %d: %s", ind, tc.name)
		}
	}
}
