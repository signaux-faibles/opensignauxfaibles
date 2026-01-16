package apdemande

import (
	"opensignauxfaibles/lib/engine"
	"time"
)

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string    `input:"ID_DA"              sql:"id_demande"           csv:"id_demande"`
	Siret              string    `input:"ETAB_SIRET"         sql:"siret"                csv:"siret"`
	EffectifEntreprise *int      `input:"EFF_ENT"                                       csv:"effectif_entreprise"`
	Effectif           *int      `input:"EFF_ETAB"                                      csv:"effectif"`
	DateStatut         time.Time `input:"DATE_STATUT"        sql:"date_statut"          csv:"date_statut"`
	PeriodeDebut       time.Time `input:"DATE_DEB"           sql:"periode_debut"        csv:"période_début"`
	PeriodeFin         time.Time `input:"DATE_FIN"           sql:"periode_fin"          csv:"période_fin"`
	HTA                *float64  `input:"HTA"                sql:"heures"               csv:"heures_autorisées"`
	MTA                *float64  `                           sql:"montant"              csv:"montants_autorisés"`
	EffectifAutorise   *int      `input:"EFF_AUTO"           sql:"effectif"             csv:"effectif_autorisé"`
	MotifRecoursSE     *int      `input:"MOTIF_RECOURS_SE"   sql:"motif_recours"        csv:"motif_recours_se"`
	HeureConsommee     *float64  `input:"S_HEURE_CONSOM_TOT"                            csv:"heure_consommee"`
	MontantConsomme    *float64  `                                                      csv:"montant_consomme"`
	EffectifConsomme   *int      `input:"S_HEURE_CONSOM_TOT"                            csv:"effectif_consomme"`
	Perimetre          *int      `input:"PERIMETRE_AP"                                  csv:"perimetre"`
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
