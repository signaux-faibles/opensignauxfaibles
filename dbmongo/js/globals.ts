// Déclaration des fonctions globales fournies par MongoDB
declare function emit(key: string, value: unknown): void

// Paramètres globaux utilisés par "compact"
/* eslint-disable @typescript-eslint/no-unused-vars */
let batches: string[]
let batchKey: string
let serie_periode: Date[]
let types: string[]
let completeTypes: { [key: string]: string[] }
/* eslint-enable @typescript-eslint/no-unused-vars */

// Types partagés

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
  apconso?: { [key: string]: any } // TODO: définir type plus précisément
  apdemande?: { [key: string]: any } // TODO: définir type plus précisément
}

type DataHash = string

type Periode = Date

type SiretOrSiren = string

type RepOrder = {
  random_order: number
  periode: Periode
  siret: SiretOrSiren
}

type Effectif = {
  numero_compte: string
  periode: Periode
  effectif: number
}

type CurrentDataState = { [key: string]: Set<DataHash> }
