package urssaf

import (
	"encoding/csv"
	"errors"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
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
func (ccsf CCSF) Type() base.ParserType {
	return base.Ccsf
}

// ParserCCSF fournit une instance utilisable par ParseFilesFromBatch.
var ParserCCSF = &ccsfParser{}

type ccsfParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes engine.Comptes
	idx     engine.ColMapping
}

func (parser *ccsfParser) Type() base.ParserType {
	return base.Ccsf
}

func (parser *ccsfParser) Init(cache *engine.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = engine.GetCompteSiretMapping(*cache, batch, engine.OpenAndReadSiretMapping)
	return err
}

func (parser *ccsfParser) Close() error {
	return parser.file.Close()
}

func (parser *ccsfParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, err = engine.IndexColumnsFromCsvHeader(parser.reader, CCSF{})
	}
	return err
}

func (parser *ccsfParser) ParseLines(parsedLineChan chan engine.ParsedLineResult) {
	engine.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *engine.ParsedLineResult) {
		parseCcsfLine(parser.idx, row, &parser.comptes, parsedLine)
	})
}

func parseCcsfLine(idx engine.ColMapping, row []string, comptes *engine.Comptes, parsedLine *engine.ParsedLineResult) {
	var err error
	ccsf := CCSF{}
	if len(row) >= 4 {
		idxRow := idx.IndexRow(row)
		ccsf.Action = idxRow.GetVal("Code_externe_action")
		ccsf.Stade = idxRow.GetVal("Code_externe_stade")
		ccsf.DateTraitement, err = engine.UrssafToDate(idxRow.GetVal("Date_de_traitement"))
		parsedLine.AddRegularError(err)
		if err != nil {
			return
		}

		ccsf.key, err = engine.GetSiretFromComptesMapping(idxRow.GetVal("Compte"), &ccsf.DateTraitement, *comptes)
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
