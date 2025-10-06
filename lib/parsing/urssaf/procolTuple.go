package urssaf

import (
	"time"

	"opensignauxfaibles/lib/base"
)

// Procol Procédures collectives, extraction URSSAF
type Procol struct {
	Siret        string    `input:"siret"         json:"-"             csv:"siret"`
	DateEffet    time.Time `input:"dt_effet"      json:"date_effet"    csv:"date_effet"`
	ActionProcol string    `input:"lib_actx_stdx" json:"action_procol" csv:"action_procol"`
	StadeProcol  string    `input:"lib_actx_stdx" json:"stade_procol"  csv:"stade_procol"`
}

// Key _id de l'objet
func (procol Procol) Key() string {
	return procol.Siret
}

// Scope de l'objet
func (procol Procol) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (procol Procol) Type() base.ParserType {
	return base.Procol
}
