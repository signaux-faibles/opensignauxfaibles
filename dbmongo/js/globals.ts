// Types partagés

type ParPériode<T> = Record<Periode, T>

type CodeAPE = string

// Détail des types de données

type CodeAPE = string

type DataHash = string

// Détail des types de données

type EntréeDefaillances = {
  code_evenement: string
  action_procol: "liquidation" | "redressement" | "sauvegarde"
  stade_procol: "abandon_procedure" | "fin_procedure" | "plan_continuation"
  date_effet: Date
}

type EntréeApConso = {
  id_conso: string
  periode: Date
  heure_consomme: number
}

type EntréeApDemande = {
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

type EntréeCompte = {
  periode: Date
  numero_compte: number
}

type EntréeInterim = {
  periode: Date
  etp: number
}

type DataHash = string

type Periode = string // Date.toString()
type Timestamp = number // Date.getTime()

type SiretOrSiren = string

type EntréeRepOrder = {
  random_order: number
  periode: Date
  siret: SiretOrSiren
}

type EntréeEffectif = {
  numero_compte: string
  periode: Date
  effectif: number
}

// Valeurs attendues par delais(), pour chaque période. (cf dbmongo/lib/urssaf/delai.go)
type EntréeDelai = {
  date_creation: Date
  date_echeance: Date
  duree_delai: number // nombre de jours entre date_creation et date_echeance
  montant_echeancier: number // exprimé en euros
}

type DebitHash = string

type EntréeCotisation = {
  periode: { start: Date; end: Date }
  du: number
}

/**
 * Représente un reste à payer (dette) sur cotisation sociale ou autre.
 */
type EntréeDebit = {
  periode: { start: Date; end: Date } // Periode pour laquelle la cotisation était attendue
  numero_ecart_negatif: number // identifiant du débit pour une période donnée (comme une sorte de numéro de facture)
  numero_historique: number // identifiant d'un remboursement (partiel ou complet) d'un débit
  numero_compte: string // identifiant URSSAF d'établissement (équivalent du SIRET)
  date_traitement: Date // Date de constatation du débit (exemple: remboursement, majoration ou autre modification du montant)
  debit_suivant: DebitHash
  // le montant est ventilé entre ces deux valeurs, exprimées en euros (€):
  part_ouvriere: number
  part_patronale: number
  montant_majorations?: number // exploité par map-reduce "public", mais pas par "reduce.algo2"
}

type Departement = string

type EntréeSirene = {
  ape: CodeAPE
  lattitude: number // TODO: une fois que les données auront été migrées, corriger l'orthographe de cette propriété (--> latitude)
  longitude: number
  departement: Departement
  raison_sociale: string
  date_creation: Date
}

type EntréeSireneEntreprise = {
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

type EntréeBdf = {
  arrete_bilan_bdf: Date
  annee_bdf: number
  exercice_bdf: number
  raison_sociale: string
  secteur: unknown
  siren: SiretOrSiren
} & EntréeBdfRatios

type EntréeBdfRatios = {
  poids_frng: number
  taux_marge: number
  delai_fournisseur: number
  dette_fiscale: number
  financier_court_terme: number
  frais_financier: number
}

type EntréeDiane = {
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
}
