package effectif

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// EffectifEnt Urssaf
type EffectifEnt struct {
	Siren       string    `input:"siren" sql:"siren"        csv:"siren"`
	Periode     time.Time `              sql:"periode"      csv:"période"`
	EffectifEnt int       `              sql:"effectif_ent" csv:"effectif_entreprise"`
}

// Key _id de l'objet
func (effectifEnt EffectifEnt) Key() string {
	return effectifEnt.Siren
}

// Scope de l'objet
func (effectifEnt EffectifEnt) Scope() engine.Scope {
	return engine.ScopeEntreprise
}

// Type de l'objet
func (effectifEnt EffectifEnt) Type() engine.ParserType {
	return engine.EffectifEnt
}
