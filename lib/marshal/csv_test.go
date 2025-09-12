package marshal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestCSVTuple struct {
	Test1 string `csv:"test1"`
	Test2 *int   `csv:"test2"`
	Test3 string
	Test4 *time.Time `csv:"test4"`
}

func (TestCSVTuple) Key() string   { return "" }
func (TestCSVTuple) Scope() string { return "" }
func (TestCSVTuple) Type() string  { return "" }

func TestExtractCSVValues(t *testing.T) {
	anInt := 1
	testCases := []struct {
		tuple          TestCSVTuple
		expectedLen    int
		expectedValues []string
	}{
		{
			TestCSVTuple{},
			3,
			[]string{"", "", ""},
		},
		{
			TestCSVTuple{"abc", &anInt, "def", &time.Time{}},
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
