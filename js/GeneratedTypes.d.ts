/**
 * This file was automatically generated by generate-types.sh.
 *
 * DO NOT MODIFY IT BY HAND.
 *
 * Instead:
 * - modify the validation/*.schema.json files;
 * - then, run generate-types.sh to regenerate this file.
 */

/**
 * Champs importés par le parseur lib/apconso/main.go de sfdata.
 */
export interface EntréeApConso {
  id_conso: string
  heure_consomme: number
  montant?: number
  effectif?: number
  periode: Date
}
/**
 * Champs importés par le parseur lib/apdemande/main.go de sfdata.
 */
export interface EntréeApDemande {
  id_demande: string
  periode: {
    start: Date
    end: Date
  }
  /**
   * Nombre total d'heures autorisées
   */
  hta: number
  /**
   * Cause d'activité partielle
   */
  motif_recours_se: number
  effectif_entreprise?: number
  effectif?: number
  date_statut?: Date
  mta?: number
  effectif_autorise?: number
  heure_consommee?: number
  montant_consommee?: number
  effectif_consomme?: number
}
/**
 * Note: CE SCHEMA EST INCOMPLET POUR L'INSTANT. Cf https://github.com/signaux-faibles/opensignauxfaibles/pull/143
 */
export interface EntréeBdf {
  siren: string
}
/**
 * Champs importés par le parseur lib/urssaf/ccsf.go de sfdata.
 */
export interface EntréeCcsf {
  /**
   * Date de début de la procédure CCSF
   */
  date_traitement: Date
  /**
   * TODO: choisir un type plus précis
   */
  stade: string
  /**
   * TODO: choisir un type plus précis
   */
  action: string
}
/**
 * Champs importés par le parseur lib/urssaf/compte.go de sfdata.
 */
export interface EntréeCompte {
  /**
   * Date à laquelle cet établissement est associé à ce numéro de compte URSSAF.
   */
  periode: Date
  /**
   * Numéro SIRET de l'établissement. Les numéros avec des Lettres sont des sirets provisoires.
   */
  siret: string
  /**
   * Compte administratif URSSAF.
   */
  numero_compte: string
}
/**
 * Champs importés par le parseur lib/urssaf/delai.go de sfdata.
 */
export interface EntréeDelai {
  /**
   * Compte administratif URSSAF.
   */
  numero_compte: string
  /**
   * Le numéro de structure est l'identifiant d'un dossier contentieux.
   */
  numero_contentieux: string
  /**
   * Date de création du délai.
   */
  date_creation: Date
  /**
   * Date d'échéance du délai.
   */
  date_echeance: Date
  /**
   * Durée du délai en jours: nombre de jours entre date_creation et date_echeance.
   */
  duree_delai: number
  /**
   * Raison sociale de l'établissement.
   */
  denomination: string
  /**
   * Délai inférieur ou supérieur à 6 mois ? Modalités INF et SUP.
   */
  indic_6m: string
  /**
   * Année de création du délai.
   */
  annee_creation: number
  /**
   * Montant global de l'échéancier, en euros.
   */
  montant_echeancier: number
  /**
   * Code externe du stade.
   */
  stade: string
  /**
   * Code externe de l'action.
   */
  action: string
}
/**
 * Champs importés par le parseur lib/urssaf/effectif.go de sfdata.
 */
export interface EntréeEffectif {
  /**
   * Compte administratif URSSAF.
   */
  numero_compte: string
  periode: Date
  /**
   * Nombre de personnes employées par l'établissement.
   */
  effectif: number
}
/**
 * Champs importés par le parseur lib/urssaf/procol.go de sfdata.
 */
export interface EntréeDéfaillances {
  /**
   * Nature de la procédure de défaillance.
   */
  action_procol: "liquidation" | "redressement" | "sauvegarde"
  /**
   * Evénement survenu dans le cadre de cette procédure.
   */
  stade_procol:
    | "abandon_procedure"
    | "solde_procedure"
    | "fin_procedure"
    | "plan_continuation"
    | "ouverture"
    | "inclusion_autre_procedure"
    | "cloture_insuffisance_actif"
  /**
   * Date effet de la procédure collective.
   */
  date_effet: Date
}
