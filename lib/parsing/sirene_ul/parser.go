package sireneul

import (
	"encoding/csv"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
)

type sireneULParser struct {
	file   *os.File
	reader *csv.Reader
	idx    engine.ColMapping
}

func NewParserSireneUL() *sireneULParser {
	return &sireneULParser{}
}

func (parser *sireneULParser) Type() base.ParserType {
	return base.SireneUl
}

func (parser *sireneULParser) Init(_ *engine.Cache, _ engine.SirenFilter, _ *base.AdminBatch) error {
	return nil
}

func (parser *sireneULParser) Close() error {
	return parser.file.Close()
}

func (parser *sireneULParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ',', true)
	if err == nil {
		parser.idx, err = engine.IndexColumnsFromCsvHeader(parser.reader, SireneUL{})
	}
	return err
}

func (parser *sireneULParser) ReadNext(res *engine.ParsedLineResult) error {
	row, err := parser.reader.Read()
	if err != nil {
		return err
	}

	idxRow := parser.idx.IndexRow(row)
	sireneul := SireneUL{}
	sireneul.Siren = idxRow.GetVal("siren")
	sireneul.RaisonSociale = idxRow.GetVal("denominationUniteLegale")
	sireneul.Prenom1UniteLegale = idxRow.GetVal("prenom1UniteLegale")
	sireneul.Prenom2UniteLegale = idxRow.GetVal("prenom2UniteLegale")
	sireneul.Prenom3UniteLegale = idxRow.GetVal("prenom3UniteLegale")
	sireneul.Prenom4UniteLegale = idxRow.GetVal("prenom4UniteLegale")
	sireneul.NomUniteLegale = idxRow.GetVal("nomUniteLegale")
	sireneul.NomUsageUniteLegale = idxRow.GetVal("nomUsageUniteLegale")
	sireneul.CodeStatutJuridique = idxRow.GetVal("categorieJuridiqueUniteLegale")

	creation, err := time.Parse("2006-01-02", idxRow.GetVal("dateCreationUniteLegale")) // note: cette date n'est pas toujours présente, et on ne souhaite pas être rapporter d'erreur en cas d'absence
	if err == nil {
		sireneul.Creation = &creation
	}

	res.AddTuple(sireneul)
	return nil
}
