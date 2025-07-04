package engine

import (
	"opensignauxfaibles/lib/marshal"
	"sync"
)

// An OutputStreamer directs a stream of output data to the desired sink
type OutputStreamer interface {
	// Stream is expected to be synchronous. Any concurrency is handled by the
	// caller.
	Stream(ch chan marshal.Tuple) error
}

func NewCombinedStreamer(out1, out2 OutputStreamer) *CombinedStreamer {
	return &CombinedStreamer{out1, out2}

}

// A CombinedStreamer streams to two different sinks
type CombinedStreamer struct {
	out1 OutputStreamer
	out2 OutputStreamer
}

// Stream dispatches data to the two OutputStreamer
// The returned error is the first error encountered, if any.
func (combined *CombinedStreamer) Stream(ch chan marshal.Tuple) error {

	outCh1 := make(chan marshal.Tuple, 10)
	outCh2 := make(chan marshal.Tuple, 10)

	errCh := make(chan error, 3)

	go func() {
		defer close(outCh1)
		defer close(outCh2)

		for tuple := range ch {
			// tuple is immutable and does not need a copy
			outCh1 <- tuple
			outCh2 <- tuple
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := combined.out1.Stream(outCh1)
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := combined.out2.Stream(outCh2)
		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	// When all go routines have ended we can safely close errCh
	close(errCh)

	err, ok := <-errCh
	if ok {
		return err
	}

	return nil
}
