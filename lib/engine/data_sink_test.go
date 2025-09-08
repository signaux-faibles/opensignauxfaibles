package engine

import (
	"context"
	"opensignauxfaibles/lib/marshal"
	"testing"

	"github.com/stretchr/testify/assert"
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

	ctx := context.Background()
	compSink, _ := compositeFactory.CreateSink("testtype")
	compSink.ProcessOutput(ctx, ch)

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

func TestCompositeSink_FailingSink(t *testing.T) {
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
	f2 := FailSinkFactory{}
	compositeFactory := &compositeSinkFactory{[]SinkFactory{f1, f2}, 1}

	ctx := context.Background()
	compSink, _ := compositeFactory.CreateSink("testtype")
	err := compSink.ProcessOutput(ctx, ch)

	assert.Error(t, err)
}
