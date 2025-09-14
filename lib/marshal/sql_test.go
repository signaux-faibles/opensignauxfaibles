package marshal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractValues(t *testing.T) {
	anInt := 1
	testCases := []struct {
		tuple          TestTuple
		expectedLen    int
		expectedValues []any
	}{
		{
			TestTuple{},
			3,
			[]any{"", (*int)(nil), (*time.Time)(nil)},
		},
		{
			TestTuple{"abc", &anInt, "def", &time.Time{}},
			3,
			[]any{"abc", anInt, time.Time{}},
		},
	}
	for _, tc := range testCases {

		extracted := ExtractTableRow(tc.tuple)
		assert.Len(t, extracted, tc.expectedLen)
		for i, expected := range tc.expectedValues {
			assert.Equal(t, extracted[i], expected)
		}
	}
}

func TestExtractTableColumns(t *testing.T) {
	tuple := TestTuple{}
	assert.Equal(t, ExtractTableColumns(tuple), []string{"test1", "test2", "test4"})
}
