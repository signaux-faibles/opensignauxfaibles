package engine

import (
	"opensignauxfaibles/lib/marshal"
	"testing"
)

func TestCombinedStreamer_Stream(t *testing.T) {

	ch := make(chan marshal.Tuple)

	go func() {
		ch <- MockTuple{
			key:     "123456789",
			scope:   "entreprise",
			tType:   "testtype",
			headers: []string{"header1", "header2", "header3"},
			values:  []string{"value1", "value2", "value3"},
		}
		ch <- MockTuple{
			key:     "987654321",
			scope:   "entreprise",
			tType:   "testtype",
			headers: []string{"header1", "header2", "header3"},
			values:  []string{"value4", "value5", "value6"},
		}
		close(ch)
	}()

	out1 := NewTestOutputStreamer()
	out2 := NewTestOutputStreamer()
	combined := NewCombinedStreamer(out1, out2)

	combined.Stream(ch)

	if out1.Count != 2 || out2.Count != 2 {
		t.Fatalf("A CombinedStreamer is expected to dispatch all data to both output sinks")
	}

}
