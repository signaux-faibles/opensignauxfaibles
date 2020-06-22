// Déclaration des fonctions globales fournies par MongoDB
declare function emit(key: unknown, value: unknown): void

// Paramètres globaux utilisés par "compact"
/* eslint-disable @typescript-eslint/no-unused-vars */
let batches: BatchKey[]
let batchKey: BatchKey
let serie_periode: Date[]
let types: string[]
let completeTypes: { [key: string]: string[] }
/* eslint-enable @typescript-eslint/no-unused-vars */

// Types partagés

type BatchKey = string

type CodeAPE = string

type CodeNAF = string

type Scope = "etablissement" | "entreprise"

type ReduceIndexFlags = {
  algo1: boolean
  algo2: boolean
}

type BatchValues = { [batchKey: string]: BatchValue }

type CompanyDataValues = {
  key: SiretOrSiren
  scope: Scope
  batch: BatchValues
}

type CompanyDataValuesWithFlags = CompanyDataValues & {
  index: ReduceIndexFlags
}

type BatchValue = {
  reporder?: { [periode: string]: RepOrder }
  compact?: { delete: { [dataType: string]: DataHash[] } }
  effectif?: { [dataHash: string]: Effectif }
  apconso?: Record<DataHash, ApConso>
  apdemande?: Record<DataHash, ApDemande>
  compte?: Record<DataHash, Compte>
  interim?: Record<Periode, Interim>
  delai?: Record<DataHash, Delai>
} & Partial<DonnéesDefaillances> &
  Partial<DonnéesCotisationsDettes> &
  Partial<DonnéesCcsf> &
  Partial<DonnéesSirene> &
  Partial<DonnéesSireneUL> &
  Partial<DonnéesEffectifEntreprise> &
  Partial<DonnéesBdf> &
  Partial<DonnéesDiane>

type BatchDataType = keyof BatchValue

type AltaresCode = string

type Action = "liquidation" | "redressement" | "sauvegarde"

type Stade = "abandon_procedure" | "fin_procedure" | "plan_continuation"

type Événement = {
  code_evenement: AltaresCode
  action_procol: Action
  stade_procol: Stade
  date_effet: Date
}

type DonnéesDefaillances = {
  altares: Record<DataHash, Événement>
  procol: Record<DataHash, Événement>
}

type ApConso = {
  id_conso: string
  periode: Date
  heure_consomme: number
}

type ApDemande = {
  id_demande: string
  periode: { start: Date; end: Date }
  hta: unknown
  motif_recours_se: unknown
}

type Compte = {
  periode: Date
  numero_compte: number
}

type Interim = {
  periode: Date
  etp: number
}

type DataHash = string

type Periode = string // Date

type SiretOrSiren = string

type RepOrder = {
  random_order: number
  periode: Date
  siret: SiretOrSiren
}

type Effectif = {
  numero_compte: string
  periode: Date
  effectif: number
}

// Valeurs attendues par delais(), pour chaque période. (cf dbmongo/lib/urssaf/delai.go)
type Delai = {
  date_creation: Date
  date_echeance: Date
  duree_delai: number // nombre de jours entre date_creation et date_echeance
  montant_echeancier: number // exprimé en euros
}

type CurrentDataState = { [key: string]: Set<DataHash> }

type DebitHash = string

type Cotisation = {
  periode: { start: Date; end: Date }
  du: number
}

type Debit = {
  periode: { start: Date; end: Date }
  numero_ecart_negatif: unknown
  numero_compte: unknown
  numero_historique: number
  date_traitement: Date
  debit_suivant: DebitHash
  part_ouvriere: number
  part_patronale: number
  montant_majorations: number
}

type DonnéesCotisationsDettes = {
  cotisation: Record<string, Cotisation>
  debit: Record<string, Debit>
}

type DonnéesCcsf = {
  ccsf: Record<DataHash, { date_traitement: Date }>
}

type Departement = string

type Sirene = {
  ape: CodeAPE
  lattitude: unknown // TODO ⚠️ typo ?
  longitude: unknown
  departement: Departement
  raison_sociale: unknown
  date_creation: Date
}

type DonnéesSirene = {
  sirene: Record<string, Sirene>
}

type SireneUL = {
  raison_sociale: string
  nom_unite_legale: string
  nom_usage_unite_legale: string
  prenom1_unite_legale: string
  prenom2_unite_legale: string
  prenom3_unite_legale: string
  prenom4_unite_legale: string
  statut_juridique: unknown
  date_creation: Date
}

type DonnéesSireneUL = {
  sirene_ul: Record<string, SireneUL>
}

type EffectifEntreprise = Record<DataHash, Effectif>

type DonnéesEffectifEntreprise = {
  effectif_ent: EffectifEntreprise
}

type EntréeBdf = {
  arrete_bilan_bdf: Date
  annee_bdf: number
  exercice_bdf: number
}

type DonnéesBdf = {
  bdf: Record<DataHash, EntréeBdf>
}

type EntréeDiane = {
  arrete_bilan_diane: Date
  couverture_ca_fdr: number
  interets: number
  excedent_brut_d_exploitation: number
  produits_financiers: number
  produit_exceptionnel: number
  charge_exceptionnelle: number
  charges_financieres: number
}

type DonnéesDiane = {
  diane: Record<DataHash, EntréeDiane>
}
