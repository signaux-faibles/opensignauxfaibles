package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
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

// ParserProcol expose le parseur et le type de fichier qu'il supporte.
var ParserProcol = marshal.Parser{FileType: "procol", FileParser: ParseProcolFile}

// ParseProcolFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseProcolFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) (marshal.ParsedLineChan, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// defer file.Close() // TODO: à réactiver
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true

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

	parsedLineChan := make(marshal.ParsedLineChan)
	go func() {
		for {
			parsedLine := marshal.ParsedLineResult{}
			row, err := reader.Read()
			if err == io.EOF {
				close(parsedLineChan)
				break
			} else if err != nil {
				parsedLine.AddError(err)
			} else {
				parseProcolLine(row, idx, &parsedLine)
				if _, err := strconv.Atoi(row[idx["siret"]]); err == nil && len(row[idx["siret"]]) == 14 { // TODO: remove validation
					if len(parsedLine.Errors) > 0 {
						parsedLine.Tuples = []marshal.Tuple{}
					}
				}
			}
			parsedLineChan <- parsedLine
		}
	}()
	return parsedLineChan, nil
}

func parseProcolLine(row []string, idx colMapping, parsedLine *marshal.ParsedLineResult) {
	var err error
	procol := Procol{}
	procol.DateEffet, err = time.Parse("02Jan2006", row[idx["dt_effet"]])
	parsedLine.AddError(err)
	procol.Siret = row[idx["siret"]]
	actionStade := row[idx["lib_actx_stdx"]]
	splitted := strings.Split(strings.ToLower(actionStade), "_")
	for i, v := range splitted {
		r, err := regexp.Compile("liquidation|redressement|sauvegarde")
		parsedLine.AddError(err)
		if match := r.MatchString(v); match {
			procol.ActionProcol = v
			procol.StadeProcol = strings.Join(append(splitted[:i], splitted[i+1:]...), "_")
			break
		}
	}
	parsedLine.AddTuple(procol)
}
