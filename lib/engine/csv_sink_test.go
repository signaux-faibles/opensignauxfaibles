package engine

import (
	"bytes"
	"opensignauxfaibles/lib/marshal"
	"testing"
)

// MockTuple implements the marshal.Tuple interface for testing
// H1, H2, H3 are three mock data
type MockTuple struct {
	H1    string `csv:"header1"`
	H2    string `csv:"header2"`
	H3    string `csv:"header3"`
	key   string
	scope string
	tType string
}

func (m MockTuple) Key() string   { return m.key }
func (m MockTuple) Scope() string { return m.scope }
func (m MockTuple) Type() string  { return m.tType }

func TestCSVSink_ProcessOutput(t *testing.T) {

	var buf bytes.Buffer
	sink := &CSVSink{
		writer: &buf,
	}

	// setup channel
	ch := make(chan marshal.Tuple, 3)
	ch <- MockTuple{
		H1:    "value1",
		H2:    "value2",
		H3:    "value3",
		key:   "123456789",
		scope: "entreprise",
		tType: "testtype",
	}
	ch <- MockTuple{
		H1:    "value4",
		H2:    "value5",
		H3:    "value6",
		key:   "987654321",
		scope: "entreprise",
		tType: "testtype",
	}
	close(ch)

	// Stream to writer
	err := sink.ProcessOutput(ch)

	if err != nil {
		t.Error("Expect data to be streamed to writer")
	}

	// Expected output
	expectedOutput := "header1,header2,header3\nvalue1,value2,value3\nvalue4,value5,value6\n"
	if buf.String() != expectedOutput {
		t.Error("Expect output data to be properly csv formatted")
	}
}
