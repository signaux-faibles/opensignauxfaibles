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
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

	"github.com/signaux-faibles/gournal"
)

// EffectifEnt Urssaf
type EffectifEnt struct {
	Siren       string    `json:"-" bson:"-"`
	Periode     time.Time `json:"periode" bson:"periode"`
	EffectifEnt int       `json:"effectif" bson:"effectif"`
}

// Key _id de l'objet
func (effectifEnt EffectifEnt) Key() string {
	return effectifEnt.Siren
}

// Scope de l'objet
func (effectifEnt EffectifEnt) Scope() string {
	return "entreprise"
}

// Type de l'objet
func (effectifEnt EffectifEnt) Type() string {
	return "effectif_ent"
}

type periodCol struct {
	dateStart time.Time
	colIndex  int
}

// ParseEffectifPeriod extrait les périodes depuis une liste de noms de colonnes csv.
func parseEffectifPeriod(fields []string) []periodCol {
	periods := []periodCol{}
	re, _ := regexp.Compile("^eff")
	for index, field := range fields {
		if re.MatchString(field) {
			date, _ := marshal.UrssafToPeriod(field[3:9])
			periods = append(periods, periodCol{dateStart: date.Start, colIndex: index})
		}
	}
	return periods
}

// ParserEffectifEnt expose le parseur et le type de fichier qu'il supporte.
var ParserEffectifEnt = marshal.Parser{FileType: "effectif_ent", FileParser: ParseEffectifEntFile}

// ParseEffectifEntFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseEffectifEntFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	filter := marshal.GetSirenFilterFromCache(*cache)
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	parseEffectifEntFile(reader, filter, tracker, outputChannel)
}

func parseEffectifEntFile(reader *csv.Reader, filter marshal.SirenFilter, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	fields, err := reader.Read()
	if err != nil {
		tracker.Add(err)
		return
	}

	var idx = colMapping{
		"siren": misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siren" }),
	}

	// Dans quels champs lire l'effectifEnt
	periods := parseEffectifPeriod(fields)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
		} else {
			effectifs := parseEffectifEntLine(periods, row, idx, filter, tracker)
			for _, eff := range effectifs {
				outputChannel <- eff
			}
		}
		tracker.Next()
	}
}

func parseEffectifEntLine(periods []periodCol, row []string, idx colMapping, filter marshal.SirenFilter, tracker *gournal.Tracker) []EffectifEnt {
	var effectifs = []EffectifEnt{}
	siren := row[idx["siren"]]
	filtered, err := marshal.IsFiltered(siren, filter)
	tracker.Add(err)
	if len(siren) != 9 {
		tracker.Add(errors.New("Format de siren incorrect : " + siren))
	} else if !filtered {
		for _, period := range periods {
			value := row[period.colIndex]
			if value != "" {
				noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
				s, err := strconv.ParseFloat(noThousandsSep, 64)
				tracker.Add(err)
				e := int(s)
				if e > 0 {
					effectifs = append(effectifs, EffectifEnt{
						Siren:       siren,
						Periode:     period.dateStart,
						EffectifEnt: e,
					})
				}
			}
		}
	}
	return effectifs
}
