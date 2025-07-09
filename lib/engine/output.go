package engine

import (
	"errors"
	"opensignauxfaibles/lib/marshal"
	"sync"
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

	// All errors (max 2) are buffered as they are processed only after both
	// streams have ended
	errCh := make(chan error, 2)

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

	errs := []error{}
	for err := range errCh {
		errs = append(errs, err)
	}

	var err error
	switch len(errs) {
	case 2:
		err = errors.Join(errs[0], errs[1])
	case 1:
		err = errs[0]
	case 0:
		err = nil
	}

	return err
}
