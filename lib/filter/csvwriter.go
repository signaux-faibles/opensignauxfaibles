package filter

import (
	"fmt"
	"io"
	"opensignauxfaibles/lib/engine"
)

// CsvWriter implements engine.FilterWriter to output filters as CSV
type CsvWriter struct {
	w io.Writer
}

// NewCsvWriter creates a new CsvFilterWriter
func NewCsvWriter(w io.Writer) *CsvWriter {
	return &CsvWriter{w: w}
}

// Write outputs the filter as CSV to the writer
func (c *CsvWriter) Write(f engine.SirenFilter) error {
	sirens := f.All()
	fmt.Fprintln(c.w, "siren")
	for siren := range sirens {
		fmt.Fprintln(c.w, siren)
	}
	return nil
}
