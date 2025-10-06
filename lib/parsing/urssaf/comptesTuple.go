package urssaf

import (
	"time"

	"opensignauxfaibles/lib/base"
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
func (compte Compte) Type() base.ParserType {
	return base.AdminUrssaf
}
