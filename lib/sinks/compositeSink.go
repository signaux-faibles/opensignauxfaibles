package sinks

import (
	"context"
	"fmt"
	"opensignauxfaibles/lib/engine"
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

// NewCompositeSinkFactory gives a SinkFactory, that creates DataSink
// instances that combine multiple sinks
func Combine(factories ...engine.SinkFactory) engine.SinkFactory {
	return &compositeSinkFactory{
		factories:  factories,
		bufferSize: DefaultBufferSize,
	}
}

type compositeSinkFactory struct {
	factories  []engine.SinkFactory
	bufferSize int
}

func (f *compositeSinkFactory) CreateSink(parserType engine.ParserType) (engine.DataSink, error) {
	var sinks []engine.DataSink
	for _, factory := range f.factories {
		sink, err := factory.CreateSink(parserType)
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, sink)
	}
	return &compositeSink{sinks, f.bufferSize}, nil
}

type compositeSink struct {
	sinks      []engine.DataSink
	bufferSize int
}

func (s *compositeSink) ProcessOutput(ctx context.Context, ch chan engine.Tuple) error {

	// Creates a new context for the ability to cancel all sinks if any sink
	// fails. For the moment this is an intended and acceptable behavior.
	subctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	var outChannels []chan engine.Tuple

	// We duplicate the channels
	for range s.sinks {
		outChannels = append(outChannels, make(chan engine.Tuple, s.bufferSize))
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

// getType return the name of the type of the input
func getType(myvar any) string {
	return reflect.TypeOf(myvar).String()
}
