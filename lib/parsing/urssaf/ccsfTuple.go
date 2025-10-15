package urssaf

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// CCSF information urssaf ccsf
type CCSF struct {
	key            string    `                                                   csv:"siret"`
	NumeroCompte   string    `input:"Compte"              json:"-"               csv:"numéro_compte"`
	DateTraitement time.Time `input:"Date_de_traitement"  json:"date_traitement" csv:"date_traitement"`
	Stade          string    `input:"Code_externe_stade"  json:"stade"           csv:"stade"`
	Action         string    `input:"Code_externe_action" json:"action"          csv:"action"`
}

// Key _id de l'objet
func (ccsf CCSF) Key() string {
	return ccsf.key
}

// Scope de l'objet
func (ccsf CCSF) Scope() engine.Scope {
	return engine.ScopeEtablissement
}

// Type de l'objet
func (ccsf CCSF) Type() engine.ParserType {
	return engine.Ccsf
}
