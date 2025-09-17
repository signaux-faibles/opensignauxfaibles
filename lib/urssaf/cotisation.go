package urssaf

import (
	"encoding/csv"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
)

// Cotisation Objet cotisation
type Cotisation struct {
	Siret        string    `                   json:"-"             sql:"siret"            csv:"siret"`
	NumeroCompte string    `input:"Compte"     json:"numero_compte"                      csv:"numéro_compte"`
	PeriodeDebut time.Time `input:"periode"    json:"periode_debut" sql:"periode_debut"  csv:"période_début"`
	PeriodeFin   time.Time `input:"periode"    json:"periode_fin"   sql:"periode_fin"    csv:"période_fin"`
	Encaisse     *float64  `input:"enc_direct" json:"encaisse"                           csv:"encaissé"`
	Du           *float64  `input:"cotis_due"  json:"du"            sql:"du"             csv:"du"`
}

// Key _id de l'objet
func (cotisation Cotisation) Key() string {
	return cotisation.Siret
}

// Scope de l'objet
func (cotisation Cotisation) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (cotisation Cotisation) Type() base.ParserType {
	return base.Cotisation
}

// ParserCotisation fournit une instance utilisable par ParseFilesFromBatch.
var ParserCotisation = &cotisationParser{}

type cotisationParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
	idx     marshal.ColMapping
}

func (parser *cotisationParser) Type() base.ParserType {
	return base.Cotisation
}

func (parser *cotisationParser) Close() error {
	return parser.file.Close()
}

func (parser *cotisationParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *cotisationParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ';', true)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, Cotisation{})
	}
	return err
}

func (parser *cotisationParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseCotisationLine(parser.idx, row, &parser.comptes, parsedLine)
	})
}

func parseCotisationLine(idx marshal.ColMapping, row []string, comptes *marshal.Comptes, parsedLine *marshal.ParsedLineResult) {
	idxRow := idx.IndexRow(row)
	cotisation := Cotisation{}

	periodeDebut, periodeFin, err := marshal.UrssafToPeriod(idxRow.GetVal("periode"))
	parsedLine.AddRegularError(err)

	siret, err := marshal.GetSiretFromComptesMapping(idxRow.GetVal("Compte"), &periodeDebut, *comptes)
	if err != nil {
		parsedLine.SetFilterError(err)
	} else {
		cotisation.Siret = siret
		cotisation.NumeroCompte = idxRow.GetVal("Compte")
		cotisation.PeriodeDebut = periodeDebut
		cotisation.PeriodeFin = periodeFin
		cotisation.Encaisse, err = idxRow.GetCommaFloat64("enc_direct")
		parsedLine.AddRegularError(err)
		cotisation.Du, err = idxRow.GetCommaFloat64("cotis_due")
		parsedLine.AddRegularError(err)
	}
	if len(parsedLine.Errors) == 0 {
		parsedLine.AddTuple(cotisation)
	}
}
