package urssaf

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// Procol Procédures collectives, extraction URSSAF
type Procol struct {
	Siret        string    `input:"siret"         json:"-"             csv:"siret"         sql:"siret"`
	DateEffet    time.Time `input:"dt_effet"      json:"date_effet"    csv:"date_effet"    sql:"date_effet"`
	ActionProcol string    `input:"lib_actx_stdx" json:"action_procol" csv:"action_procol" sql:"action_procol"`
	StadeProcol  string    `input:"lib_actx_stdx" json:"stade_procol"  csv:"stade_procol"  sql:"stade_procol"`
}

// Key _id de l'objet
func (procol Procol) Key() string {
	return procol.Siret
}

// Scope de l'objet
func (procol Procol) Scope() engine.Scope {
	return engine.ScopeEtablissement
}

// Type de l'objet
func (procol Procol) Type() engine.ParserType {
	return engine.Procol
}
