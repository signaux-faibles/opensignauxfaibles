package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
)

// Procol Procédures collectives, extraction URSSAF
type Procol struct {
	DateEffet    time.Time `json:"date_effet" bson:"date_effet"`
	ActionProcol string    `json:"action_procol" bson:"action_procol"`
	StadeProcol  string    `json:"stade_procol" bson:"stade_procol"`
	Siret        string    `json:"-" bson:"-"`
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
	idx    colMapping
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
	parser.file, parser.reader, err = openProcolFile(filePath)
	if err == nil {
		parser.idx, err = parseProcolColMapping(parser.reader)
	}
	return err
}

func openProcolFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true
	return file, reader, err
}

func parseProcolColMapping(reader *csv.Reader) (colMapping, error) {
	fields, err := reader.Read()
	if err != nil {
		return nil, err
	}
	var idx = colMapping{
		"dt_effet":      misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "dt_effet" }),
		"lib_actx_stdx": misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "lib_actx_stdx" }),
		"siret":         misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siret" }),
	}
	if misc.SliceMin(idx["dt_effet"], idx["lib_actx_stdx"], idx["siret"]) == -1 {
		return nil, errors.New("format de fichier incorrect")
	}
	return idx, nil
}

func (parser *procolParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(base.NewRegularError(err))
		} else {
			parseProcolLine(row, parser.idx, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []marshal.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseProcolLine(row []string, idx colMapping, parsedLine *marshal.ParsedLineResult) {
	var err error
	procol := Procol{}
	procol.DateEffet, err = time.Parse("02Jan2006", row[idx["dt_effet"]])
	parsedLine.AddError(base.NewRegularError(err))
	procol.Siret = row[idx["siret"]]
	actionStade := row[idx["lib_actx_stdx"]]
	splitted := strings.Split(strings.ToLower(actionStade), "_")
	for i, v := range splitted {
		r, err := regexp.Compile("liquidation|redressement|sauvegarde")
		parsedLine.AddError(base.NewRegularError(err))
		if match := r.MatchString(v); match {
			procol.ActionProcol = v
			procol.StadeProcol = strings.Join(append(splitted[:i], splitted[i+1:]...), "_")
			break
		}
	}
	parsedLine.AddTuple(procol)
}
