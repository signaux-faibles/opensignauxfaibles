package urssaf

import (
	"encoding/csv"
	"os"
	"strconv"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/sfregexp"
)

type effectifEntParser struct {
	file    *os.File
	reader  *csv.Reader
	periods []periodCol
	idx     engine.ColMapping
}

func NewParserEffectifEnt() *effectifEntParser {
	return &effectifEntParser{}
}

func (parser *effectifEntParser) Type() base.ParserType {
	return base.EffectifEnt
}

func (parser *effectifEntParser) Init(_ *engine.Cache, _ engine.SirenFilter, _ *base.AdminBatch) error {
	return nil
}

func (parser *effectifEntParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, parser.periods, err = parseEffectifColMapping(parser.reader, EffectifEnt{})
	}
	return err
}

func (parser *effectifEntParser) ReadNext(res *engine.ParsedLineResult) error {
	row, err := parser.reader.Read()
	if err != nil {
		return err
	}

	for _, period := range parser.periods {
		value := row[period.colIndex]

		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			s, err := strconv.ParseFloat(noThousandsSep, 64)
			res.AddRegularError(err)
			e := int(s)
			if e > 0 {
				idxRow := parser.idx.IndexRow(row)
				res.AddTuple(EffectifEnt{
					Siren:       idxRow.GetVal("siren"),
					Periode:     period.dateStart,
					EffectifEnt: e,
				})
			}
		}
	}
	return nil
}

func (parser *effectifEntParser) Close() error {
	return parser.file.Close()
}
