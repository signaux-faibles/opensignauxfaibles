package urssaf

import (
	"time"

	"opensignauxfaibles/lib/base"
)

// Cotisation Objet cotisation
type Cotisation struct {
	Siret        string    `                   json:"-"             sql:"siret"            csv:"siret"`
	NumeroCompte string    `input:"Compte"     json:"numero_compte"                      csv:"numéro_compte"`
	PeriodeDebut time.Time `input:"periode"    json:"periode_debut" sql:"periode_debut"  csv:"période_début"`
	PeriodeFin   time.Time `input:"periode"    json:"periode_fin"   sql:"periode_fin"    csv:"période_fin"`
	Encaisse     *float64  `input:"enc_direct" json:"encaisse"                           csv:"encaissé"`
	Du           *float64  `input:"cotis_due"  json:"du"            sql:"du"             csv:"du"`
}

// Key _id de l'objet
func (cotisation Cotisation) Key() string {
	return cotisation.Siret
}

// Scope de l'objet
func (cotisation Cotisation) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (cotisation Cotisation) Type() base.ParserType {
	return base.Cotisation
}
