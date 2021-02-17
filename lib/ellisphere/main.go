package ellisphere

import (
	"github.com/pkg/errors"
	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/tealeg/xlsx/v3"
)

// Colonnes:
//  0: GRP Code
//  1: GRP Personne PouM
//  2: GRP SIREN 9
//  3: GRP RefID
//  4: GRP Raison Sociale
//  5: GRP Adresse
//  6: GRP Code Postal
//  7: GRP Ville
//  8: GRP Pays
//  9: Niveau de détention
// 10: % Financier
// 11: Tranche % Financier
// 12: FIL Code
// 13: FIL Personne PouM
// 14: FIL SIREN 9
// 15: FIL RefID
// 16: FIL Raison Sociale
// 17: FIL Adresse
// 18: FIL Code Postal
// 19: FIL Ville
// 20: FIL Pays

// Ellisphere informations groupe pour une entreprise
type Ellisphere struct {
	Siren               string  `xlsx:"14" json:"-"                                bson:"-"`
	CodeGroupe          string  `xlsx:"0"  json:"code_groupe,omitempty"            bson:"code_groupe,omitempty"`
	SirenGroupe         string  `xlsx:"2"  json:"siren_groupe,omitempty"           bson:"siren_groupe,omitempty"`
	RefIDGroupe         string  `xlsx:"3"  json:"refid_groupe,omitempty"           bson:"refid_groupe,omitempty"`
	RaisocGroupe        string  `xlsx:"4"  json:"raison_sociale_groupe,omitempty"  bson:"raison_sociale_groupe,omitempty"`
	AdresseGroupe       string  `xlsx:"5"  json:"adresse_groupe,omitempty"         bson:"adresse_groupe,omitempty"`
	PersonnePouMGroupe  string  `xlsx:"1"  json:"personne_pou_m_groupe,omitempty"  bson:"personne_pou_m_groupe,omitempty"`
	NiveauDetention     int     `xlsx:"9"  json:"niveau_detention,omitempty"       bson:"niveau_detention,omitempty"`
	PartFinanciere      float64 `xlsx:"10" json:"part_financiere,omitempty"        bson:"part_financiere,omitempty"`
	CodeFiliere         string  `xlsx:"12" json:"code_filiere,omitempty"           bson:"code_filiere,omitempty"`
	RefIDFiliere        string  `xlsx:"15" json:"refid_filiere,omitempty"          bson:"refid_filiere,omitempty"`
	PersonnePouMFiliere string  `xlsx:"13" json:"personne_pou_m_filiere,omitempty" bson:"personne_pou_m_filiere,omitempty"`
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

func (parser *ellisphereParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

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
			parsedLine.AddRegularError(err)
			if len(parsedLine.Errors) == 0 {
				parsedLine.AddTuple(ellisphere)
			}
			parsedLineChan <- parsedLine
			return nil
		},
	)
	close(parsedLineChan)
}
