package urssaf

import (
	"opensignauxfaibles/lib/base"
	"time"
)

// Delai tuple fichier urssaf
// Implémente engine.Tuple
type Delai struct {
	Siret             string    `                                    json:"-"                  sql:"siret"                csv:"siret"`
	NumeroCompte      string    `input:"Numero_compte_externe"       json:"numero_compte"                                 csv:"numéro_compte"`
	NumeroContentieux string    `input:"Numero_structure"            json:"numero_contentieux"                            csv:"numéro_contentieux"`
	DateCreation      time.Time `input:"Date_creation"               json:"date_creation"      sql:"date_creation"        csv:"date_création"`
	DateEcheance      time.Time `input:"Date_echeance"               json:"date_echeance"      sql:"date_echeance"        csv:"date_échéance"`
	DureeDelai        *int      `input:"Duree_delai"                 json:"duree_delai"        sql:"duree_delai"          csv:"durée_délai"`
	Denomination      string    `input:"Denomination_premiere_ligne" json:"denomination"                                  csv:"dénomination"`
	Indic6m           string    `input:"Indic_6M"                    json:"indic_6m"                                      csv:"indic_6mois"`
	AnneeCreation     *int      `input:"Annee_creation"              json:"annee_creation"                                csv:"année_création"`
	MontantEcheancier *float64  `input:"Montant_global_echeancier"   json:"montant_echeancier" sql:"montant_echeancier"   csv:"montant_échéancier"`
	Stade             string    `input:"Code_externe_stade"          json:"stade"              sql:"stade"                csv:"stade"`
	Action            string    `input:"Code_externe_action"         json:"action"             sql:"action"               csv:"action"`
}

// Key _id de l'objet
func (delai Delai) Key() string {
	return delai.Siret
}

// Scope de l'objet
func (delai Delai) Scope() base.Scope {
	return base.ScopeEtablissement
}

// Type de l'objet
func (delai Delai) Type() base.ParserType {
	return base.Delai
}
