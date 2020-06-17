// Déclaration des fonctions globales fournies par MongoDB
declare function emit(key: unknown, value: object): void

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
  reporder: { [periode: string]: RepOrder }
  compact: { delete: { [dataType: string]: DataHash[] } }
  effectif: { [dataHash: string]: Effectif }
  apconso: { [key: string]: any } // TODO: définir type plus précisément
  apdemande: { [key: string]: any } // TODO: définir type plus précisément
  compte: Record<DataHash, Compte>
  interim: Record<Periode, Interim>
}

type BatchDataType = keyof BatchValue

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
  periode: Periode
  effectif: number
}

type CurrentDataState = { [key: string]: Set<DataHash> }
