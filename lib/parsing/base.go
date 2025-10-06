package parsing

import (
	"encoding/csv"
	"io"
	"opensignauxfaibles/lib/engine"
)

// ParseHandler is implemented by concrete parsers to handle their specific parsing logic
type ParseHandler interface {
	ParseRow(row []string, result *engine.ParsedLineResult) error
}

// BaseParser provides common Parser implementation
type BaseParser struct {
	ParseHandler

	Reader *csv.Reader
}

func (b *BaseParser) ReadNext(res *engine.ParsedLineResult) error {
	row, err := b.Reader.Read()
	if err == io.EOF {
		return err
	}

	if err != nil {
		// Do not interrupt parsing (`return err`) if a single line is malformed
		res.AddRegularError(err)
		return nil
	}

	return b.ParseRow(row, res)
}
