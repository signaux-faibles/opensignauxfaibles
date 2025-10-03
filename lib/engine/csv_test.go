package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractCSVValues(t *testing.T) {
	anInt := 1
	testCases := []struct {
		tuple          TestTuple
		expectedLen    int
		expectedValues []string
	}{
		{
			TestTuple{},
			3,
			[]string{"", "", ""},
		},
		{
			TestTuple{"abc", &anInt, "def", &time.Time{}},
			3,
			[]string{"abc", "1", "0001-01-01"},
		},
	}
	for _, tc := range testCases {

		extracted := ExtractCSVRow(tc.tuple)

		if assert.Len(t, extracted, tc.expectedLen) {
			for i, expected := range tc.expectedValues {
				assert.Equal(t, extracted[i], expected)
			}
		}
	}
}

func TestExtractCSVHeaders(t *testing.T) {
	tuple := TestTuple{}
	assert.Equal(t, ExtractCSVHeaders(tuple), []string{"test1", "test2", "test4"})
}
