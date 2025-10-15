package sinks

import (
	"bytes"
	"context"
	"opensignauxfaibles/lib/engine"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockTuple implements the Tuple interface for testing
// H1, H2, H3 are three mock data
type MockTuple struct {
	H1    string `csv:"header1"`
	H2    string `csv:"header2"`
	H3    string `csv:"header3"`
	key   string
	scope engine.Scope
	tType engine.ParserType
}

func (m MockTuple) Key() string           { return m.key }
func (m MockTuple) Scope() engine.Scope     { return m.scope }
func (m MockTuple) Type() engine.ParserType { return m.tType }

func TestCSVSink_ProcessOutput(t *testing.T) {

	var buf bytes.Buffer
	sink := &CSVSink{
		writer: &buf,
	}

	// setup channel
	ch := make(chan engine.Tuple, 3)
	ch <- MockTuple{
		H1:    "value1",
		H2:    "value2",
		H3:    "value3",
		key:   "123456789",
		scope: engine.ScopeEntreprise,
		tType: "testtype",
	}
	ch <- MockTuple{
		H1:    "value4",
		H2:    "value5",
		H3:    "value6",
		key:   "987654321",
		scope: engine.ScopeEntreprise,
		tType: "testtype",
	}
	close(ch)

	ctx := context.Background()
	// Stream to writer
	err := sink.ProcessOutput(ctx, ch)

	if err != nil {
		t.Error("Expect data to be streamed to writer")
	}

	// Expected output
	expectedOutput := "header1,header2,header3\nvalue1,value2,value3\nvalue4,value5,value6\n"
	if buf.String() != expectedOutput {
		t.Error("Expect output data to be properly csv formatted")
	}
}

func TestExtractCSVValues(t *testing.T) {
	anInt := 1
	testCases := []struct {
		tuple          engine.TestTuple
		expectedLen    int
		expectedValues []string
	}{
		{
			engine.TestTuple{},
			3,
			[]string{"", "", ""},
		},
		{
			engine.TestTuple{"abc", &anInt, "def", &time.Time{}},
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
	tuple := engine.TestTuple{}
	assert.Equal(t, ExtractCSVHeaders(tuple), []string{"test1", "test2", "test4"})
}
