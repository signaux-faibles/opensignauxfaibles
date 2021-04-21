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
	ID                 string       `col:"ID_DA" json:"id_demande" bson:"id_demande"`
	Siret              string       `col:"ETAB_SIRET" json:"-" bson:"-"`
	EffectifEntreprise *int         `col:"EFF_ENT" json:"effectif_entreprise" bson:"effectif_entreprise"`
	Effectif           *int         `col:"EFF_ETAB" json:"effectif" bson:"effectif"`
	DateStatut         time.Time    `col:"DATE_STATUT" json:"date_statut" bson:"date_statut"`
	Periode            misc.Periode `cols:"DATE_DEB,DATE_FIN" json:"periode" bson:"periode"`
	HTA                *float64     `col:"HTA" json:"hta" bson:"hta"`
	MTA                *float64     `json:"mta" bson:"mta"`
	EffectifAutorise   *int         `col:"EFF_AUTO" json:"effectif_autorise" bson:"effectif_autorise"`
	MotifRecoursSE     *int         `col:"MOTIF_RECOURS_SE" json:"motif_recours_se" bson:"motif_recours_se"`
	HeureConsommee     *float64     `col:"S_HEURE_CONSOM_TOT" json:"heure_consommee" bson:"heure_consommee"`
	MontantConsomme    *float64     `json:"montant_consommee" bson:"montant_consommee"`
	EffectifConsomme   *int         `col:"S_HEURE_CONSOM_TOT" json:"effectif_consomme" bson:"effectif_consomme"`
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
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ',', true)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, APDemande{})
	}
	return err
}

func (parser *apdemandeParser) Close() error {
	return parser.file.Close()
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
		} else {
			idxRow := parser.idx.IndexRow(row)
			if idxRow.GetVal("ETAB_SIRET") == "" {
				parsedLine.AddRegularError(errors.New("invalidLine"))
			} else {
				parseApDemandeLine(idxRow, &parsedLine)
				if len(parsedLine.Errors) > 0 {
					parsedLine.Tuples = []marshal.Tuple{}
				}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseApDemandeLine(idxRow marshal.IndexedRow, parsedLine *marshal.ParsedLineResult) {

	apdemande := APDemande{}
	apdemande.ID = idxRow.GetVal("ID_DA")
	apdemande.Siret = idxRow.GetVal("ETAB_SIRET")
	var err error
	apdemande.EffectifEntreprise, err = idxRow.GetInt("EFF_ENT")
	parsedLine.AddRegularError(err)
	apdemande.Effectif, err = idxRow.GetInt("EFF_ETAB")
	parsedLine.AddRegularError(err)
	apdemande.DateStatut, err = time.Parse("02/01/2006", idxRow.GetVal("DATE_STATUT"))
	parsedLine.AddRegularError(err)
	apdemande.Periode = misc.Periode{}
	apdemande.Periode.Start, err = time.Parse("02/01/2006", idxRow.GetVal("DATE_DEB"))
	parsedLine.AddRegularError(err)
	apdemande.Periode.End, err = time.Parse("02/01/2006", idxRow.GetVal("DATE_FIN"))
	parsedLine.AddRegularError(err)
	apdemande.HTA, err = idxRow.GetFloat64("HTA")
	parsedLine.AddRegularError(err)
	apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("MTA"), ",", "."))
	parsedLine.AddRegularError(err)
	apdemande.EffectifAutorise, err = idxRow.GetInt("EFF_AUTO")
	parsedLine.AddRegularError(err)
	apdemande.MotifRecoursSE, err = idxRow.GetInt("MOTIF_RECOURS_SE")
	parsedLine.AddRegularError(err)
	apdemande.HeureConsommee, err = idxRow.GetFloat64("S_HEURE_CONSOM_TOT")
	parsedLine.AddRegularError(err)
	apdemande.EffectifConsomme, err = idxRow.GetInt("S_EFF_CONSOM_TOT")
	parsedLine.AddRegularError(err)
	apdemande.MontantConsomme, err = misc.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("S_MONTANT_CONSOM_TOT"), ",", "."))
	parsedLine.AddRegularError(err)
	parsedLine.AddTuple(apdemande)
}
