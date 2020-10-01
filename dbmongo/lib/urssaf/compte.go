package urssaf

import (
	"time"

	"github.com/signaux-faibles/gournal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
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

// ParserCompte retourne les comptes lus depuis un fichier admin_urssaf.
func ParserCompte(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)
	go func() {

		defer close(outputChannel)
		defer close(eventChannel)
		if len(batch.Files["admin_urssaf"]) > 0 {
			event := marshal.Event{
				Code:    "compteParser",
				Channel: eventChannel,
			}
			tracker := gournal.NewTracker(
				map[string]string{"path": "Admin_urssaf", "batchKey": batch.ID.Key},
				engine.TrackerReports)

			periode_init := batch.Params.DateDebut
			periodes := misc.GenereSeriePeriode(periode_init, time.Now()) //[]time.Time
			event.Info("Comptes urssaf : traitement")
			mapping, err := marshal.GetCompteSiretMapping(cache, batch, marshal.OpenAndReadSiretMapping)
			tracker.Add(err)
			for c := range mapping {
				for _, p := range periodes {
					compte := Compte{}
					compte.NumeroCompte = c
					compte.Periode = p
					var err error
					compte.Siret, err = marshal.GetSiret(c, &p, cache, batch)
					tracker.Add(base.NewCriticError(err, "erreur"))

					outputChannel <- compte
				}
				tracker.Next()
			}
			event.Info(tracker.Report("abstract"))
		}
	}()
	return outputChannel, eventChannel
}
