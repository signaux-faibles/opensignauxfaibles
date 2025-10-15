package effectif

import (
	"time"

	"opensignauxfaibles/lib/base"
)

// Effectif Urssaf
type Effectif struct {
	Siret        string    `input:"siret"  json:"-"             sql:"siret"          csv:"siret"`
	NumeroCompte string    `input:"compte" json:"numero_compte"                      csv:"compte"`
	Periode      time.Time `               json:"periode"       sql:"periode"        csv:"période"`
	Effectif     int       `               json:"effectif"      sql:"effectif"       csv:"effectif"`
}

// Key _id de l'objet
func (effectif Effectif) Key() string {
	return effectif.Siret
}

// Scope de l'objet
func (effectif Effectif) Scope() base.Scope {
	return base.ScopeEtablissement
}

// Type de l'objet
func (effectif Effectif) Type() base.ParserType {
	return base.Effectif
}
