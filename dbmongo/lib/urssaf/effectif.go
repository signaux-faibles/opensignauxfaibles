package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

	"github.com/signaux-faibles/gournal"
)

// Effectif Urssaf
type Effectif struct {
	Siret        string    `json:"-" bson:"-"`
	NumeroCompte string    `json:"numero_compte" bson:"numero_compte"`
	Periode      time.Time `json:"periode" bson:"periode"`
	Effectif     int       `json:"effectif" bson:"effectif"`
}

// Key _id de l'objet
func (effectif Effectif) Key() string {
	return effectif.Siret
}

// Scope de l'objet
func (effectif Effectif) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (effectif Effectif) Type() string {
	return "effectif"
}

// ParserEffectif retourne un channel fournissant des données extraites
func ParserEffectif(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	return marshal.ParseFilesFromBatch(cache, batch, marshal.Parser{FileType: "effectif", FileParser: ParseEffectifFile})
}

// ParseEffectifFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseEffectifFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	filter := marshal.GetSirenFilterFromCache(*cache)
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	parseEffectifFile(reader, filter, tracker, outputChannel)
}

func parseEffectifFile(reader *csv.Reader, filter map[string]bool, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	fields, err := reader.Read()
	if err != nil {
		tracker.Add(err)
		return
	}

	var idx = colMapping{
		"siret":  misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siret" }),
		"compte": misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "compte" }),
	}

	if misc.SliceMin(idx["siret"], idx["compte"]) == -1 {
		tracker.Add(errors.New("erreur à l'analyse du fichier, abandon, l'un " +
			"des champs obligatoires n'a pu etre trouve:" +
			" siretIndex = " + strconv.Itoa(idx["siret"]) +
			", compteIndex = " + strconv.Itoa(idx["compte"])))
		return
	}

	// Dans quels champs lire l'effectif
	periods := parseEffectifPeriod(fields)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
		} else {
			effectifs := parseEffectifLine(periods, row, idx, filter, tracker)
			for _, eff := range effectifs {
				outputChannel <- eff
			}
		}
		tracker.Next()
	}
}

func parseEffectifLine(periods []periodCol, row []string, idx colMapping, filter map[string]bool, tracker *gournal.Tracker) []Effectif {
	var effectifs = []Effectif{}
	siret := row[idx["siret"]]
	validSiret := sfregexp.RegexpDict["siret"].MatchString(siret)
	if !validSiret {
		tracker.Add(base.NewRegularError(errors.New("Le siret/siren est invalide")))
	} else if filter != nil || !marshal.FilterHas(siret, filter) {
		for _, period := range periods {
			value := row[period.colIndex]
			if value != "" {
				noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
				e, err := strconv.Atoi(noThousandsSep)
				tracker.Add(err)
				if e > 0 {
					effectifs = append(effectifs, Effectif{
						Siret:        siret,
						NumeroCompte: row[idx["compte"]],
						Periode:      period.dateStart,
						Effectif:     e,
					})
				}
			}
		}
	}
	return effectifs
}
