package urssaf

import (
	"encoding/csv"
	"os"
	"regexp"
	"strings"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
)

type procolParser struct {
	file   *os.File
	reader *csv.Reader
	idx    engine.ColMapping
}

func NewParserProcol() *procolParser {
	return &procolParser{}
}

func (parser *procolParser) Type() base.ParserType {
	return base.Procol
}

func (parser *procolParser) Init(_ *engine.Cache, _ engine.SirenFilter, _ *base.AdminBatch) error {
	return nil
}

func (parser *procolParser) Close() error {
	return parser.file.Close()
}

func (parser *procolParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ';', true)
	if err == nil {
		parser.idx, err = parseProcolColMapping(parser.reader)
	}
	return err
}

func parseProcolColMapping(reader *csv.Reader) (engine.ColMapping, error) {
	fields, err := reader.Read()
	if err != nil {
		return engine.ColMapping{}, err
	}
	return engine.ValidateAndIndexColumnsFromInputTags(engine.LowercaseFields(fields), Procol{})
}

func (parser *procolParser) ReadNext(res *engine.ParsedLineResult) error {
	row, err := parser.reader.Read()
	if err != nil {
		return err
	}

	idxRow := parser.idx.IndexRow(row)
	procol := Procol{}
	procol.DateEffet, err = time.Parse("02Jan2006", idxRow.GetVal("dt_effet"))
	res.AddRegularError(err)
	procol.Siret = idxRow.GetVal("siret")
	actionStade := idxRow.GetVal("lib_actx_stdx")
	splitted := strings.Split(strings.ToLower(actionStade), "_")
	for i, v := range splitted {
		r, err := regexp.Compile("liquidation|redressement|sauvegarde")
		res.AddRegularError(err)
		if match := r.MatchString(v); match {
			procol.ActionProcol = v
			procol.StadeProcol = strings.Join(append(splitted[:i], splitted[i+1:]...), "_")
			break
		}
	}
	if len(res.Errors) == 0 {
		res.AddTuple(procol)
	}
	return nil
}
