package urssaf

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/misc"
	"time"
)

// Compte tuple fichier ursaff
type Compte struct {
	key          string    `hash:"-"`
	Siret        string    `json:"siret" bson:"siret"`
	NumeroCompte string    `json:"numero_compte" bson:"numero_compte"`
	Periode      time.Time `json:"periode" bson:"periode"`
}

// Key _id de l'objet
func (compte Compte) Key() string {
	return compte.key
}

// Scope de l'objet
func (compte Compte) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (compte Compte) Type() string {
	return "compte"
}

func parseCompte(batch engine.AdminBatch, mapping Comptes) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	//	event := engine.Event{
	//	Code:    "compteParser",
	//	Channel: eventChannel,
	//}

	go func() {
		periode_init := batch.Params.DateDebut
		periodes := misc.GenereSeriePeriode(periode_init, time.Now()) //[]time.Time
		for c := range mapping {
			for _, p := range periodes {
				compte := Compte{}
				compte.NumeroCompte = c
				compte.Periode = p
				var err error
				compte.Siret, err = mapping.GetSiret(c, p)
				if err == nil {
          compte.key = compte.Siret
					outputChannel <- compte
				}
			}
		}
		close(outputChannel)
		close(eventChannel)
	}()
	return outputChannel, eventChannel
}
