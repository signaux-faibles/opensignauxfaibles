package apconso

import (
	"time"

	"opensignauxfaibles/lib/base"
)

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `input:"ID_DA"      json:"id_conso"       sql:"id_demande"       csv:"ID"`
	Siret          string    `input:"ETAB_SIRET" json:"-"              sql:"siret"            csv:"Siret"`
	HeureConsommee *float64  `input:"HEURES"     json:"heure_consomme" sql:"heures"           csv:"HeureConsommee"`
	Montant        *float64  `input:"MONTANTS"   json:"montant"        sql:"montant"          csv:"Montant"`
	Effectif       *int      `input:"EFFECTIFS"  json:"effectif"       sql:"effectif"         csv:"Effectif"`
	Periode        time.Time `input:"MOIS"       json:"periode"        sql:"periode"          csv:"Periode"`
}

// Key id de l'objet
func (apconso APConso) Key() string {
	return apconso.Siret
}

// Type de données
func (apconso APConso) Type() base.ParserType {
	return base.Apconso
}

// Scope de l'objet
func (apconso APConso) Scope() string {
	return "etablissement"
}
