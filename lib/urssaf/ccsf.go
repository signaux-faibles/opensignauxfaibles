package urssaf

import (
	"encoding/csv"
	"errors"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
)

// CCSF information urssaf ccsf
type CCSF struct {
	key            string    `                                                   csv:"siret"`
	NumeroCompte   string    `input:"Compte"              json:"-"               csv:"numéro_compte"`
	DateTraitement time.Time `input:"Date_de_traitement"  json:"date_traitement" csv:"date_traitement"`
	Stade          string    `input:"Code_externe_stade"  json:"stade"           csv:"stade"`
	Action         string    `input:"Code_externe_action" json:"action"          csv:"action"`
}

// Key _id de l'objet
func (ccsf CCSF) Key() string {
	return ccsf.key
}

// Scope de l'objet
func (ccsf CCSF) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (ccsf CCSF) Type() string {
	return "ccsf"
}

// ParserCCSF fournit une instance utilisable par ParseFilesFromBatch.
var ParserCCSF = &ccsfParser{}

type ccsfParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
	idx     marshal.ColMapping
}

func (parser *ccsfParser) Type() string {
	return "ccsf"
}

func (parser *ccsfParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *ccsfParser) Close() error {
	return parser.file.Close()
}

func (parser *ccsfParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, CCSF{})
	}
	return err
}

func (parser *ccsfParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseCcsfLine(parser.idx, row, &parser.comptes, parsedLine)
	})
}

func parseCcsfLine(idx marshal.ColMapping, row []string, comptes *marshal.Comptes, parsedLine *marshal.ParsedLineResult) {
	var err error
	ccsf := CCSF{}
	if len(row) >= 4 {
		idxRow := idx.IndexRow(row)
		ccsf.Action = idxRow.GetVal("Code_externe_action")
		ccsf.Stade = idxRow.GetVal("Code_externe_stade")
		ccsf.DateTraitement, err = marshal.UrssafToDate(idxRow.GetVal("Date_de_traitement"))
		parsedLine.AddRegularError(err)
		if err != nil {
			return
		}

		ccsf.key, err = marshal.GetSiretFromComptesMapping(idxRow.GetVal("Compte"), &ccsf.DateTraitement, *comptes)
		if err != nil {
			// Compte filtré
			parsedLine.SetFilterError(err)
			return
		}
		ccsf.NumeroCompte = idxRow.GetVal("Compte")
	} else {
		parsedLine.AddRegularError(errors.New("ligne non conforme, moins de 4 champs"))
	}
	if len(parsedLine.Errors) == 0 {
		parsedLine.AddTuple(ccsf)
	}
}
