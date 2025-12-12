package urssaf

import (
	"io"
	"regexp"
	"slices"
	"strings"
	"time"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type ProcolParser struct{}

func NewProcolParser() engine.Parser {
	return &ProcolParser{}
}

func (parser *ProcolParser) Type() engine.ParserType { return engine.Procol }
func (parser *ProcolParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:        r,
		RowParser:     &procolRowParser{},
		Comma:         ';',
		LazyQuotes:    false,
		CaseSensitive: false,
		DestTuple:     Procol{},
	}
}

type procolRowParser struct{}

func (rp *procolRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {

	var err error
	idxRow := idx.IndexRow(row)
	procol := Procol{}
	procol.DateEffet, err = time.Parse("02Jan2006", idxRow.GetVal("Dt_effet"))
	res.AddRegularError(err)

	procol.Siret = idxRow.GetVal("Siret")

	actionStade := idxRow.GetVal("Lib_actx_stdx")
	splitted := strings.Split(strings.ToLower(actionStade), "_")

	for i, v := range splitted {
		r, err := regexp.Compile("liquidation|redressement|sauvegarde")
		res.AddRegularError(err)
		if match := r.MatchString(v); match {
			procol.ActionProcol = v
			procol.StadeProcol = strings.Join(slices.Delete(slices.Clone(splitted), i, i+1), "_")
			break
		}
	}
	if len(res.Errors) == 0 {
		res.AddTuple(procol)
	}
}
