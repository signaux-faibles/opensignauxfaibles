package sirenehisto

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// SireneHisto is one event in the history of changes to "établissements" SIRENE
// We keep only changes related to being active or closed
type SireneHisto struct {
	Siret     string     `input:"siret"                          json:"siret"                    sql:"siret"`
	DateDebut *time.Time `input:"dateDebut"                      json:"date_debut,omitempty"     sql:"date_debut"`
	DateFin   *time.Time `input:"dateFin"                        json:"date_fin,omitempty"       sql:"date_fin"`
	EstActif  bool       `input:"etatAdministratifEtablissement" json:"est_actif,omitempty"      sql:"est_actif"`
	// Is the change related to a the "établissement" closing (EstActif true -> false) or (re)opening (EstActif false -> true) ?
	ChangementStatutActif bool ` input:"changementEtatAdministratifEtablissement" json:"changement_statut_actif"  sql:"changement_statut_actif"`
}

// Key id de l'objet
func (sh SireneHisto) Key() string {
	return sh.Siret
}

// Type de données
func (sh SireneHisto) Type() engine.ParserType {
	return engine.SireneHisto
}

// Scope de l'objet
func (sh SireneHisto) Scope() engine.Scope {
	return engine.ScopeEtablissement
}
