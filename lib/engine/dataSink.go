package engine

import (
	"context"
)

// SinkFactory creates DataSink instances configured for specific parser types
type SinkFactory interface {
	// CreateSink returns a new DataSink instance configured for the given parser type
	CreateSink(parserType ParserType) (DataSink, error)
}

// A DataSink directs a stream of output data to the desired sink
type DataSink interface {
	// ProcessOutput reads from the input channel and writes to the sink
	//
	// It is expected to be synchronous. Any concurrency is handled by the
	// caller.
	//
	// The channel ch must be completely consumed.
	ProcessOutput(ctx context.Context, ch chan Tuple) error
}
