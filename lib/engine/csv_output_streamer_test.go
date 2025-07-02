package engine

import (
	"bytes"
	"opensignauxfaibles/lib/marshal"
	"testing"
)

// mockTuple implements the marshal.Tuple interface for testing
type mockTuple struct {
	key     string
	scope   string
	tType   string
	headers []string
	values  []string
}

func (m mockTuple) Key() string       { return m.key }
func (m mockTuple) Scope() string     { return m.scope }
func (m mockTuple) Type() string      { return m.tType }
func (m mockTuple) Headers() []string { return m.headers }
func (m mockTuple) Values() []string  { return m.values }

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
		key:     "123456789",
		scope:   "entreprise",
		tType:   "testtype",
		headers: []string{"header1", "header2", "header3"},
		values:  []string{"value1", "value2", "value3"},
	}
	ch <- mockTuple{
		key:     "987654321",
		scope:   "entreprise",
		tType:   "testtype",
		headers: []string{"header1", "header2", "header3"},
		values:  []string{"value4", "value5", "value6"},
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
