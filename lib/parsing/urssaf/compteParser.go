package urssaf

import (
	"io"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/misc"
)

type comptesParser struct {
	periodes []time.Time
	mapping  engine.Comptes
}

func NewParserComptes() *comptesParser {
	return &comptesParser{}
}

func (parser *comptesParser) Type() base.ParserType {
	return base.AdminUrssaf
}

func (parser *comptesParser) Init(cache *engine.Cache, filter engine.SirenFilter, batch *base.AdminBatch) (err error) {
	if len(batch.Files["admin_urssaf"]) > 0 {
		parser.periodes = misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin)
		parser.mapping, err = engine.GetCompteSiretMapping(*cache, batch, filter, engine.OpenAndReadSiretMapping)
	}
	return err
}

func (parser *comptesParser) Open(filePath base.BatchFile) error {
	// Ce parseur produit des tuples à partir des mappings compte<->siret déjà
	// parsés par engine.GetCompteSiretMapping(). => pas de fichier à ouvrir.
	return nil
}

func (parser *comptesParser) Close() error {
	return nil
}

func (parser *comptesParser) ReadNext(res *engine.ParsedLineResult) error {
	// First, we sort the mapping entries by account number, to make sure that
	// tuples are always processed in the same order, and therefore that errors
	// (e.g. "siret invalide") are reported at consistent Cycle/line numbers.
	// cf https://github.com/signaux-faibles/opensignauxfaibles/pull/225#issuecomment-720594272
	accounts := parser.mapping.GetSortedKeys()
	for accountIndex := range accounts {

		account := accounts[accountIndex]
		for _, p := range parser.periodes {
			var err error
			compte := Compte{}
			compte.NumeroCompte = account
			compte.Periode = p
			compte.Siret, err = engine.GetSiretFromComptesMapping(account, &p, parser.mapping)
			res.AddRegularError(err)
			res.AddTuple(compte)
		}
	}
	return io.EOF
}
