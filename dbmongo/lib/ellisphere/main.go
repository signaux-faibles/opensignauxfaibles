package ellisphere

import (
	"github.com/pkg/errors"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"
	"github.com/tealeg/xlsx/v3"

	"github.com/signaux-faibles/gournal"
)

// Ellisphere informations groupe pour une entreprise
type Ellisphere struct {
	Siren               string  `json:"-" bson:"-" xlsx:"14"`
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

// Parser expose le parseur et le type de fichier qu'il supporte.
var Parser = marshal.Parser{FileType: "ellisphere", FileParser: ParseFile}

// ParseFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseFile(path string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker) marshal.LineParser {
	filter := marshal.GetSirenFilterFromCache(*cache)
	xlsxFile, err := xlsx.OpenFile(path)
	if err != nil {
		tracker.Add(err)
		return nil
	}
	if len(xlsxFile.Sheets) != 1 {
		tracker.Add(errors.Errorf("the source has %d sheets, should have only 1", len(xlsxFile.Sheets)))
		return nil
	}

	returned := false
	tuples := parseEllisphereSheet(xlsxFile.Sheets[0], filter, tracker)

	return func() []marshal.Tuple {
		if !returned {
			returned = true
			return tuples
		}
		return nil // EOF
	}
}

func parseEllisphereSheet(sheet *xlsx.Sheet, filter marshal.SirenFilter, tracker *gournal.Tracker) []marshal.Tuple {
	tuples := []marshal.Tuple{}
	sheet.ForEachRow(
		func(row *xlsx.Row) error {
			var ellisphere Ellisphere
			err := row.ReadStruct(&ellisphere)

			if !sfregexp.ValidSiren(ellisphere.Siren) {
				tracker.Add(errors.New("siren invalide : " + ellisphere.Siren))
			}
			if err == nil {
				tuples = append(tuples, ellisphere)
			}
			tracker.Add(err)
			tracker.Next()
			return nil
		},
	)
	return tuples
}
