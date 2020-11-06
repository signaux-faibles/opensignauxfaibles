package ellisphere

import (
	"github.com/pkg/errors"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/tealeg/xlsx/v3"
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

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &ellisphereParser{}

type ellisphereParser struct {
	sheet *xlsx.Sheet
}

func (parser *ellisphereParser) GetFileType() string {
	return "ellisphere"
}

func (parser *ellisphereParser) Init(cache *marshal.Cache, batch *base.AdminBatch) {}

func (parser *ellisphereParser) Open(filePath string) (err error) {
	xlsxFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return err
	}
	if len(xlsxFile.Sheets) != 1 {
		return errors.Errorf("the source has %d sheets, should have only 1", len(xlsxFile.Sheets))
	}
	parser.sheet = xlsxFile.Sheets[0]
	return nil
}

func (parser *ellisphereParser) Close() error {
	return nil
}

func (parser *ellisphereParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	parser.sheet.ForEachRow(
		func(row *xlsx.Row) error {
			parsedLine := marshal.ParsedLineResult{}
			var ellisphere Ellisphere
			err := row.ReadStruct(&ellisphere)
			parsedLine.AddError(base.NewRegularError(err))
			if len(parsedLine.Errors) == 0 {
				parsedLine.AddTuple(ellisphere)
			}
			parsedLineChan <- parsedLine
			return nil
		},
	)
	close(parsedLineChan)
}
