package engine

import (
	"context"
	"opensignauxfaibles/lib/marshal"

	"golang.org/x/sync/errgroup"
)

// SinkFactory creates DataSink instances configured for specific parser types
type SinkFactory interface {
	// CreateSink returns a new DataSink instance configured for the given parser type
	CreateSink(parserType string) (DataSink, error)
}

// A DataSink directs a stream of output data to the desired sink
type DataSink interface {
	// ProcessOutput reads from the input channel and writes to the sink
	//
	// It is expected to be synchronous. Any concurrency is handled by the
	// caller.
	//
	// The channel ch must be completely consumed.
	ProcessOutput(ctx context.Context, ch chan marshal.Tuple) error
}

// NewCompositeSinkFactory gives a SinkFactory, that creates DataSink
// instances that combine multiple sinks
func NewCompositeSinkFactory(factories ...SinkFactory) SinkFactory {
	return &compositeSinkFactory{
		factories: factories,
	}
}

type compositeSinkFactory struct {
	factories []SinkFactory
}

func (f *compositeSinkFactory) CreateSink(parserType string) (DataSink, error) {
	var sinks []DataSink
	for _, factory := range f.factories {
		sink, err := factory.CreateSink(parserType)
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, sink)
	}
	return &compositeSink{sinks}, nil
}

type compositeSink struct {
	sinks []DataSink
}

func (s *compositeSink) ProcessOutput(ctx context.Context, ch chan marshal.Tuple) error {

	// Creates a new context for the ability to cancel all sinks if any sink
	// fails. For the moment this is an intended and acceptable behavior.
	subctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var outChannels []chan marshal.Tuple

	// We duplicate the channels
	for range s.sinks {
		outChannels = append(outChannels, make(chan marshal.Tuple, 1000))
	}

	go func() {
		for _, outCh := range outChannels {
			defer close(outCh)
		}

		for tuple := range ch {
			for _, outCh := range outChannels {
				select {
				case <-ctx.Done(): // something went wrong upstream or downstream
					return
				case <-subctx.Done(): // something went wrong with a sink
					return
				case outCh <- tuple:
				}
			}
		}
	}()

	var g errgroup.Group

	for i, sink := range s.sinks {
		g.Go(
			func() error {
				err := sink.ProcessOutput(subctx, outChannels[i])

				if err != nil {
					cancel() // cancels all sinks
				}

				return err
			},
		)
	}

	err := g.Wait()

	return err
}

type DiscardDataSink struct {
	counter int
}

func (s *DiscardDataSink) ProcessOutput(ctx context.Context, ch chan marshal.Tuple) error {
	for range ch {
		s.counter++
	}
	return nil
}

type DiscardReportSink struct {
	counter int
}

func (s DiscardReportSink) Process(ch chan marshal.Report) error {
	for range ch {
		s.counter++
	}
	return nil
}
