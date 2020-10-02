package apdemande

import (
	"github.com/pkg/errors"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/tealeg/xlsx/v3"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// Ellisphere informations groupe pour une entreprise
type Ellisphere struct {
	Siren               string  `json:"siren" bson:"-" xlsx:"14"`
	CodeGroupe          string  `json:"code_groupe,omitempty" bson:"code_groupe,omitempty" xlsx:"0"`
	SirenGroupe         string  `json:"siren_groupe,omitempty" bson:"siren_groupe,omitempty" xlsx:"2"`
	RefIDGroupe         string  `json:"refid_groupe,omitempty" bson:"refid_groupe,omitempty" xlsx:"3"`
	RaisocGroupe        string  `json:"raison_sociale_groupe,omitempty" bson:"raison_sociale_groupe,omitempty" xlsx:"4"`
	AdresseGroupe       string  `json:"adresse_groupe,omitempty" bson:"adresse_groupe,omitempty" xlsx:"5"`
	PersonnePouMGroupe  string  `json:"personne_pou_m_groupe,omitempty" bson:"personne_pou_m_groupe,omitempty" xlsx:"1"`
	NiveauDetention     int     `json:"niveau_detention,omitempty" bson:"niveau_detention,omitempty" xlsx:"9"`
	PartFinanciere      float64 `json:"part_financiere,omitempty" bson:"part_financiere,omitempty" xlsx:"10"`
	CodeFiliere         string  `json:"code_filiere,omitempty" bson:"code_filiere,omitempty" xlsx:"12"`
	RefIDFiliere        string  `json:"refid_filiere,omitempty" bson:"refid_filiere,omitempty" xlsx:"15"`
	PersonnePouMFiliere string  `json:"personne_pou_m_filiere,omitempty" bson:"personne_pou_m_filiere,omitempty" xlsx:"13"`
}

// Key id de l'objet
func (ellisphere Ellisphere) Key() string {
	return ellisphere.Siren
}

// Type de données
func (ellisphere Ellisphere) Type() string {
	return "ellisphere"
}

// Scope de l'objet
func (ellisphere Ellisphere) Scope() string {
	return "entreprise"
}

// Parser produit les lignes
func Parser(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)
	event := marshal.Event{
		Code:    "parserEllisphere",
		Channel: eventChannel,
	}

	go func() {
		defer close(outputChannel)
		defer close(eventChannel)

		for _, path := range batch.Files["ellisphere"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				engine.TrackerReports)

			xlsxFile, err := xlsx.OpenFile(viper.GetString("APP_DATA") + path)
			tracker.Error(err)

			if len(xlsxFile.Sheets) != 1 {
				tracker.Error(errors.Errorf("the source has %d sheets, should have only 1", len(xlsxFile.Sheets)))
				continue
			}
			sheet := xlsxFile.Sheets[0]
			err = sheet.ForEachRow(
				func(row *xlsx.Row) error {
					var ellisphere Ellisphere
					err := row.ReadStruct(&ellisphere)
					tracker.Error(err)
					outputChannel <- ellisphere
					return nil
				},
			)

			if err != nil {
				tracker.Error(err)
			}

			tracker.Next()

			event.Info(tracker.Report("abstract"))
		}
	}()

	return outputChannel, eventChannel
}

// file, err := os.Open(viper.GetString("APP_DATA") + path)
// if err != nil {
// 	event.Critical(path + ": erreur à l'ouverture du fichier: " + err.Error())
// 	return
// }
// reader := csv.NewReader(bufio.NewReader(file))
// reader.Comma = ','
// reader.LazyQuotes = true

// event.Info(path + ": ouverture")

// header, err := reader.Read()
// if err != nil {
// 	tracker.Error(err)
// 	event.Debug(tracker.Report("invalidLine"))
// 	break
// }

// f := make(map[string]int)
// for idx, field := range header {
// 	f[field] = idx
// }
// fields := []string{
// 	"ID_DA",
// 	"ETAB_SIRET",
// 	"EFF_ENT",
// 	"EFF_ETAB",
// 	"DATE_STATUT",
// 	"DATE_DEB",
// 	"DATE_FIN",
// 	"HTA",
// 	"EFF_AUTO",
// 	"MOTIF_RECOURS_SE",
// 	"S_HEURE_CONSOM_TOT",
// 	"S_EFF_CONSOM_TOT",
// }

// for _, field := range fields {
// 	if _, found := f[field]; !found {
// 		event.Critical("Import du fichier " + path + ". " + field + " non trouvé. Abandon.")
// 		continue
// 	}
// }
// for {
// 	row, err := reader.Read()
// 	if err == io.EOF {
// 		break
// 	} else if err != nil {
// 		tracker.Error(err)
// 		event.Debug(tracker.Report("invalidLine"))
// 		break
// 	}

// 	if row[f["ETAB_SIRET"]] != "" {
// 		apdemande := APDemande{}
// 		apdemande.ID = row[f["ID_DA"]]
// 		apdemande.Siret = row[f["ETAB_SIRET"]]
// 		apdemande.EffectifEntreprise, err = misc.ParsePInt(row[f["EFF_ENT"]])
// 		tracker.Error(err)
// 		apdemande.Effectif, err = misc.ParsePInt(row[f["EFF_ETAB"]])
// 		tracker.Error(err)
// 		apdemande.DateStatut, err = time.Parse("02/01/2006", row[f["DATE_STATUT"]])
// 		tracker.Error(err)
// 		apdemande.Periode = misc.Periode{}
// 		apdemande.Periode.Start, err = time.Parse("02/01/2006", row[f["DATE_DEB"]])
// 		tracker.Error(err)
// 		apdemande.Periode.End, err = time.Parse("02/01/2006", row[f["DATE_FIN"]])
// 		tracker.Error(err)
// 		apdemande.HTA, err = misc.ParsePFloat(row[f["HTA"]])
// 		tracker.Error(err)
// 		apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(row[f["MTA"]], ",", "."))
// 		tracker.Error(err)
// 		apdemande.EffectifAutorise, err = misc.ParsePInt(row[f["EFF_AUTO"]])
// 		tracker.Error(err)
// 		apdemande.MotifRecoursSE, err = misc.ParsePInt(row[f["MOTIF_RECOURS_SE"]])
// 		tracker.Error(err)
// 		apdemande.HeureConsommee, err = misc.ParsePFloat(row[f["S_HEURE_CONSOM_TOT"]])
// 		tracker.Error(err)
// 		apdemande.EffectifConsomme, err = misc.ParsePInt(row[f["S_EFF_CONSOM_TOT"]])
// 		tracker.Error(err)
// 		apdemande.MontantConsomme, err = misc.ParsePFloat(strings.ReplaceAll(row[f["S_MONTANT_CONSOM_TOT"]], ",", "."))
// 		tracker.Error(err)
