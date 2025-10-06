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

type apdemandeParser struct {
	file   *os.File
	reader *csv.Reader
	idx    engine.ColMapping
}

func NewParserApdemande() *apdemandeParser {
	return &apdemandeParser{}
}

func (parser *apdemandeParser) Type() base.ParserType {
	return base.Apdemande
}

func (parser *apdemandeParser) Init(_ *engine.Cache, _ engine.SirenFilter, _ *base.AdminBatch) error {
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

func (parser *apdemandeParser) ReadNext(res *engine.ParsedLineResult) error {
	row, err := parser.reader.Read()
	if err != nil {
		return err
	}

	idxRow := parser.idx.IndexRow(row)

	if idxRow.GetVal("ETAB_SIRET") == "" {
		res.AddRegularError(errors.New("invalidLine: siret unspecified"))
		return nil
	}

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
	apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("MTA"), ",", "."))
	res.AddRegularError(err)
	apdemande.EffectifAutorise, err = idxRow.GetIntFromFloat("EFF_AUTO")
	res.AddRegularError(err)
	motifRecoursSE, err := idxRow.GetInt("MOTIF_RECOURS_SE")
	res.AddRegularError(err)

	if motifRecoursSE != nil {
		if *motifRecoursSE >= 1 && *motifRecoursSE <= 7 {
			apdemande.MotifRecoursSE = motifRecoursSE
		} else {
			res.AddRegularError(fmt.Errorf("property \"MOTIF_RECOURS_SE\" a une valeur invalide : %d. Valeur ignorée", *motifRecoursSE))
		}
	}

	apdemande.HeureConsommee, err = idxRow.GetFloat64("S_HEURE_CONSOM_TOT")
	res.AddRegularError(err)
	apdemande.EffectifConsomme, err = idxRow.GetIntFromFloat("S_EFF_CONSOM_TOT")
	res.AddRegularError(err)
	apdemande.MontantConsomme, err = misc.ParsePFloat(strings.ReplaceAll(idxRow.GetVal("S_MONTANT_CONSOM_TOT"), ",", "."))
	res.AddRegularError(err)
	apdemande.Perimetre, err = idxRow.GetInt("PERIMETRE_AP")
	res.AddRegularError(err)

	if len(res.Errors) == 0 {
		res.AddTuple(apdemande)
	}
	return nil
}
