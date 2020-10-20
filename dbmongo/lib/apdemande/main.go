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

	"github.com/signaux-faibles/gournal"
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

// ParseFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	parseApDemandeFile(reader, tracker, outputChannel)
}

type colMapping map[string]int

func parseApDemandeFile(reader *csv.Reader, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	header, err := reader.Read()
	if err != nil {
		tracker.Add(err)
		return
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
			tracker.Add(errors.New("Colonne " + field + " non trouvée. Abandon."))
			return
		}
	}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
		} else if row[idx["ETAB_SIRET"]] == "" {
			tracker.Add(errors.New("invalidLine"))
		} else {
			// TODO: filtrer et/ou valider siret ?
			apdemande := parseApDemandeLine(row, tracker, idx)
			if !tracker.HasErrorInCurrentCycle() {
				outputChannel <- apdemande
			}
		}
		tracker.Next()
	}
}

func parseApDemandeLine(row []string, tracker *gournal.Tracker, idx colMapping) APDemande {
	apdemande := APDemande{}
	apdemande.ID = row[idx["ID_DA"]]
	apdemande.Siret = row[idx["ETAB_SIRET"]]
	var err error
	apdemande.EffectifEntreprise, err = misc.ParsePInt(row[idx["EFF_ENT"]])
	tracker.Add(err)
	apdemande.Effectif, err = misc.ParsePInt(row[idx["EFF_ETAB"]])
	tracker.Add(err)
	apdemande.DateStatut, err = time.Parse("02/01/2006", row[idx["DATE_STATUT"]])
	tracker.Add(err)
	apdemande.Periode = misc.Periode{}
	apdemande.Periode.Start, err = time.Parse("02/01/2006", row[idx["DATE_DEB"]])
	tracker.Add(err)
	apdemande.Periode.End, err = time.Parse("02/01/2006", row[idx["DATE_FIN"]])
	tracker.Add(err)
	apdemande.HTA, err = misc.ParsePFloat(row[idx["HTA"]])
	tracker.Add(err)
	apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(row[idx["MTA"]], ",", "."))
	tracker.Add(err)
	apdemande.EffectifAutorise, err = misc.ParsePInt(row[idx["EFF_AUTO"]])
	tracker.Add(err)
	apdemande.MotifRecoursSE, err = misc.ParsePInt(row[idx["MOTIF_RECOURS_SE"]])
	tracker.Add(err)
	apdemande.HeureConsommee, err = misc.ParsePFloat(row[idx["S_HEURE_CONSOM_TOT"]])
	tracker.Add(err)
	apdemande.EffectifConsomme, err = misc.ParsePInt(row[idx["S_EFF_CONSOM_TOT"]])
	tracker.Add(err)
	apdemande.MontantConsomme, err = misc.ParsePFloat(strings.ReplaceAll(row[idx["S_MONTANT_CONSOM_TOT"]], ",", "."))
	tracker.Add(err)
	return apdemande
}
