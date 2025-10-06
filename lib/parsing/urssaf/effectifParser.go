package urssaf

import (
	"encoding/csv"
	"os"
	"strconv"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/sfregexp"
)

type effectifParser struct {
	file    *os.File
	reader  *csv.Reader
	periods []periodCol
	idx     engine.ColMapping
}

func NewParserEffectif() *effectifParser {
	return &effectifParser{}
}

func (parser *effectifParser) Type() base.ParserType {
	return base.Effectif
}

func (parser *effectifParser) Init(cache *engine.Cache, _ engine.SirenFilter, batch *base.AdminBatch) error {
	return nil
}

func (parser *effectifParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, parser.periods, err = parseEffectifColMapping(parser.reader, Effectif{})
	}
	return err
}

func (parser *effectifParser) ReadNext(res *engine.ParsedLineResult) error {
	row, err := parser.reader.Read()

	if err != nil {
		return err
	}

	for _, period := range parser.periods {
		value := row[period.colIndex]

		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			e, err := strconv.Atoi(noThousandsSep)
			res.AddRegularError(err)
			if e > 0 {
				idxRow := parser.idx.IndexRow(row)
				res.AddTuple(Effectif{
					Siret:        idxRow.GetVal("siret"),
					NumeroCompte: idxRow.GetVal("compte"),
					Periode:      period.dateStart,
					Effectif:     e,
				})
			}
		}
	}
	return nil
}

func (parser *effectifParser) Close() error {
	return parser.file.Close()
}
