package engine

import (
	"bytes"
	"opensignauxfaibles/lib/marshal"
	"testing"
)

// mockTuple implements the marshal.Tuple interface for testing
type mockTuple struct {
	H1     string `csv:"header1"`
	H2     string `csv:"header2"`
	H3     string `csv:"header3"`
	key    string
	scope  string
	tType  string
	values []string
}

func (m mockTuple) Key() string       { return m.key }
func (m mockTuple) Scope() string     { return m.scope }
func (m mockTuple) Type() string      { return m.tType }
func (m mockTuple) Headers() []string { return []string{"header1", "header2", "header3"} }
func (m mockTuple) Values() []string  { return []string{m.H1, m.H2, m.H3} }

func TestCSVOutputStreamer_Stream(t *testing.T) {

	// setup streamer func()
	var buf bytes.Buffer
	streamer := &CSVOutputStreamer{
		relativeDirPath: "test",
		writer:          &buf,
	}

	// setup channel
	ch := make(chan marshal.Tuple, 3)
	ch <- mockTuple{
		H1:    "value1",
		H2:    "value2",
		H3:    "value3",
		key:   "123456789",
		scope: "entreprise",
		tType: "testtype",
	}
	ch <- mockTuple{
		H1:    "value4",
		H2:    "value5",
		H3:    "value6",
		key:   "987654321",
		scope: "entreprise",
		tType: "testtype",
	}
	close(ch)

	// Stream to writer
	err := streamer.Stream(ch)

	if err != nil {
		t.Error("Expect data to be streamed to writer")
	}

	// Expected output
	expectedOutput := "header1,header2,header3\nvalue1,value2,value3\nvalue4,value5,value6\n"
	if buf.String() != expectedOutput {
		t.Error("Expect output data to be properly csv formatted")
	}
}
