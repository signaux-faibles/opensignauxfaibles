package urssaf

import (
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/misc"
)

// Compte tuple fichier ursaff
type Compte struct {
	Siret        string    `json:"siret"         csv:"siret"`
	NumeroCompte string    `json:"numero_compte" csv:"numéro_compte"`
	Periode      time.Time `json:"periode"       csv:"période"`
}

// Key _id de l'objet
func (compte Compte) Key() string {
	return compte.Siret
}

// Scope de l'objet
func (compte Compte) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (compte Compte) Type() string {
	return "admin_urssaf"
}

// ParserCompte fournit une instance utilisable par ParseFilesFromBatch.
var ParserCompte = &comptesParser{}

type comptesParser struct {
	periodes []time.Time
	mapping  marshal.Comptes
}

func (parser *comptesParser) Type() string {
	return "admin_urssaf"
}

func (parser *comptesParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	if len(batch.Files["admin_urssaf"]) > 0 {
		parser.periodes = misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin)
		parser.mapping, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	}
	return err
}

func (parser *comptesParser) Open(filePath base.BatchFile) error {
	// Ce parseur produit des tuples à partir des mappings compte<->siret déjà
	// parsés par marshal.GetCompteSiretMapping(). => pas de fichier à ouvrir.
	return nil
}

func (parser *comptesParser) Close() error {
	return nil
}

func (parser *comptesParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	// First, we sort the mapping entries by account number, to make sure that
	// tuples are always processed in the same order, and therefore that errors
	// (e.g. "siret invalide") are reported at consistent Cycle/line numbers.
	// cf https://github.com/signaux-faibles/opensignauxfaibles/pull/225#issuecomment-720594272
	accounts := parser.mapping.GetSortedKeys()
	for accountIndex := range accounts {
		parsedLine := marshal.ParsedLineResult{}
		account := accounts[accountIndex]
		for _, p := range parser.periodes {
			var err error
			compte := Compte{}
			compte.NumeroCompte = account
			compte.Periode = p
			compte.Siret, err = marshal.GetSiretFromComptesMapping(account, &p, parser.mapping)
			parsedLine.AddRegularError(err)
			parsedLine.AddTuple(compte)
		}
		parsedLineChan <- parsedLine
	}
	close(parsedLineChan) // EOF
}
