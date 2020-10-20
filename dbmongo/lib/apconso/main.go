package apconso

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
)

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `json:"id_conso"         bson:"id_conso"`
	Siret          string    `json:"-"                bson:"-"`
	HeureConsommee *float64  `json:"heure_consomme"   bson:"heure_consomme"`
	Montant        *float64  `json:"montant"          bson:"montant"`
	Effectif       *int      `json:"effectif"         bson:"effectif"`
	Periode        time.Time `json:"periode"          bson:"periode"`
}

// Key id de l'objet
func (apconso APConso) Key() string {
	return apconso.Siret
}

// Type de données
func (apconso APConso) Type() string {
	return "apconso"
}

// Scope de l'objet
func (apconso APConso) Scope() string {
	return "etablissement"
}

// Parser produit des lignes de consommation d'activité partielle
func Parser(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	return marshal.ParseFilesFromBatch(cache, batch, marshal.Parser{FileType: "apconso", FileParser: ParseFile})
}

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
	parseApConsoFile(reader, tracker, outputChannel)
}

type colMapping map[string]int

func parseApConsoFile(reader *csv.Reader, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	header, err := reader.Read()
	if err != nil {
		tracker.Add(err)
		return
	}
	var idx = colMapping{}
	idx["ID"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "ID_DA" })
	idx["Siret"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "ETAB_SIRET" })
	idx["Periode"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "MOIS" })
	idx["HeureConsommee"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "HEURES" })
	idx["Montant"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "MONTANTS" })
	idx["Effectif"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "EFFECTIFS" })

	if misc.SliceMin(idx["ID"], idx["Siret"], idx["Periode"], idx["HeureConsommee"], idx["Montant"], idx["Effectif"]) == -1 {
		tracker.Add(errors.New("entête non conforme, fichier ignoré"))
		return
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
		} else if len(row) > 0 {
			// TODO: filtrer et/ou valider siret ?
			apconso := parseApConsoLine(row, tracker, idx)
			if !tracker.HasErrorInCurrentCycle() && apconso.Siret != "" {
				outputChannel <- apconso
			}
		}
		tracker.Next()
	}
}

func parseApConsoLine(row []string, tracker *gournal.Tracker, idx colMapping) APConso {
	apconso := APConso{}
	apconso.ID = row[idx["ID"]]
	apconso.Siret = row[idx["Siret"]]
	var err error
	apconso.Periode, err = time.Parse("01/2006", row[idx["Periode"]])
	tracker.Add(err)
	apconso.HeureConsommee, err = misc.ParsePFloat(row[idx["HeureConsommee"]])
	tracker.Add(err)
	apconso.Montant, err = misc.ParsePFloat(row[idx["Montant"]])
	tracker.Add(err)
	apconso.Effectif, err = misc.ParsePInt(row[idx["Effectif"]])
	tracker.Add(err)
	return apconso
}
