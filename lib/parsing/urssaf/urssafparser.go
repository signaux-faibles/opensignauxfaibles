package urssaf

import (
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

// UrssafRowParser est la base commune pour les parsers de l'URSSAF
type UrssafRowParser struct {
	comptes *engine.Comptes
}

func (p *UrssafRowParser) GetComptes() *engine.Comptes {
	return p.comptes
}

func (p *UrssafRowParser) setComptes(comptes engine.Comptes) {
	p.comptes = &comptes
}

// UrssaParserInst factorise l'initialisation des comptes pour tous les
// parsers URSSAF
type UrssafParserInst struct {
	parsing.CsvParserInst
}

func (p *UrssafParserInst) Init(cache *engine.Cache, filter engine.SirenFilter, batch *base.AdminBatch) error {
	err := p.CsvParserInst.Init(cache, filter, batch)
	if err != nil {
		return err
	}

	if urssafRowParser, ok := p.RowParser.(interface{ setComptes(engine.Comptes) }); ok {
		comptes, err := engine.GetCompteSiretMapping(*cache, batch, filter, engine.OpenAndReadSiretMapping)
		if err != nil {
			return err
		}

		urssafRowParser.setComptes(comptes)
	}

	return err
}
