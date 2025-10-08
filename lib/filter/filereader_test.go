package filter

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadFilter(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected SirenFilter
		wantErr  bool
	}{
		{
			name:     "valid sirens",
			input:    "012345678\n876543210",
			expected: SirenFilter{"012345678": true, "876543210": true},
			wantErr:  false,
		},
		{
			name:     "invalid siren (wrong length)",
			input:    "0123456789\n876543210",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid header (or invalid siren)",
			input:    "abcdefghi\n876543210",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "valid siren with header",
			input:    "siren\n012345678",
			expected: SirenFilter{"012345678": true},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testFilter := make(SirenFilter)
			err := parseCSVFilter(strings.NewReader(tc.input), testFilter)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("readFilter should fail on incorrect siren")
				}
				return
			}

			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			if !reflect.DeepEqual(testFilter, tc.expected) {
				t.Fatalf("Filter not read as expected, failure")
			}
		})
	}
}
