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
	Siren               string     `col:"siren"                         json:"siren,omitempty"                  bson:"siren,omitempty"`
	Nic                 string     `                                    json:"nic,omitempty"                    bson:"nic,omitempty"`
	RaisonSociale       string     `col:"denominationUniteLegale"       json:"raison_sociale"                   bson:"raison_sociale"`
	Prenom1UniteLegale  string     `col:"prenom1UniteLegale"            json:"prenom1_unite_legale,omitempty"   bson:"prenom1_unite_legale,omitempty"`
	Prenom2UniteLegale  string     `col:"prenom2UniteLegale"            json:"prenom2_unite_legale,omitempty"   bson:"prenom2_unite_legale,omitempty"`
	Prenom3UniteLegale  string     `col:"prenom3UniteLegale"            json:"prenom3_unite_legale,omitempty"   bson:"prenom3_unite_legale,omitempty"`
	Prenom4UniteLegale  string     `col:"prenom4UniteLegale"            json:"prenom4_unite_legale,omitempty"   bson:"prenom4_unite_legale,omitempty"`
	NomUniteLegale      string     `col:"nomUniteLegale"                json:"nom_unite_legale,omitempty"       bson:"nom_unite_legale,omitempty"`
	NomUsageUniteLegale string     `col:"nomUsageUniteLegale"           json:"nom_usage_unite_legale,omitempty" bson:"nom_usage_unite_legale,omitempty"`
	CodeStatutJuridique string     `col:"categorieJuridiqueUniteLegale" json:"statut_juridique"                 bson:"statut_juridique"`
	Creation            *time.Time `col:"dateCreationUniteLegale"       json:"date_creation,omitempty"          bson:"date_creation,omitempty"`
}

func (sirene_ul SireneUL) Headers() []string {
	return marshal.ExtractColTags(sirene_ul)
}

func (sirene_ul SireneUL) Values() []string {
	creation := ""
	if sirene_ul.Creation != nil {
		creation = sirene_ul.Creation.Format(time.DateOnly)
	}
	return []string{
		sirene_ul.Siren,
		sirene_ul.RaisonSociale,
		sirene_ul.Prenom1UniteLegale,
		sirene_ul.Prenom2UniteLegale,
		sirene_ul.Prenom3UniteLegale,
		sirene_ul.Prenom4UniteLegale,
		sirene_ul.NomUniteLegale,
		sirene_ul.NomUsageUniteLegale,
		sirene_ul.CodeStatutJuridique,
		creation,
	}
}

// Key id de l'objet
func (sirene_ul SireneUL) Key() string {
	return sirene_ul.Siren
}

// Type de données
func (sirene_ul SireneUL) Type() string {
	return "sirene_ul"
}

// Scope de l'objet
func (sirene_ul SireneUL) Scope() string {
	return "entreprise"
}

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &sireneUlParser{}

type sireneUlParser struct {
	file   *os.File
	reader *csv.Reader
	idx    marshal.ColMapping
}

func (parser *sireneUlParser) GetFileType() string {
	return "sirene_ul"
}

func (parser *sireneUlParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *sireneUlParser) Close() error {
	return parser.file.Close()
}

func (parser *sireneUlParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ',', true)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, SireneUL{})
	}
	return err
}

func (parser *sireneUlParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
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
