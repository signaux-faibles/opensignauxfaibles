package urssaf

import (
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
)

// Compte tuple fichier ursaff
type Compte struct {
	Siret        string    `json:"siret" bson:"siret"`
	NumeroCompte string    `json:"numero_compte" bson:"numero_compte"`
	Periode      time.Time `json:"periode" bson:"periode"`
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
	return "compte"
}

// ParserCompte expose le parseur et le type de fichier qu'il supporte.
var ParserCompte = marshal.Parser{FileType: "admin_urssaf", FileParser: ParseCompteFile}

// ParseCompteFile permet de générer des tuples à partir des mappings
// compte<->siret déjà parsés par marshal.GetCompteSiretMapping().
func ParseCompteFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) marshal.OpenFileResult {
	var err error
	var periodes []time.Time
	var mapping marshal.Comptes
	closeFct := func() error { return nil }
	if len(batch.Files["admin_urssaf"]) > 0 {
		periodes = misc.GenereSeriePeriode(batch.Params.DateDebut, time.Now()) //[]time.Time
		mapping, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	}
	return marshal.OpenFileResult{
		Error: err,
		ParseLines: func(parsedLineChan chan base.ParsedLineResult) {
			parseCompteLines(periodes, &mapping, parsedLineChan)
		},
		Close: closeFct,
	}
}

func parseCompteLines(periodes []time.Time, mapping *marshal.Comptes, parsedLineChan chan base.ParsedLineResult) {
	// First, we sort the mapping entries by account number, to make sure that
	// tuples are always processed in the same order, and therefore that errors
	// (e.g. "siret invalide") are reported at consistent Cycle/line numbers.
	// cf https://github.com/signaux-faibles/opensignauxfaibles/pull/225#issuecomment-720594272
	accounts := mapping.GetSortedKeys()
	for accountIndex := range accounts {
		parsedLine := base.ParsedLineResult{}
		account := accounts[accountIndex]
		for _, p := range periodes {
			var err error
			compte := Compte{}
			compte.NumeroCompte = account
			compte.Periode = p
			compte.Siret, err = marshal.GetSiretFromComptesMapping(account, &p, *mapping)
			parsedLine.AddError(base.NewCriticError(err, "erreur"))
			parsedLine.AddTuple(compte)
		}
		parsedLineChan <- parsedLine
	}
	close(parsedLineChan) // EOF
}
