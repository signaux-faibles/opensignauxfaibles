package apdemande

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type ApdemandeParser struct{}

func NewApdemandeParser() engine.Parser {
	return &ApdemandeParser{}
}

func (p *ApdemandeParser) Type() engine.ParserType { return engine.Apdemande }
func (p *ApdemandeParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:        r,
		RowParser:     &apdemandeRowParser{},
		Comma:         ',',
		LazyQuotes:    true,
		CaseSensitive: false,
		DestTuple:     APDemande{},
	}
}

type apdemandeRowParser struct{}

func (parser *apdemandeRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) error {
	idxRow := idx.IndexRow(row)

	if idxRow.GetVal("ETAB_SIRET") == "" {
		res.AddRegularError(errors.New("invalidLine: siret unspecified"))
		return nil
	}

	var err error

	apdemande := APDemande{}
	apdemande.ID = idxRow.GetVal("ID_DA")
	apdemande.Siret = idxRow.GetVal("ETAB_SIRET")
	apdemande.EffectifEntreprise, err = idxRow.GetIntFromFloat("EFF_ENT")
	res.AddRegularError(err)
	apdemande.Effectif, err = idxRow.GetIntFromFloat("EFF_ETAB")
	res.AddRegularError(err)
	apdemande.DateStatut, err = time.Parse("2006-01-02", idxRow.GetVal("DATE_STATUT"))
	res.AddRegularError(err)
	apdemande.PeriodeDebut, err = time.Parse("2006-01-02", idxRow.GetVal("DATE_DEB"))
	res.AddRegularError(err)
	apdemande.PeriodeFin, err = time.Parse("2006-01-02", idxRow.GetVal("DATE_FIN"))
	res.AddRegularError(err)
	apdemande.HTA, err = idxRow.GetFloat64("HTA")
	res.AddRegularError(err)
	apdemande.MTA, err = parsing.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("MTA"), ",", "."))
	res.AddRegularError(err)
	apdemande.EffectifAutorise, err = idxRow.GetIntFromFloat("EFF_AUTO")
	res.AddRegularError(err)
	motifRecoursSE, err := idxRow.GetInt("MOTIF_RECOURS_SE")
	res.AddRegularError(err)

	if motifRecoursSE != nil {
		if *motifRecoursSE >= 1 && *motifRecoursSE <= 7 {
			apdemande.MotifRecoursSE = motifRecoursSE
		} else {
			res.AddRegularError(fmt.Errorf("property \"MOTIF_RECOURS_SE\" has invalid value : %d. Value ignored", *motifRecoursSE))
		}
	}

	apdemande.HeureConsommee, err = idxRow.GetFloat64("S_HEURE_CONSOM_TOT")
	res.AddRegularError(err)
	apdemande.EffectifConsomme, err = idxRow.GetIntFromFloat("S_EFF_CONSOM_TOT")
	res.AddRegularError(err)
	apdemande.MontantConsomme, err = parsing.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("S_MONTANT_CONSOM_TOT"), ",", "."))
	res.AddRegularError(err)
	apdemande.Perimetre, err = idxRow.GetInt("PERIMETRE_AP")
	res.AddRegularError(err)

	if len(res.Errors) == 0 {
		res.AddTuple(apdemande)
	}
	return nil
}
