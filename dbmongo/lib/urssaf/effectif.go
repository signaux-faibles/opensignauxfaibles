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

// ParserEffectif expose le parseur et le type de fichier qu'il supporte.
var ParserEffectif = marshal.Parser{FileType: "effectif", FileParser: ParseEffectifFile}

// ParseEffectifFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseEffectifFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) (marshal.ParsedLineChan, error) {
	filter := marshal.GetSirenFilterFromCache(*cache)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// defer file.Close() // TODO: à réactiver
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	fields, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var idx = colMapping{
		"siret":  misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siret" }),
		"compte": misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "compte" }),
	}

	if misc.SliceMin(idx["siret"], idx["compte"]) == -1 {
		return nil, errors.New("erreur à l'analyse du fichier, abandon, l'un " +
			"des champs obligatoires n'a pu etre trouve:" +
			" siretIndex = " + strconv.Itoa(idx["siret"]) +
			", compteIndex = " + strconv.Itoa(idx["compte"]))
	}

	// Dans quels champs lire l'effectif
	periods := parseEffectifPeriod(fields)

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
				parseEffectifLine(periods, row, idx, filter, &parsedLine)
			}
			parsedLineChan <- parsedLine
		}
	}()
	return parsedLineChan, nil
}

func parseEffectifLine(periods []periodCol, row []string, idx colMapping, filter marshal.SirenFilter, parsedLine *marshal.ParsedLineResult) {
	siret := row[idx["siret"]]
	validSiret := sfregexp.RegexpDict["siret"].MatchString(siret) // TODO: remove validation
	if !validSiret {
		parsedLine.AddError(base.NewRegularError(errors.New("Le siret/siren est invalide")))
	} else {
		for _, period := range periods {
			value := row[period.colIndex]
			if value != "" {
				noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
				e, err := strconv.Atoi(noThousandsSep)
				parsedLine.AddError(err)
				if e > 0 {
					parsedLine.AddTuple(Effectif{
						Siret:        siret,
						NumeroCompte: row[idx["compte"]],
						Periode:      period.dateStart,
						Effectif:     e,
					})
				}
			}
		}
	}
}
