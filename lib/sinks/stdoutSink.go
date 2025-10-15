package sinks

import (
	"context"
	"encoding/csv"
	"os"

	"opensignauxfaibles/lib/engine"
)

// StdoutSinkFactory creates sinks that output CSV tuples to stdout
type StdoutSinkFactory struct{}

func NewStdoutSinkFactory() engine.SinkFactory {
	return &StdoutSinkFactory{}
}

func (f *StdoutSinkFactory) CreateSink(parserType engine.ParserType) (engine.DataSink, error) {
	return &StdoutSink{}, nil
}

// StdoutSink outputs each tuple as CSV to stdout
type StdoutSink struct{}

func (s *StdoutSink) ProcessOutput(ctx context.Context, ch chan engine.Tuple) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	headersWritten := false
	for tuple := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if !headersWritten {
				headers := ExtractCSVHeaders(tuple)
				if err := w.Write(headers); err != nil {
					return err
				}
				headersWritten = true
			}
			
			row := ExtractCSVRow(tuple)
			if err := w.Write(row); err != nil {
				return err
			}
		}
	}
	return nil
}