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

// Procol Procédures collectives, extraction URSSAF
type Procol struct {
	Siret        string    `input:"siret"         json:"-"             csv:"siret"`
	DateEffet    time.Time `input:"dt_effet"      json:"date_effet"    csv:"date_effet"`
	ActionProcol string    `input:"lib_actx_stdx" json:"action_procol" csv:"action_procol"`
	StadeProcol  string    `input:"lib_actx_stdx" json:"stade_procol"  csv:"stade_procol"`
}

// Key _id de l'objet
func (procol Procol) Key() string {
	return procol.Siret
}

// Scope de l'objet
func (procol Procol) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (procol Procol) Type() base.ParserType {
	return base.Procol
}

// ParserProcol fournit une instance utilisable par ParseFilesFromBatch.
var ParserProcol = &procolParser{}

type procolParser struct {
	file   *os.File
	reader *csv.Reader
	idx    engine.ColMapping
}

func (parser *procolParser) Type() base.ParserType {
	return base.Procol
}

func (parser *procolParser) Init(_ *engine.Cache, _ *base.AdminBatch) error {
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

func (parser *procolParser) ParseLines(parsedLineChan chan engine.ParsedLineResult) {
	engine.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *engine.ParsedLineResult) {
		parseProcolLine(row, parser.idx, parsedLine)
	})
}

func parseProcolLine(row []string, idx engine.ColMapping, parsedLine *engine.ParsedLineResult) {
	var err error
	idxRow := idx.IndexRow(row)
	procol := Procol{}
	procol.DateEffet, err = time.Parse("02Jan2006", idxRow.GetVal("dt_effet"))
	parsedLine.AddRegularError(err)
	procol.Siret = idxRow.GetVal("siret")
	actionStade := idxRow.GetVal("lib_actx_stdx")
	splitted := strings.Split(strings.ToLower(actionStade), "_")
	for i, v := range splitted {
		r, err := regexp.Compile("liquidation|redressement|sauvegarde")
		parsedLine.AddRegularError(err)
		if match := r.MatchString(v); match {
			procol.ActionProcol = v
			procol.StadeProcol = strings.Join(append(splitted[:i], splitted[i+1:]...), "_")
			break
		}
	}
	if len(parsedLine.Errors) == 0 {
		parsedLine.AddTuple(procol)
	}
}
