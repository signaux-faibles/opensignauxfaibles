/* eslint-disable no-use-before-define */

// Types de données de base

export type Periode = string // Date.getTime().toString()
export type Timestamp = number // Date.getTime()
export type ParPériode<T> = Record<Periode, T>

export type Departement = string

export type Siret = string
export type SiretOrSiren = Siret | string
export type CodeAPE = string

export type DataHash = string
export type ParHash<T> = Record<DataHash, T>

// Données importées pour une entreprise ou établissement

export type Scope = "etablissement" | "entreprise"

export type CompanyDataValues = {
  key: SiretOrSiren
  scope: Scope
  batch: BatchValues
}

export type CompanyDataValuesWithFlags = CompanyDataValues & IndexFlags

export type IndexFlags = {
  index: {
    algo2: boolean // pour spécifier quelles données seront à calculer puis inclure dans Features, par Reduce.algo2
  }
}

// Données importées par les parseurs, pour chaque source de données

export type BatchKey = string

export type BatchValues = Record<BatchKey, BatchValue>

export type DataType = keyof BatchValueProps // => 'reporder' | 'effectif' | 'apconso' | ...

export type BatchValue = Partial<BatchValueProps>

export type BatchValueProps = {
  reporder: ParPériode<EntréeRepOrder> // RepOrder est généré par "compact", et non importé => Usage de Periode en guise de hash d'indexation
  effectif: ParHash<EntréeEffectif>
  apconso: ParHash<EntréeApConso>
  apdemande: ParHash<EntréeApDemande>
  compte: ParHash<EntréeCompte>
  interim: ParHash<EntréeInterim>
  delai: ParHash<EntréeDelai>
  procol: ParHash<EntréeDéfaillances>
  cotisation: ParHash<EntréeCotisation>
  debit: ParHash<EntréeDebit>
  ccsf: ParHash<{ date_traitement: Date }>
  sirene: ParHash<EntréeSirene>
  sirene_ul: ParHash<EntréeSireneEntreprise>
  effectif_ent: ParHash<EntréeEffectif>
  bdf: ParHash<EntréeBdf>
  diane: ParHash<EntréeDiane>
  ellisphere: ParHash<EntréeEllisphere>
}

// Détail des types de données

export type EntréeDéfaillances = {
  action_procol: "liquidation" | "redressement" | "sauvegarde"
  stade_procol:
    | "abandon_procedure"
    | "solde_procedure"
    | "fin_procedure"
    | "plan_continuation"
    | "ouverture"
    | "inclusion_autre_procedure"
    | "cloture_insuffisance_actif"
  date_effet: Date
}

export type EntréeApConso = {
  id_conso: string
  periode: Date
  heure_consomme: number
}

export type EntréeApDemande = {
  id_demande: string
  periode: { start: Date; end: Date }
  hta: number /* Nombre total d'heures autorisées */
  motif_recours_se: number /* Cause d'activité partielle */
  effectif_entreprise?: number
  effectif?: number
  date_statut?: Date
  mta?: number
  effectif_autorise?: number
  heure_consomme?: number
  montant_consomme?: number
  effectif_consomme?: number
}

export type EntréeCompte = {
  periode: Date
  numero_compte: number
}

export type EntréeInterim = {
  periode: Date
  etp: number
}

export type EntréeRepOrder = {
  random_order: number
  periode: Date
  siret: SiretOrSiren
}

export type EntréeEffectif = {
  numero_compte: string
  periode: Date
  effectif: number
}

// Valeurs attendues par delais(), pour chaque période. (cf lib/urssaf/delai.go)
export type EntréeDelai = {
  date_creation: Date
  date_echeance: Date
  duree_delai: number // nombre de jours entre date_creation et date_echeance
  montant_echeancier: number // exprimé en euros
}

export type EntréeCotisation = {
  periode: { start: Date; end: Date }
  du: number
}

/**
 * Représente un reste à payer (dette) sur cotisation sociale ou autre.
 */
export type EntréeDebit = {
  periode: { start: Date; end: Date } // Periode pour laquelle la cotisation était attendue
  numero_ecart_negatif: number // identifiant du débit pour une période donnée (comme une sorte de numéro de facture)
  numero_historique: number // identifiant d'un remboursement (partiel ou complet) d'un débit
  numero_compte: string // identifiant URSSAF d'établissement (équivalent du SIRET)
  date_traitement: Date // Date de constatation du débit (exemple: remboursement, majoration ou autre modification du montant)
  debit_suivant: string // Hash d'un autre débit
  // le montant est ventilé entre ces deux valeurs, exprimées en euros (€):
  part_ouvriere: number
  part_patronale: number
  montant_majorations?: number // exploité par map-reduce "public", mais pas par "reduce.algo2"
}

export type EntréeSirene = {
  ape: CodeAPE
  latitude: number
  longitude: number
  departement: Departement
  raison_sociale: string
  date_creation: Date
}

export type EntréeSireneEntreprise = {
  raison_sociale: string
  nom_unite_legale: string
  nom_usage_unite_legale: string
  prenom1_unite_legale: string
  prenom2_unite_legale: string
  prenom3_unite_legale: string
  prenom4_unite_legale: string
  statut_juridique: string | null // code numérique sérialisé en chaine de caractères
  date_creation: Date
}

export type EntréeBdf = {
  arrete_bilan_bdf: Date
  annee_bdf: number
  exercice_bdf: number
  raison_sociale: string
  secteur: string
  siren: SiretOrSiren
} & EntréeBdfRatios

export type EntréeBdfRatios = {
  poids_frng: number
  taux_marge: number
  delai_fournisseur: number
  dette_fiscale: number
  financier_court_terme: number
  frais_financier: number
}

export type EntréeDiane = {
  exercice_diane: number
  arrete_bilan_diane: Date
  couverture_ca_fdr?: number | null
  interets?: number | null
  excedent_brut_d_exploitation?: number | null
  produits_financiers?: number | null
  produit_exceptionnel?: number | null
  charge_exceptionnelle?: number | null
  charges_financieres?: number | null
  ca?: number | null
  concours_bancaire_courant?: number | null
  valeur_ajoutee?: number | null
  dette_fiscale_et_sociale?: number | null
  marquee: unknown // TODO: propriété non trouvée en sortie du parseur Diane => à supprimer ?
  nom_entreprise: string
  numero_siren: SiretOrSiren
  statut_juridique: string
  procedure_collective: boolean
}

export type EntréeEllisphere = {
  siren: string
  code_groupe: string
  siren_groupe: string
  refid_groupe: string
  raison_sociale_groupe: string
  adresse_groupe: string
  personne_pou_m_groupe: string
  niveau_detention: number
  part_financiere: number
  code_filiere: string
  refid_filiere: string
  personne_pou_m_filiere: string
}
