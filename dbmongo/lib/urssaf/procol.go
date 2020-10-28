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

	"github.com/signaux-faibles/gournal"
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
func ParseProcolFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker) marshal.LineParser {
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return nil
	}
	// defer file.Close() // TODO: à réactiver
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true

	fields, err := reader.Read()
	if err != nil {
		tracker.Add(err)
		return nil
	}

	var idx = colMapping{
		"dt_effet":      misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "dt_effet" }),
		"lib_actx_stdx": misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "lib_actx_stdx" }),
		"siret":         misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siret" }),
	}

	if misc.SliceMin(idx["dt_effet"], idx["lib_actx_stdx"], idx["siret"]) == -1 {
		tracker.Add(errors.New("format de fichier incorrect"))
		return nil
	}

	return func() []marshal.Tuple {
		row, err := reader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			tracker.Add(err)
		} else {
			procol := parseProcolLine(row, tracker, idx)
			if _, err := strconv.Atoi(row[idx["siret"]]); err == nil && len(row[idx["siret"]]) == 14 {
				if !tracker.HasErrorInCurrentCycle() {
					return []marshal.Tuple{procol}
				}
			}
		}
		return []marshal.Tuple{}
	}
}

func parseProcolLine(row []string, tracker *gournal.Tracker, idx colMapping) Procol {
	var err error
	procol := Procol{}
	procol.DateEffet, err = time.Parse("02Jan2006", row[idx["dt_effet"]])
	tracker.Add(err)
	procol.Siret = row[idx["siret"]]
	actionStade := row[idx["lib_actx_stdx"]]
	splitted := strings.Split(strings.ToLower(actionStade), "_")
	for i, v := range splitted {
		r, err := regexp.Compile("liquidation|redressement|sauvegarde")
		tracker.Add(err)
		if match := r.MatchString(v); match {
			procol.ActionProcol = v
			procol.StadeProcol = strings.Join(append(splitted[:i], splitted[i+1:]...), "_")
			break
		}
	}
	return (procol)
}
