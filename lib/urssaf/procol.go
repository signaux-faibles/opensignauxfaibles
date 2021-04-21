package urssaf

import (
	"encoding/csv"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

// Procol Procédures collectives, extraction URSSAF
type Procol struct {
	DateEffet    time.Time `col:"dt_effet"      json:"date_effet"    bson:"date_effet"`
	ActionProcol string    `col:"lib_actx_stdx" json:"action_procol" bson:"action_procol"`
	StadeProcol  string    `col:"lib_actx_stdx" json:"stade_procol"  bson:"stade_procol"`
	Siret        string    `col:"siret"         json:"-"             bson:"-"`
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
func (procol Procol) Type() string {
	return "procol"
}

// ParserProcol fournit une instance utilisable par ParseFilesFromBatch.
var ParserProcol = &procolParser{}

type procolParser struct {
	file   *os.File
	reader *csv.Reader
	idx    marshal.ColMapping
}

func (parser *procolParser) GetFileType() string {
	return "procol"
}

func (parser *procolParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *procolParser) Close() error {
	return parser.file.Close()
}

func (parser *procolParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ';', true)
	if err == nil {
		parser.idx, err = parseProcolColMapping(parser.reader)
	}
	return err
}

func parseProcolColMapping(reader *csv.Reader) (marshal.ColMapping, error) {
	fields, err := reader.Read()
	if err != nil {
		return marshal.ColMapping{}, err
	}
	return marshal.ValidateAndIndexColumnsFromColTags(marshal.LowercaseFields(fields), Procol{})
}

func (parser *procolParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseProcolLine(row, parser.idx, parsedLine)
	})
}

func parseProcolLine(row []string, idx marshal.ColMapping, parsedLine *marshal.ParsedLineResult) {
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
	parsedLine.AddTuple(procol)
}
