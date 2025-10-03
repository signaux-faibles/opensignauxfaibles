package apdemande

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/misc"
)

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string    `input:"ID_DA"              json:"id_demande"          sql:"id_demande"           csv:"id_demande"`
	Siret              string    `input:"ETAB_SIRET"         json:"-"                   sql:"siret"                csv:"siret"`
	EffectifEntreprise *int      `input:"EFF_ENT"            json:"effectif_entreprise"                            csv:"effectif_entreprise"`
	Effectif           *int      `input:"EFF_ETAB"           json:"effectif"                                       csv:"effectif"`
	DateStatut         time.Time `input:"DATE_STATUT"        json:"date_statut"         sql:"date_statut"          csv:"date_statut"`
	PeriodeDebut       time.Time `input:"DATE_DEB"           json:"periode_debut"       sql:"periode_debut"        csv:"période_début"`
	PeriodeFin         time.Time `input:"DATE_FIN"           json:"periode_fin"         sql:"periode_fin"          csv:"période_fin"`
	HTA                *float64  `input:"HTA"                json:"hta"                 sql:"heures"               csv:"heures_autorisées"`
	MTA                *float64  `                           json:"mta"                 sql:"montant"              csv:"montants_autorisés"`
	EffectifAutorise   *int      `input:"EFF_AUTO"           json:"effectif_autorise"   sql:"effectif"             csv:"effectif_autorisé"`
	MotifRecoursSE     *int      `input:"MOTIF_RECOURS_SE"   json:"motif_recours_se"    sql:"motif_recours"        csv:"motif_recours_se"`
	HeureConsommee     *float64  `input:"S_HEURE_CONSOM_TOT" json:"heures_consommees"                              csv:"heure_consommee"`
	MontantConsomme    *float64  `                           json:"montant_consomme"                               csv:"montant_consomme"`
	EffectifConsomme   *int      `input:"S_HEURE_CONSOM_TOT" json:"effectif_consomme"                              csv:"effectif_consomme"`
	Perimetre          *int      `input:"PERIMETRE_AP"       json:"perimetre"                                      csv:"perimetre"`
}

// Key id de l'objet
func (apdemande APDemande) Key() string {
	return apdemande.Siret
}

// Type de données
func (apdemande APDemande) Type() base.ParserType {
	return base.Apdemande
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
	idx    engine.ColMapping
}

func (parser *apdemandeParser) Type() base.ParserType {
	return base.Apdemande
}

func (parser *apdemandeParser) Init(cache *engine.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *apdemandeParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ',', true)
	if err == nil {
		parser.idx, err = engine.IndexColumnsFromCsvHeader(parser.reader, APDemande{})
	}
	return err
}

func (parser *apdemandeParser) Close() error {
	return parser.file.Close()
}

func (parser *apdemandeParser) ParseLines(parsedLineChan chan engine.ParsedLineResult) {
	engine.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *engine.ParsedLineResult) {
		parser.parseLine(row, parsedLine)
	})
}

func (parser *apdemandeParser) parseLine(row []string, parsedLine *engine.ParsedLineResult) {
	idxRow := parser.idx.IndexRow(row)
	if idxRow.GetVal("ETAB_SIRET") == "" {
		parsedLine.AddRegularError(errors.New("invalidLine"))
	} else {
		parseApDemandeLine(idxRow, parsedLine)
		if len(parsedLine.Errors) > 0 {
			parsedLine.Tuples = []engine.Tuple{}
		}
	}
}

func parseApDemandeLine(idxRow engine.IndexedRow, parsedLine *engine.ParsedLineResult) {

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
	apdemande.PeriodeDebut, err = time.Parse("2006-01-02", idxRow.GetVal("DATE_DEB"))
	parsedLine.AddRegularError(err)
	apdemande.PeriodeFin, err = time.Parse("2006-01-02", idxRow.GetVal("DATE_FIN"))
	parsedLine.AddRegularError(err)
	apdemande.HTA, err = idxRow.GetFloat64("HTA")
	parsedLine.AddRegularError(err)
	apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("MTA"), ",", "."))
	parsedLine.AddRegularError(err)
	apdemande.EffectifAutorise, err = idxRow.GetIntFromFloat("EFF_AUTO")
	parsedLine.AddRegularError(err)
	motifRecoursSE, err := idxRow.GetInt("MOTIF_RECOURS_SE")
	parsedLine.AddRegularError(err)

	if motifRecoursSE != nil {
		if *motifRecoursSE >= 1 && *motifRecoursSE <= 7 {
			apdemande.MotifRecoursSE = motifRecoursSE
		} else {
			parsedLine.AddRegularError(fmt.Errorf("property \"MOTIF_RECOURS_SE\" a une valeur invalide : %d. Valeur ignorée", *motifRecoursSE))
		}
	}

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
