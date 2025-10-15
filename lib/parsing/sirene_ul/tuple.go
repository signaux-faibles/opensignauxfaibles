package sireneul

import (
	"time"

	"opensignauxfaibles/lib/base"
)

// SireneUL informations sur les entreprises
type SireneUL struct {
	Siren               string     ` input:"siren"                         json:"siren,omitempty"                  sql:"siren"                    csv:"Siren"`
	RaisonSociale       string     ` input:"denominationUniteLegale"       json:"raison_sociale"                   sql:"raison_sociale"           csv:"RaisonSociale"`
	Prenom1UniteLegale  string     ` input:"prenom1UniteLegale"            json:"prenom1_unite_legale,omitempty"   sql:"prenom1_unite_legale"     csv:"Prenom1UniteLegale"`
	Prenom2UniteLegale  string     ` input:"prenom2UniteLegale"            json:"prenom2_unite_legale,omitempty"   sql:"prenom2_unite_legale"     csv:"Prenom2UniteLegale"`
	Prenom3UniteLegale  string     ` input:"prenom3UniteLegale"            json:"prenom3_unite_legale,omitempty"   sql:"prenom3_unite_legale"     csv:"Prenom3UniteLegale"`
	Prenom4UniteLegale  string     ` input:"prenom4UniteLegale"            json:"prenom4_unite_legale,omitempty"   sql:"prenom4_unite_legale"     csv:"Prenom4UniteLegale"`
	NomUniteLegale      string     ` input:"nomUniteLegale"                json:"nom_unite_legale,omitempty"       sql:"nom_unite_legale"         csv:"NomUniteLegale"`
	NomUsageUniteLegale string     ` input:"nomUsageUniteLegale"           json:"nom_usage_unite_legale,omitempty" sql:"nom_usage_unite_legale"   csv:"NomUsageUniteLegale"`
	CodeStatutJuridique string     ` input:"categorieJuridiqueUniteLegale" json:"statut_juridique"                 sql:"statut_juridique"         csv:"CodeStatutJuridique"`
	Creation            *time.Time ` input:"dateCreationUniteLegale"       json:"date_creation,omitempty"          sql:"creation"                 csv:"Creation"`
}

// Key id de l'objet
func (sireneUL SireneUL) Key() string {
	return sireneUL.Siren
}

// Type de données
func (sireneUL SireneUL) Type() base.ParserType {
	return base.SireneUl
}

// Scope de l'objet
func (sireneUL SireneUL) Scope() base.Scope {
	return base.ScopeEntreprise
}
