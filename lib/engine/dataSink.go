package engine

import (
	"context"
	"fmt"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"reflect"

	"golang.org/x/sync/errgroup"
)

// DefaultBufferSize defines the size of the buffer of values that each sink
// has.
// Sinks may proceed input at different paces, and to avoid to have all sinks
// waiting for the slowest one, the slower sinks buffer values while the
// faster ones proceed.
// However, we do not want the buffer to grow unconstrained, so when the
// buffer is full, all sinks have to go at the slowest sink's pace.
const DefaultBufferSize = 100000

// SinkFactory creates DataSink instances configured for specific parser types
type SinkFactory interface {
	// CreateSink returns a new DataSink instance configured for the given parser type
	CreateSink(parserType base.ParserType) (DataSink, error)
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
// instances that combine multiple sinks.
// It also implements the `Finalizer` interface that runs any finalization
// function from each individual sinks (if any)
func NewCompositeSinkFactory(factories ...SinkFactory) SinkFactory {
	return &compositeSinkFactory{
		factories:  factories,
		bufferSize: DefaultBufferSize,
	}
}

type compositeSinkFactory struct {
	factories  []SinkFactory
	bufferSize int
}

func (f *compositeSinkFactory) CreateSink(parserType base.ParserType) (DataSink, error) {
	var sinks []DataSink
	for _, factory := range f.factories {
		sink, err := factory.CreateSink(parserType)
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, sink)
	}
	return &compositeSink{sinks, f.bufferSize}, nil
}

func (f *compositeSinkFactory) Finalize() error {
	for _, factory := range f.factories {
		if finalizer, ok := factory.(Finalizer); ok {
			err := finalizer.Finalize()

			if err != nil {
				return err
			}
		}
	}
	return nil
}

type compositeSink struct {
	sinks      []DataSink
	bufferSize int
}

func (s *compositeSink) ProcessOutput(ctx context.Context, ch chan marshal.Tuple) error {

	// Creates a new context for the ability to cancel all sinks if any sink
	// fails. For the moment this is an intended and acceptable behavior.
	subctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	var outChannels []chan marshal.Tuple

	// We duplicate the channels
	for range s.sinks {
		outChannels = append(outChannels, make(chan marshal.Tuple, s.bufferSize))
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
					err = fmt.Errorf("sink %s failed: %v ; cancelling data processing for all sinks", getType(sink), err)
					cancel(err)
				}

				return err
			},
		)
	}

	err := g.Wait()

	return err
}

// DiscardSinkFactory discards all data, regardless of the parser
type DiscardSinkFactory struct{}

func (f *DiscardSinkFactory) CreateSink(parserType base.ParserType) (DataSink, error) {
	return &DiscardDataSink{}, nil
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

// getType return the name of the type of the input
func getType(myvar any) string {
	return reflect.TypeOf(myvar).String()
}

// Finalizer is an (optional) interface for sink factories that need to run a final
// operation after all processes have completed
type Finalizer interface {
	Finalize() error
}
