package apconso

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `input:"ID_DA"      sql:"id_demande"       csv:"ID"`
	Siret          string    `input:"ETAB_SIRET" sql:"siret"            csv:"Siret"`
	HeureConsommee *float64  `input:"HEURES"     sql:"heures"           csv:"HeureConsommee"`
	Montant        *float64  `input:"MONTANTS"   sql:"montant"          csv:"Montant"`
	Effectif       *int      `input:"EFFECTIFS"  sql:"effectif"         csv:"Effectif"`
	Periode        time.Time `input:"MOIS"       sql:"periode"          csv:"Periode"`
}

// Key id de l'objet
func (apconso APConso) Key() string {
	return apconso.Siret
}

// Type de données
func (apconso APConso) Type() engine.ParserType {
	return engine.Apconso
}

// Scope de l'objet
func (apconso APConso) Scope() engine.Scope {
	return engine.ScopeEtablissement
}
