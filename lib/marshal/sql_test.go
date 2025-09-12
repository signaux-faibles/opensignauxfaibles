package marshal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestSQLTuple struct {
	Test1 string `sql:"test1"`
	Test2 *int   `sql:"test2"`
	Test3 string
	Test4 *time.Time `sql:"test4"`
}

func (TestSQLTuple) Key() string   { return "" }
func (TestSQLTuple) Scope() string { return "" }
func (TestSQLTuple) Type() string  { return "" }

func TestExtractValues(t *testing.T) {
	anInt := 1
	testCases := []struct {
		tuple          TestSQLTuple
		expectedLen    int
		expectedValues []any
	}{
		{
			TestSQLTuple{},
			3,
			[]any{"", (*int)(nil), (*time.Time)(nil)},
		},
		{
			TestSQLTuple{"abc", &anInt, "def", &time.Time{}},
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
