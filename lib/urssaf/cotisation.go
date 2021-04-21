package urssaf

import (
	"encoding/csv"
	"os"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
)

// Cotisation Objet cotisation
type Cotisation struct {
	key          string       `                 hash:"-"`
	NumeroCompte string       `col:"Compte"     json:"numero_compte" bson:"numero_compte"`
	Periode      misc.Periode `col:"periode"    json:"periode"       bson:"periode"`
	Encaisse     *float64     `col:"enc_direct" json:"encaisse"      bson:"encaisse"`
	Du           *float64     `col:"cotis_due"  json:"du"            bson:"du"`
}

// Key _id de l'objet
func (cotisation Cotisation) Key() string {
	return cotisation.key
}

// Scope de l'objet
func (cotisation Cotisation) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (cotisation Cotisation) Type() string {
	return "cotisation"
}

// ParserCotisation fournit une instance utilisable par ParseFilesFromBatch.
var ParserCotisation = &cotisationParser{}

type cotisationParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
	idx     marshal.ColMapping
}

func (parser *cotisationParser) GetFileType() string {
	return "cotisation"
}

func (parser *cotisationParser) Close() error {
	return parser.file.Close()
}

func (parser *cotisationParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *cotisationParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ';', true)
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

	periode, err := marshal.UrssafToPeriod(idxRow.GetVal("periode"))
	date := periode.Start
	parsedLine.AddRegularError(err)

	siret, err := marshal.GetSiretFromComptesMapping(idxRow.GetVal("Compte"), &date, *comptes)
	if err != nil {
		parsedLine.SetFilterError(err)
	} else {
		cotisation.key = siret
		cotisation.NumeroCompte = idxRow.GetVal("Compte")
		cotisation.Periode, err = marshal.UrssafToPeriod(idxRow.GetVal("periode"))
		parsedLine.AddRegularError(err)
		cotisation.Encaisse, err = idxRow.GetCommaFloat64("enc_direct")
		parsedLine.AddRegularError(err)
		cotisation.Du, err = idxRow.GetCommaFloat64("cotis_due")
		parsedLine.AddRegularError(err)
	}
	if len(parsedLine.Errors) == 0 {
		parsedLine.AddTuple(cotisation)
	}
}
