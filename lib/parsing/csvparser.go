package parsing

import (
	"encoding/csv"
	"io"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
)

// RowParser parses a single row
type RowParser interface {
	ParseRow(row []string, result *engine.ParsedLineResult, idx ColIndex) error
}

// CsvParserInstance provides a CSV Parser implementation base
// Assumes the CSV has a header row
type CsvParserInst struct {
	io.Reader
	RowParser

	Comma      rune
	LazyQuotes bool
	DestTuple  any

	idx ColIndex

	csvReader *csv.Reader

	header []string
}

func (p *CsvParserInst) Init(cache *engine.Cache, filter engine.SirenFilter, batch *base.AdminBatch) error {
	p.csvReader = csv.NewReader(p)
	p.csvReader.Comma = p.Comma
	p.csvReader.LazyQuotes = p.LazyQuotes

	var err error

	// Read first row as header
	p.header, err = p.csvReader.Read()
	if err != nil {
		return err
	}

	p.idx, err = HeaderIndexer{p.DestTuple}.Index(p.Header())
	return err
}

func (p *CsvParserInst) Header() []string {
	return p.header
}

func (p *CsvParserInst) ReadNext(res *engine.ParsedLineResult) error {
	row, err := p.csvReader.Read()
	if err == io.EOF {
		return err
	}

	if err != nil {
		// Do not interrupt parsing (with `return err`) if a single line is malformed
		res.AddRegularError(err)
		return nil
	}

	return p.ParseRow(row, res, p.idx)
}
