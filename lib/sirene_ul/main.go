package sireneul

import (
	"encoding/csv"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
)

// SireneUL informations sur les entreprises
type SireneUL struct {
	Siren               string     ` input:"siren"                         json:"siren,omitempty"                  sql:"siren"                    csv:"Siren"`
	Nic                 string     ` json:"nic,omitempty"                                                            sql:"nic"                      `
	RaisonSociale       string     ` input:"denominationUniteLegale"       json:"raison_sociale"                   sql:"raison_sociale"           csv:"RaisonSociale"`
	Prenom1UniteLegale  string     ` input:"prenom1UniteLegale"            json:"prenom1_unite_legale,omitempty"   sql:"prenom1_unite_legale"     csv:"Prenom1UniteLegale"`
	Prenom2UniteLegale  string     ` input:"prenom2UniteLegale"            json:"prenom2_unite_legale,omitempty"   sql:"prenom2_unite_legale"     csv:"Prenom2UniteLegale"`
	Prenom3UniteLegale  string     ` input:"prenom3UniteLegale"            json:"prenom3_unite_legale,omitempty"   sql:"prenom3_unite_legale"     csv:"Prenom3UniteLegale"`
	Prenom4UniteLegale  string     ` input:"prenom4UniteLegale"            json:"prenom4_unite_legale,omitempty"   sql:"prenom4_unite_legale"     csv:"Prenom4UniteLegale"`
	NomUniteLegale      string     ` input:"nomUniteLegale"                json:"nom_unite_legale,omitempty"       sql:"nom_unite_legale"         csv:"NomUniteLegale"`
	NomUsageUniteLegale string     ` input:"nomUsageUniteLegale"           json:"nom_usage_unite_legale,omitempty" sql:"nom_usage_unite_legale"   csv:"NomUsageUniteLegale"`
	CodeStatutJuridique string     ` input:"categorieJuridiqueUniteLegale" json:"statut_juridique"                 sql:"statut_juridique"         csv:"CodeStatutJuridique"`
	Creation            *time.Time ` input:"dateCreationUniteLegale"       json:"date_creation,omitempty"          sql:"creation"                 csv:"Creation"`
}

// Key id de l'objet
func (sireneUL SireneUL) Key() string {
	return sireneUL.Siren
}

// Type de données
func (sireneUL SireneUL) Type() string {
	return "sirene_ul"
}

// Scope de l'objet
func (sireneUL SireneUL) Scope() string {
	return "entreprise"
}

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &sireneULParser{}

type sireneULParser struct {
	file   *os.File
	reader *csv.Reader
	idx    marshal.ColMapping
}

func (parser *sireneULParser) Type() string {
	return "sirene_ul"
}

func (parser *sireneULParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *sireneULParser) Close() error {
	return parser.file.Close()
}

func (parser *sireneULParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ',', true)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, SireneUL{})
	}
	return err
}

func (parser *sireneULParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseSireneUlLine(parser.idx, row, parsedLine)
	})
}

func parseSireneUlLine(idx marshal.ColMapping, row []string, parsedLine *marshal.ParsedLineResult) {
	idxRow := idx.IndexRow(row)
	sireneul := SireneUL{}
	sireneul.Siren = idxRow.GetVal("siren")
	sireneul.RaisonSociale = idxRow.GetVal("denominationUniteLegale")
	sireneul.Prenom1UniteLegale = idxRow.GetVal("prenom1UniteLegale")
	sireneul.Prenom2UniteLegale = idxRow.GetVal("prenom2UniteLegale")
	sireneul.Prenom3UniteLegale = idxRow.GetVal("prenom3UniteLegale")
	sireneul.Prenom4UniteLegale = idxRow.GetVal("prenom4UniteLegale")
	sireneul.NomUniteLegale = idxRow.GetVal("nomUniteLegale")
	sireneul.NomUsageUniteLegale = idxRow.GetVal("nomUsageUniteLegale")
	sireneul.CodeStatutJuridique = idxRow.GetVal("categorieJuridiqueUniteLegale")
	creation, err := time.Parse("2006-01-02", idxRow.GetVal("dateCreationUniteLegale")) // note: cette date n'est pas toujours présente, et on ne souhaite pas être rapporter d'erreur en cas d'absence
	if err == nil {
		sireneul.Creation = &creation
	}
	parsedLine.AddTuple(sireneul)
}
