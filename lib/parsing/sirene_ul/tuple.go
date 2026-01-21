package sireneul

import (
	"opensignauxfaibles/lib/engine"
	"time"
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
	CategorieJuridique  string     ` input:"categorieJuridiqueUniteLegale" json:"categorie_juridique"              sql:"categorie_juridique"         csv:"CodeStatutJuridique"`
	Creation            *time.Time ` input:"dateCreationUniteLegale"       json:"date_creation,omitempty"          sql:"creation"                 csv:"Creation"`
	APE                 string     `                                       json:"ape,omitempty"                    sql:"activite_principale"      csv:"APE"`
	EstActif            bool       ` input:"etatAdministratifUniteLegale"  json:"est_actif"                        sql:"est_actif"                csv:"EstActif"`
}

// Key id de l'objet
func (sireneUL SireneUL) Key() string {
	return sireneUL.Siren
}

// Type de données
func (sireneUL SireneUL) Type() engine.ParserType {
	return engine.SireneUl
}

// Scope de l'objet
func (sireneUL SireneUL) Scope() engine.Scope {
	return engine.ScopeEntreprise
}
