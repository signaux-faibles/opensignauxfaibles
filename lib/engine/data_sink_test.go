package engine

import (
	"opensignauxfaibles/lib/marshal"
	"testing"
)

func TestCompositeSink(t *testing.T) {

	ch := make(chan marshal.Tuple)

	go func() {
		ch <- MockTuple{
			key:   "123456789",
			scope: "entreprise",
			tType: "testtype",
			H1:    "value1",
			H2:    "value2",
			H3:    "value3",
		}
		ch <- MockTuple{
			key:   "987654321",
			scope: "entreprise",
			tType: "testtype",
			H1:    "value4",
			H2:    "value5",
			H3:    "value6",
		}
		close(ch)
	}()

	f1 := TestSinkFactory{}
	f2 := TestSinkFactory{}
	compositeFactory := NewCompositeSinkFactory(f1, f2)

	compSink, _ := compositeFactory.CreateSink("testtype")
	compSink.ProcessOutput(ch)

	allSinks := compSink.(*compositeSink).sinks
	if len(allSinks) != 2 {
		t.Fatalf("The composite sink is expected to store all output stinks it is composed of")
	}

	for _, sink := range allSinks {
		n := sink.(*DiscardDataSink).counter
		if n != 2 {
			t.Fatalf("A composite sink is expected to dispatch all data (2 tuples) to all output sinks, got %d", n)
		}
	}
}
