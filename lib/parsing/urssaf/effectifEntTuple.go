package urssaf

import (
	"time"

	"opensignauxfaibles/lib/base"
)

// EffectifEnt Urssaf
type EffectifEnt struct {
	Siren       string    `input:"siren" json:"-"        sql:"siren"        csv:"siren"`
	Periode     time.Time `              json:"periode"  sql:"periode"      csv:"période"`
	EffectifEnt int       `              json:"effectif" sql:"effectif_ent" csv:"effectif_entreprise"`
}

// Key _id de l'objet
func (effectifEnt EffectifEnt) Key() string {
	return effectifEnt.Siren
}

// Scope de l'objet
func (effectifEnt EffectifEnt) Scope() string {
	return "entreprise"
}

// Type de l'objet
func (effectifEnt EffectifEnt) Type() base.ParserType {
	return base.EffectifEnt
}
