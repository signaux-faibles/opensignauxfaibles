package apdemande

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/misc"
)

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string       `input:"ID_DA"              json:"id_demande"           csv:"ID_DA"`
	Siret              string       `input:"ETAB_SIRET"         json:"-"                    csv:"ETAB_SIRET"`
	EffectifEntreprise *int         `input:"EFF_ENT"            json:"effectif_entreprise"  csv:"EFF_ENT"`
	Effectif           *int         `input:"EFF_ETAB"           json:"effectif"             csv:"EFF_ETAB"`
	DateStatut         time.Time    `input:"DATE_STATUT"        json:"date_statut"          csv:"DATE_STATUT"`
	Periode            misc.Periode `                           json:"periode"`
	HTA                *float64     `input:"HTA"                json:"hta"                  csv:"HTA"`
	MTA                *float64     `                           json:"mta"                  csv:"MTA"`
	EffectifAutorise   *int         `input:"EFF_AUTO"           json:"effectif_autorise"    csv:"EFF_AUTO"`
	MotifRecoursSE     *int         `input:"MOTIF_RECOURS_SE"   json:"motif_recours_se"     csv:"MOTIF_RECOURS_SE"`
	HeureConsommee     *float64     `input:"S_HEURE_CONSOM_TOT" json:"heure_consommee"      csv:"S_HEURE_CONSOM_TOT"`
	MontantConsomme    *float64     `                           json:"montant_consommee"    csv:"MONTANT_CONSOMME"`
	EffectifConsomme   *int         `input:"S_HEURE_CONSOM_TOT" json:"effectif_consomme"    csv:"S_HEURE_CONSOM_TOT"`
	Perimetre          *int         `input:"PERIMETRE_AP"       json:"perimetre"            csv:"PERIMETRE_AP"`
}

func (apdemande APDemande) Values() []string {
	return []string{
		apdemande.ID,
		apdemande.Siret,
		marshal.IntToCSV(apdemande.EffectifEntreprise),
		marshal.IntToCSV(apdemande.Effectif),
		marshal.TimeToCSV(&apdemande.DateStatut),
		marshal.TimeToCSV(&apdemande.Periode.Start),
		marshal.TimeToCSV(&apdemande.Periode.End),
		marshal.FloatToCSV(apdemande.HTA),
		marshal.FloatToCSV(apdemande.MTA),
		marshal.IntToCSV(apdemande.EffectifAutorise),
		marshal.IntToCSV(apdemande.MotifRecoursSE),
		marshal.FloatToCSV(apdemande.HeureConsommee),
		marshal.FloatToCSV(apdemande.MontantConsomme),
		marshal.IntToCSV(apdemande.EffectifConsomme),
		marshal.IntToCSV(apdemande.Perimetre),
	}
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

func (parser *apdemandeParser) Type() string {
	return "apdemande"
}

func (parser *apdemandeParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *apdemandeParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ',', true)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, APDemande{})
	}
	return err
}

func (parser *apdemandeParser) Close() error {
	return parser.file.Close()
}

func (parser *apdemandeParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parser.parseLine(row, parsedLine)
	})
}

func (parser *apdemandeParser) parseLine(row []string, parsedLine *marshal.ParsedLineResult) {
	idxRow := parser.idx.IndexRow(row)
	if idxRow.GetVal("ETAB_SIRET") == "" {
		parsedLine.AddRegularError(errors.New("invalidLine"))
	} else {
		parseApDemandeLine(idxRow, parsedLine)
		if len(parsedLine.Errors) > 0 {
			parsedLine.Tuples = []marshal.Tuple{}
		}
	}
}

func parseApDemandeLine(idxRow marshal.IndexedRow, parsedLine *marshal.ParsedLineResult) {

	apdemande := APDemande{}
	apdemande.ID = idxRow.GetVal("ID_DA")
	apdemande.Siret = idxRow.GetVal("ETAB_SIRET")
	var err error
	apdemande.EffectifEntreprise, err = idxRow.GetIntFromFloat("EFF_ENT")
	parsedLine.AddRegularError(err)
	apdemande.Effectif, err = idxRow.GetIntFromFloat("EFF_ETAB")
	parsedLine.AddRegularError(err)
	apdemande.DateStatut, err = time.Parse("2006-01-02", idxRow.GetVal("DATE_STATUT"))
	parsedLine.AddRegularError(err)
	apdemande.Periode = misc.Periode{}
	apdemande.Periode.Start, err = time.Parse("2006-01-02", idxRow.GetVal("DATE_DEB"))
	parsedLine.AddRegularError(err)
	apdemande.Periode.End, err = time.Parse("2006-01-02 00:00:00", idxRow.GetVal("DATE_FIN"))
	parsedLine.AddRegularError(err)
	apdemande.HTA, err = idxRow.GetFloat64("HTA")
	parsedLine.AddRegularError(err)
	apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("MTA"), ",", "."))
	parsedLine.AddRegularError(err)
	apdemande.EffectifAutorise, err = idxRow.GetIntFromFloat("EFF_AUTO")
	parsedLine.AddRegularError(err)
	apdemande.MotifRecoursSE, err = idxRow.GetInt("MOTIF_RECOURS_SE")
	parsedLine.AddRegularError(err)
	apdemande.HeureConsommee, err = idxRow.GetFloat64("S_HEURE_CONSOM_TOT")
	parsedLine.AddRegularError(err)
	apdemande.EffectifConsomme, err = idxRow.GetIntFromFloat("S_EFF_CONSOM_TOT")
	parsedLine.AddRegularError(err)
	apdemande.MontantConsomme, err = misc.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("S_MONTANT_CONSOM_TOT"), ",", "."))
	parsedLine.AddRegularError(err)
	apdemande.Perimetre, err = idxRow.GetInt("PERIMETRE_AP")
	parsedLine.AddRegularError(err)
	parsedLine.AddTuple(apdemande)
}
