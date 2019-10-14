package marshal

import (
	"testing"
)

// Helper functions must return same type as they take in input.
// They must be able to deal with zero values of this type.

func TestDepartement(t *testing.T) {
	testCases := []struct {
		codePostal string
		expected   string
	}{
		{"87300", "87"},
		{"87", ""},
		{"", ""},
	}

	for ind, tc := range testCases {
		if Departement(tc.codePostal) != tc.expected {
			t.Errorf("Test %d failed", ind)
		}
	}
}

func TestStripPoint(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"87.300", "87300"},
		{"22.18A", "2218A"},
		{"2.2.1.8.A", "2218A"},
		{"", ""},
	}

	for ind, tc := range testCases {
		if StripPoint(tc.input) != tc.expected {
			t.Errorf("Test %d failed", ind)
		}
	}
}

func TestDivide100(t *testing.T) {
	if Divide100(nil) != nil {
		t.Errorf("Test 0 failed")
	}
	var a = 520.3
	var b = 5.203
	if *Divide100(&a)-b > 1e-8 {
		t.Errorf("Test 1 failed")
	}
}
