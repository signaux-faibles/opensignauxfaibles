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

// ParseCompteFile permet de lancer le parsing du fichier demandé.
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
	accounts := mapKeys(*mapping)
	accountIndex := 0
	for {
		parsedLine := base.ParsedLineResult{}
		if accountIndex >= len(*mapping) {
			close(parsedLineChan) // EOF
			break
		}
		account := accounts[accountIndex]
		accountIndex++
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
}

func mapKeys(mymap marshal.Comptes) []string {
	keys := make([]string, len(mymap))
	i := 0
	for k := range mymap {
		keys[i] = k
		i++
	}
	return keys
}
