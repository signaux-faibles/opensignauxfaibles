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
    DonnéesSireneUL &
    DonnéesEffectifEntreprise &
    DonnéesBdf &
    DonnéesDiane
>

type BatchDataType = Exclude<keyof BatchValue, "compact"> // => 'reporder' | 'effectif' | 'apconso' | ...

// Définition des types de données

type DonnéesRepOrder = {
  reporder: Record<Periode, EntréeRepOrder>
}

type DonnéesCompact = {
  compact: { delete: { [dataType: string]: DataHash[] } } // TODO: utiliser un type Record<~BatchDataType, DataHash[]>
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
  cotisation: Record<string, EntréeCotisation> // TODO: utiliser un type plus précis que string
}

type DonnéesDebit = {
  debit: Record<string, EntréeDebit> // TODO: utiliser un type plus précis que string
}

type DonnéesCcsf = {
  ccsf: Record<DataHash, { date_traitement: Date }>
}

type DonnéesSirene = {
  sirene: Record<string, EntréeSirene> // TODO: utiliser un type plus précis que string
}

type DonnéesSireneUL = {
  sirene_ul: Record<string, EntréeSireneUL> // TODO: utiliser un type plus précis que string
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

type Periode = string // Date

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

type EntréeDebit = {
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

type Departement = string

type EntréeSirene = {
  ape: CodeAPE
  lattitude: number // TODO: une fois que les données auront été migrées, corriger l'orthographe de cette propriété (--> latitude)
  longitude: unknown
  departement: Departement
  raison_sociale: string
  date_creation: Date
}

type EntréeSireneUL = {
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
