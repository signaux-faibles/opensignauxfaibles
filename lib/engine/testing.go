package engine

import (
	"context"
	"fmt"
	"opensignauxfaibles/lib/marshal"
	"time"
)

// DiscardTuple ignore les données
func DiscardTuple(tuples chan marshal.Tuple) {
	go func() {
		for range tuples {
		}
	}()
}

// FailDataSink is a sink that always fails
type FailDataSink struct{}

func (s *FailDataSink) ProcessOutput(ctx context.Context, ch chan marshal.Tuple) error {
	time.Sleep(500 * time.Millisecond)
	return fmt.Errorf("this sink always fails")
}

type FailSinkFactory struct{}

func (FailSinkFactory) CreateSink(parserType string) (DataSink, error) {
	return &FailDataSink{}, nil
}
