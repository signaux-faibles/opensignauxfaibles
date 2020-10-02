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
			sheet.ForEachRow(
				func(row *xlsx.Row) error {
					var ellisphere Ellisphere
					err := row.ReadStruct(&ellisphere)
					if err == nil {
						outputChannel <- ellisphere
					}
					tracker.Error(err)
					tracker.Next()
					return nil
				},
			)

			event.Info(tracker.Report("abstract"))
		}
	}()

	return outputChannel, eventChannel
}
