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

type colMapping map[string]int

// ParseFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) (marshal.ParsedLineChan, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// defer file.Close() // TODO: à réactiver
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

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
			} else if row[idx["ETAB_SIRET"]] == "" {
				parsedLine.AddError(errors.New("invalidLine")) // TODO: retirer validation
			} else {
				parseApDemandeLine(row, idx, &parsedLine)
				if len(parsedLine.Errors) > 0 {
					parsedLine.Tuples = []marshal.Tuple{}
				}
			}
			parsedLineChan <- parsedLine
		}
	}()
	return parsedLineChan, nil
}

func parseApDemandeLine(row []string, idx colMapping, parsedLine *marshal.ParsedLineResult) {
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
