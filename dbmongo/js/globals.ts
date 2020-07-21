// Déclaration des fonctions globales fournies par MongoDB

declare function emit(key: unknown, value: unknown): void

// Types partagés

type ParPériode<T> = { [période: string]: T }

type BatchKey = string

type CodeAPE = string
type CodeAPENiveau2 = string
type CodeAPENiveau3 = string
type CodeAPENiveau4 = string

type CodeNAF = string

type NAF = {
  n5to1: Record<CodeAPE, CodeNAF>
  n1: Record<CodeNAF, string>
  n2: Record<CodeAPENiveau2, string>
  n3: Record<CodeAPENiveau3, string>
  n4: Record<CodeAPENiveau4, string>
  n5: Record<CodeAPE, string>
}

type Scope = "etablissement" | "entreprise"

// Données importées pour une entreprise ou établissement
type CompanyDataValues = {
  key: SiretOrSiren
  scope: Scope
  batch: BatchValues
}

type CompanyDataValuesWithFlags = CompanyDataValues & {
  index: {
    algo1: boolean
    algo2: boolean
  }
}

type BatchValues = Record<BatchKey, BatchValue>

type BatchValue = Partial<
  DonnéesRepOrder &
    DonnéesCompact &
    DonnéesEffectif &
    DonnéesApConso &
    DonnéesApDemande &
    DonnéesCompte &
    DonnéesInterim &
    DonnéesDelai &
    DonnéesDefaillances &
    DonnéesCotisation &
    DonnéesDebit &
    DonnéesCcsf &
    DonnéesSirene &
    DonnéesSireneEntreprise &
    DonnéesEffectifEntreprise &
    DonnéesBdf &
    DonnéesDiane
>

type DataType = Exclude<keyof BatchValue, "compact"> // => 'reporder' | 'effectif' | 'apconso' | ...

// Définition des types de données

type DonnéesRepOrder = {
  reporder: Record<Periode, EntréeRepOrder>
}

type DonnéesCompact = {
  compact: { delete: Partial<Record<DataType, DataHash[]>> }
}

type DonnéesEffectif = {
  effectif: Record<DataHash, EntréeEffectif>
}

type DonnéesApConso = {
  apconso: Record<DataHash, EntréeApConso>
}

type DonnéesApDemande = {
  apdemande: Record<DataHash, EntréeApDemande>
}

type DonnéesCompte = {
  compte: Record<DataHash, EntréeCompte>
}

type DonnéesInterim = {
  interim: Record<Periode, EntréeInterim>
}

type DonnéesDelai = {
  delai: Record<DataHash, EntréeDelai>
}

type DonnéesDefaillances = {
  altares: Record<DataHash, EntréeDefaillances>
  procol: Record<DataHash, EntréeDefaillances>
}

type DonnéesCotisation = {
  cotisation: Record<DataHash, EntréeCotisation>
}

type DonnéesDebit = {
  debit: Record<DataHash, EntréeDebit>
}

type DonnéesCcsf = {
  ccsf: Record<DataHash, { date_traitement: Date }>
}

type DonnéesSirene = {
  sirene: Record<DataHash, EntréeSirene>
}

type DonnéesSireneEntreprise = {
  sirene_ul: Record<DataHash, EntréeSireneEntreprise>
}

type DonnéesEffectifEntreprise = {
  effectif_ent: Record<DataHash, EntréeEffectif>
}

type DonnéesBdf = {
  bdf: Record<DataHash, EntréeBdf>
}

type DonnéesDiane = {
  diane: Record<DataHash, EntréeDiane>
}

// Détail des types de données

type AltaresCode = string

type Action = "liquidation" | "redressement" | "sauvegarde"

type Stade = "abandon_procedure" | "fin_procedure" | "plan_continuation"

type EntréeDefaillances = {
  code_evenement: AltaresCode
  action_procol: Action
  stade_procol: Stade
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
  hta: unknown
  motif_recours_se: unknown
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

type CurrentDataState = { [key: string]: Set<DataHash> }

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
}

type Departement = string

type EntréeSirene = {
  ape: CodeAPE
  lattitude: number // TODO: une fois que les données auront été migrées, corriger l'orthographe de cette propriété (--> latitude)
  longitude: unknown
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
}

type EntréeDiane = {
  exercice_diane: number
  arrete_bilan_diane: Date
  couverture_ca_fdr: number
  interets: number
  excedent_brut_d_exploitation: number
  produits_financiers: number
  produit_exceptionnel: number
  charge_exceptionnelle: number
  charges_financieres: number
}
