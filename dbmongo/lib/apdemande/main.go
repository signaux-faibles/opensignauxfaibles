package apdemande

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
)

// Periode Période de temps avec un début et une fin

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string       `json:"id_demande" bson:"id_demande"`
	Siret              string       `json:"-" bson:"-"`
	EffectifEntreprise *int         `json:"effectif_entreprise" bson:"effectif_entreprise"`
	Effectif           *int         `json:"effectif" bson:"effectif"`
	DateStatut         time.Time    `json:"date_statut" bson:"date_statut"`
	Periode            misc.Periode `json:"periode" bson:"periode"`
	HTA                *float64     `json:"hta" bson:"hta"`
	MTA                *float64     `json:"mta" bson:"mta"`
	EffectifAutorise   *int         `json:"effectif_autorise" bson:"effectif_autorise"`
	MotifRecoursSE     *int         `json:"motif_recours_se" bson:"motif_recours_se"`
	HeureConsommee     *float64     `json:"heure_consommee" bson:"heure_consommee"`
	MontantConsomme    *float64     `json:"montant_consommee" bson:"montant_consommee"`
	EffectifConsomme   *int         `json:"effectif_consomme" bson:"effectif_consomme"`
}

// Key id de l'objet
func (apdemande APDemande) Key() string {
	return apdemande.Siret
}

// Type de données
func (apdemande APDemande) Type() string {
	return "apdemande"
}

// Scope de l'objet
func (apdemande APDemande) Scope() string {
	return "etablissement"
}

// Parser expose le parseur et le type de fichier qu'il supporte.
var Parser = marshal.Parser{FileType: "apdemande", FileParser: ParseFile}

// ParseFile permet de lancer le parsing du fichier demandé.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) marshal.OpenFileResult {
	var idx colMapping
	file, reader, err := openFile(filePath)
	if err == nil {
		idx, err = parseColMapping(reader)
	}
	return marshal.OpenFileResult{
		Error: err,
		ParseLines: func(parsedLineChan chan base.ParsedLineResult) {
			parseLines(reader, idx, parsedLineChan)
		},
		Close: func() {
			file.Close()
		},
	}
}

func openFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	return file, reader, nil
}

type colMapping map[string]int

func parseColMapping(reader *csv.Reader) (colMapping, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	var idx = colMapping{}
	for i, field := range header {
		idx[field] = i
	}
	fields := []string{
		"ID_DA",
		"ETAB_SIRET",
		"EFF_ENT",
		"EFF_ETAB",
		"DATE_STATUT",
		"DATE_DEB",
		"DATE_FIN",
		"HTA",
		"EFF_AUTO",
		"MOTIF_RECOURS_SE",
		"S_HEURE_CONSOM_TOT",
		"S_EFF_CONSOM_TOT",
	}
	for _, field := range fields {
		if _, found := idx[field]; !found {
			return nil, errors.New("Colonne " + field + " non trouvée. Abandon.")
		}
	}
	return idx, nil
}

func parseLines(reader *csv.Reader, idx colMapping, parsedLineChan chan base.ParsedLineResult) {
	for {
		parsedLine := base.ParsedLineResult{}
		row, err := reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(err)
		} else if row[idx["ETAB_SIRET"]] == "" {
			parsedLine.AddError(errors.New("invalidLine")) // TODO: retirer validation
		} else {
			parseApDemandeLine(row, idx, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []base.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseApDemandeLine(row []string, idx colMapping, parsedLine *base.ParsedLineResult) {
	apdemande := APDemande{}
	apdemande.ID = row[idx["ID_DA"]]
	apdemande.Siret = row[idx["ETAB_SIRET"]]
	var err error
	apdemande.EffectifEntreprise, err = misc.ParsePInt(row[idx["EFF_ENT"]])
	parsedLine.AddError(err)
	apdemande.Effectif, err = misc.ParsePInt(row[idx["EFF_ETAB"]])
	parsedLine.AddError(err)
	apdemande.DateStatut, err = time.Parse("02/01/2006", row[idx["DATE_STATUT"]])
	parsedLine.AddError(err)
	apdemande.Periode = misc.Periode{}
	apdemande.Periode.Start, err = time.Parse("02/01/2006", row[idx["DATE_DEB"]])
	parsedLine.AddError(err)
	apdemande.Periode.End, err = time.Parse("02/01/2006", row[idx["DATE_FIN"]])
	parsedLine.AddError(err)
	apdemande.HTA, err = misc.ParsePFloat(row[idx["HTA"]])
	parsedLine.AddError(err)
	apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(row[idx["MTA"]], ",", "."))
	parsedLine.AddError(err)
	apdemande.EffectifAutorise, err = misc.ParsePInt(row[idx["EFF_AUTO"]])
	parsedLine.AddError(err)
	apdemande.MotifRecoursSE, err = misc.ParsePInt(row[idx["MOTIF_RECOURS_SE"]])
	parsedLine.AddError(err)
	apdemande.HeureConsommee, err = misc.ParsePFloat(row[idx["S_HEURE_CONSOM_TOT"]])
	parsedLine.AddError(err)
	apdemande.EffectifConsomme, err = misc.ParsePInt(row[idx["S_EFF_CONSOM_TOT"]])
	parsedLine.AddError(err)
	apdemande.MontantConsomme, err = misc.ParsePFloat(strings.ReplaceAll(row[idx["S_MONTANT_CONSOM_TOT"]], ",", "."))
	parsedLine.AddError(err)
	parsedLine.AddTuple(apdemande)
}
