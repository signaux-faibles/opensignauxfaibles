package effectif

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// Effectif Urssaf
type Effectif struct {
	Siret        string    `input:"siret"  sql:"siret"          csv:"siret"`
	NumeroCompte string    `input:"compte"                      csv:"compte"`
	Periode      time.Time `               sql:"periode"        csv:"période"`
	Effectif     int       `               sql:"effectif"       csv:"effectif"`
}

// Key _id de l'objet
func (effectif Effectif) Key() string {
	return effectif.Siret
}

// Scope de l'objet
func (effectif Effectif) Scope() engine.Scope {
	return engine.ScopeEtablissement
}

// Type de l'objet
func (effectif Effectif) Type() engine.ParserType {
	return engine.Effectif
}
