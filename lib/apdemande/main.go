package apdemande

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
)

// Periode Période de temps avec un début et une fin

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string       `json:"id_demande" bson:"id_demande"`
	Siret              string       `json:"-" bson:"-"`
	EffectifEntreprise *int         `json:"effectif_entreprise" bson:"effectif_entreprise"`
	Effectif           *int         `json:"effectif" bson:"effectif"`
	DateStatut         time.Time    `json:"date_statut" bson:"date_statut"`
	Periode            misc.Periode `json:"periode" bson:"periode"`
	HTA                *float64     `json:"hta" bson:"hta"`
	MTA                *float64     `json:"mta" bson:"mta"`
	EffectifAutorise   *int         `json:"effectif_autorise" bson:"effectif_autorise"`
	MotifRecoursSE     *int         `json:"motif_recours_se" bson:"motif_recours_se"`
	HeureConsommee     *float64     `json:"heure_consommee" bson:"heure_consommee"`
	MontantConsomme    *float64     `json:"montant_consommee" bson:"montant_consommee"`
	EffectifConsomme   *int         `json:"effectif_consomme" bson:"effectif_consomme"`
}

// Key id de l'objet
func (apdemande APDemande) Key() string {
	return apdemande.Siret
}

// Type de données
func (apdemande APDemande) Type() string {
	return "apdemande"
}

// Scope de l'objet
func (apdemande APDemande) Scope() string {
	return "etablissement"
}

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &apdemandeParser{}

type apdemandeParser struct {
	file   *os.File
	reader *csv.Reader
	idx    marshal.ColMapping
}

func (parser *apdemandeParser) GetFileType() string {
	return "apdemande"
}

func (parser *apdemandeParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *apdemandeParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = openFile(filePath)
	if err == nil {
		parser.idx, err = parseColMapping(parser.reader)
	}
	return err
}

func (parser *apdemandeParser) Close() error {
	return parser.file.Close()
}

func openFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	return file, reader, nil
}

func parseColMapping(reader *csv.Reader) (marshal.ColMapping, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	var idx = marshal.GetFieldBindings(header)
	requiredFields := []string{
		"ID_DA",
		"ETAB_SIRET",
		"EFF_ENT",
		"EFF_ETAB",
		"DATE_STATUT",
		"DATE_DEB",
		"DATE_FIN",
		"HTA",
		"EFF_AUTO",
		"MOTIF_RECOURS_SE",
		"S_HEURE_CONSOM_TOT",
		"S_EFF_CONSOM_TOT",
	}
	if _, err := idx.HasFields(requiredFields); err != nil {
		return nil, err
	}
	return idx, nil
}

func (parser *apdemandeParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else if row[parser.idx["ETAB_SIRET"]] == "" {
			parsedLine.AddRegularError(errors.New("invalidLine"))
		} else {
			parseApDemandeLine(row, parser.idx, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []marshal.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseApDemandeLine(row []string, idx marshal.ColMapping, parsedLine *marshal.ParsedLineResult) {
	apdemande := APDemande{}
	apdemande.ID = row[idx["ID_DA"]]
	apdemande.Siret = row[idx["ETAB_SIRET"]]
	var err error
	apdemande.EffectifEntreprise, err = misc.ParsePInt(row[idx["EFF_ENT"]])
	parsedLine.AddRegularError(err)
	apdemande.Effectif, err = misc.ParsePInt(row[idx["EFF_ETAB"]])
	parsedLine.AddRegularError(err)
	apdemande.DateStatut, err = time.Parse("02/01/2006", row[idx["DATE_STATUT"]])
	parsedLine.AddRegularError(err)
	apdemande.Periode = misc.Periode{}
	apdemande.Periode.Start, err = time.Parse("02/01/2006", row[idx["DATE_DEB"]])
	parsedLine.AddRegularError(err)
	apdemande.Periode.End, err = time.Parse("02/01/2006", row[idx["DATE_FIN"]])
	parsedLine.AddRegularError(err)
	apdemande.HTA, err = misc.ParsePFloat(row[idx["HTA"]])
	parsedLine.AddRegularError(err)
	apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(row[idx["MTA"]], ",", "."))
	parsedLine.AddRegularError(err)
	apdemande.EffectifAutorise, err = misc.ParsePInt(row[idx["EFF_AUTO"]])
	parsedLine.AddRegularError(err)
	apdemande.MotifRecoursSE, err = misc.ParsePInt(row[idx["MOTIF_RECOURS_SE"]])
	parsedLine.AddRegularError(err)
	apdemande.HeureConsommee, err = misc.ParsePFloat(row[idx["S_HEURE_CONSOM_TOT"]])
	parsedLine.AddRegularError(err)
	apdemande.EffectifConsomme, err = misc.ParsePInt(row[idx["S_EFF_CONSOM_TOT"]])
	parsedLine.AddRegularError(err)
	apdemande.MontantConsomme, err = misc.ParsePFloat(strings.ReplaceAll(row[idx["S_MONTANT_CONSOM_TOT"]], ",", "."))
	parsedLine.AddRegularError(err)
	parsedLine.AddTuple(apdemande)
}
