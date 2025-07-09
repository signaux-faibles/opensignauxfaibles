package engine

import (
	"opensignauxfaibles/lib/marshal"

	"golang.org/x/sync/errgroup"
)

// An OutputStreamer directs a stream of output data to the desired sink
type OutputStreamer interface {
	// Stream is expected to be synchronous. Any concurrency is handled by the
	// caller.
	Stream(ch chan marshal.Tuple) error
}

// A CombinedStreamer streams to two different sinks
type CombinedStreamer struct {
	out1 OutputStreamer
	out2 OutputStreamer
}

func NewCombinedStreamer(out1, out2 OutputStreamer) *CombinedStreamer {
	return &CombinedStreamer{out1, out2}
}

// Stream dispatches data to two sinks, with two `OutputStreamer`s that run concurrently.
func (combined *CombinedStreamer) Stream(ch chan marshal.Tuple) error {

	outCh1 := make(chan marshal.Tuple, 1000)
	outCh2 := make(chan marshal.Tuple, 1000)

	go func() {
		defer close(outCh1)
		defer close(outCh2)

		for tuple := range ch {
			// tuple is immutable and does not need a copy
			outCh1 <- tuple
			outCh2 <- tuple
		}
	}()

	var g errgroup.Group

	g.Go(
		func() error {
			err := combined.out1.Stream(outCh1)
			return err
		},
	)

	g.Go(
		func() error {
			err := combined.out2.Stream(outCh2)
			return err
		},
	)

	err := g.Wait()

	return err
}
