// Package urssaf expose tous les parsers pour extraire les données des
// URSSAF
package urssaf

import (
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

// UrssafRowParser est la base commune pour les parsers de l'URSSAF
type UrssafRowParser struct {
	comptes *Comptes
}

func (p *UrssafRowParser) GetComptes() *Comptes {
	return p.comptes
}

func (p *UrssafRowParser) setComptes(comptes Comptes) {
	p.comptes = &comptes
}

// UrssafParserInst factorise l'initialisation des comptes pour tous les
// parsers URSSAF
type UrssafParserInst struct {
	parsing.CsvParserInst
}

func (p *UrssafParserInst) Init(cache *engine.Cache, filter engine.SirenFilter, batch *engine.AdminBatch) error {
	err := p.CsvParserInst.Init(cache, filter, batch)
	if err != nil {
		return err
	}

	if urssafRowParser, ok := p.RowParser.(interface{ setComptes(Comptes) }); ok {
		comptes, err := GetCompteSiretMapping(*cache, batch, filter, OpenAndReadSiretMapping)
		if err != nil {
			return err
		}

		urssafRowParser.setComptes(comptes)
	}

	return err
}
