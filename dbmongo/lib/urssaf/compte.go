package urssaf

import (
	"time"

	"github.com/signaux-faibles/gournal"
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

// ParseCompteFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseCompteFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	if len(batch.Files["admin_urssaf"]) > 0 {
		periodes := misc.GenereSeriePeriode(batch.Params.DateDebut, time.Now()) //[]time.Time
		mapping, err := marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
		tracker.Add(err)
		for c := range mapping {
			for _, p := range periodes {
				var err error
				compte := Compte{}
				compte.NumeroCompte = c
				compte.Periode = p
				compte.Siret, err = marshal.GetSiret(c, &p, *cache, batch)
				tracker.Add(base.NewCriticError(err, "erreur"))
				outputChannel <- compte
			}
			tracker.Next()
		}
	}
}
