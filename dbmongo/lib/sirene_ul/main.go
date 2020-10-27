package sireneul

import (
	//"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

	"github.com/signaux-faibles/gournal"
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

// Parser expose le parseur et le type de fichier qu'il supporte.
var Parser = marshal.Parser{FileType: "sirene_ul", FileParser: ParseFile}

// ParseFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	filter := marshal.GetSirenFilterFromCache(*cache)
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	parseSireneULFile(reader, filter, tracker, outputChannel)
}

func parseSireneULFile(reader *csv.Reader, filter marshal.SirenFilter, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	_, err := reader.Read()
	if err == io.EOF {
		return
	} else if err != nil {
		tracker.Add(err)
		return
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
		} else if !sfregexp.ValidSiren(row[0]) {
			tracker.Add(errors.New("siren invalide : " + row[0]))
		} else if !filter.Skips(row[0]) {
			outputChannel <- parseSireneUlLine(row, tracker)
		}
		tracker.Next()
	}
}

func parseSireneUlLine(row []string, tracker *gournal.Tracker) SireneUL {
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
	creation, err := time.Parse("2006-01-02", row[3])
	if err == nil {
		sireneul.Creation = &creation
	}
	tracker.Add(err)
	return sireneul
}
