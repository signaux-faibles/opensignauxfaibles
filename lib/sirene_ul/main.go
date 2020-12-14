package sireneul

import (
	//"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

// SireneUL informations sur les entreprises
type SireneUL struct {
	Siren               string     `json:"siren,omitempty"         bson:"siren,omitempty"`
	Nic                 string     `json:"nic,omitempty"           bson:"nic,omitempty"`
	RaisonSociale       string     `json:"raison_sociale"          bson:"raison_sociale"`
	Prenom1UniteLegale  string     `json:"prenom1_unite_legale,omitempty"      bson:"prenom1_unite_legale,omitempty"`
	Prenom2UniteLegale  string     `json:"prenom2_unite_legale,omitempty"      bson:"prenom2_unite_legale,omitempty"`
	Prenom3UniteLegale  string     `json:"prenom3_unite_legale,omitempty"      bson:"prenom3_unite_legale,omitempty"`
	Prenom4UniteLegale  string     `json:"prenom4_unite_legale,omitempty"      bson:"prenom4_unite_legale,omitempty"`
	NomUniteLegale      string     `json:"nom_unite_legale,omitempty"          bson:"nom_unite_legale,omitempty"`
	NomUsageUniteLegale string     `json:"nom_usage_unite_legale,omitempty"     bson:"nom_usage_unite_legale,omitempty"`
	CodeStatutJuridique string     `json:"statut_juridique"        bson:"statut_juridique"`
	Creation            *time.Time `json:"date_creation,omitempty" bson:"date_creation,omitempty"`
}

var expectedHeader = "siren,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale,sigleUniteLegale,sexeUniteLegale,prenom1UniteLegale,prenom2UniteLegale,prenom3UniteLegale,prenom4UniteLegale,prenomUsuelUniteLegale,pseudonymeUniteLegale,identifiantAssociationUniteLegale,trancheEffectifsUniteLegale,anneeEffectifsUniteLegale,dateDernierTraitementUniteLegale,nombrePeriodesUniteLegale,categorieEntreprise,anneeCategorieEntreprise,dateDebut,etatAdministratifUniteLegale,nomUniteLegale,nomUsageUniteLegale,denominationUniteLegale,denominationUsuelle1UniteLegale,denominationUsuelle2UniteLegale,denominationUsuelle3UniteLegale,categorieJuridiqueUniteLegale,activitePrincipaleUniteLegale,nomenclatureActivitePrincipaleUniteLegale,nicSiegeUniteLegale,economieSocialeSolidaireUniteLegale,caractereEmployeurUniteLegale"

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
	parser.file, err = os.Open(filePath)
	if err != nil {
		return err
	}
	parser.reader = csv.NewReader(parser.file)
	parser.reader.Comma = ','
	parser.reader.LazyQuotes = true

	// parse header
	row, err := parser.reader.Read()
	if err != nil {
		return err // may be io.EOF
	} else if strings.Join(row, ",") != expectedHeader {
		return errors.New("sirene_ul header does not match the parser's expectations")
	}

	return err
}

func (parser *sireneUlParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parseSireneUlLine(row, &parsedLine)
		}
		parsedLineChan <- parsedLine
	}
}

func parseSireneUlLine(row []string, parsedLine *marshal.ParsedLineResult) {
	sireneul := SireneUL{}
	sireneul.Siren = row[0]
	sireneul.RaisonSociale = row[23]
	sireneul.Prenom1UniteLegale = row[6]
	sireneul.Prenom2UniteLegale = row[7]
	sireneul.Prenom3UniteLegale = row[8]
	sireneul.Prenom4UniteLegale = row[9]
	sireneul.NomUniteLegale = row[21]
	sireneul.NomUsageUniteLegale = row[22]
	sireneul.CodeStatutJuridique = row[27]
	creation, err := time.Parse("2006-01-02", row[3]) // note: cette date n'est pas toujours présente, et on ne souhaite pas être rapporter d'erreur en cas d'absence
	if err == nil {
		sireneul.Creation = &creation
	}
	parsedLine.AddTuple(sireneul)
}
