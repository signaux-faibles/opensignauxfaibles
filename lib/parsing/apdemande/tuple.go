package apdemande

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string    `input:"ID_DA"              json:"id_demande"          sql:"id_demande"           csv:"id_demande"`
	Siret              string    `input:"ETAB_SIRET"         json:"-"                   sql:"siret"                csv:"siret"`
	EffectifEntreprise *int      `input:"EFF_ENT"            json:"effectif_entreprise"                            csv:"effectif_entreprise"`
	Effectif           *int      `input:"EFF_ETAB"           json:"effectif"                                       csv:"effectif"`
	DateStatut         time.Time `input:"DATE_STATUT"        json:"date_statut"         sql:"date_statut"          csv:"date_statut"`
	PeriodeDebut       time.Time `input:"DATE_DEB"           json:"periode_debut"       sql:"periode_debut"        csv:"période_début"`
	PeriodeFin         time.Time `input:"DATE_FIN"           json:"periode_fin"         sql:"periode_fin"          csv:"période_fin"`
	HTA                *float64  `input:"HTA"                json:"hta"                 sql:"heures"               csv:"heures_autorisées"`
	MTA                *float64  `                           json:"mta"                 sql:"montant"              csv:"montants_autorisés"`
	EffectifAutorise   *int      `input:"EFF_AUTO"           json:"effectif_autorise"   sql:"effectif"             csv:"effectif_autorisé"`
	MotifRecoursSE     *int      `input:"MOTIF_RECOURS_SE"   json:"motif_recours_se"    sql:"motif_recours"        csv:"motif_recours_se"`
	HeureConsommee     *float64  `input:"S_HEURE_CONSOM_TOT" json:"heures_consommees"                              csv:"heure_consommee"`
	MontantConsomme    *float64  `                           json:"montant_consomme"                               csv:"montant_consomme"`
	EffectifConsomme   *int      `input:"S_HEURE_CONSOM_TOT" json:"effectif_consomme"                              csv:"effectif_consomme"`
	Perimetre          *int      `input:"PERIMETRE_AP"       json:"perimetre"                                      csv:"perimetre"`
}

// Key id de l'objet
func (apdemande APDemande) Key() string {
	return apdemande.Siret
}

// Type de données
func (apdemande APDemande) Type() engine.ParserType {
	return engine.Apdemande
}

// Scope de l'objet
func (apdemande APDemande) Scope() engine.Scope {
	return engine.ScopeEtablissement
}
