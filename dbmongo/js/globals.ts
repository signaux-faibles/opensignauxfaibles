// Déclaration des fonctions globales fournies par MongoDB
declare function emit(key: any, value: any): void
declare function print(...any): void

// Déclaration des fonctions globales fournies par JSC
declare function debug(string) // supported by jsc, to print in stdout

// Déclaration de variables globales
/* eslint-disable @typescript-eslint/no-unused-vars */
let f: {
  [key: string]: Function
}

// Paramètres globaux utilisés par "compact"
let batches: string[]
let batchKey: string
let serie_periode: Date[]
let types: string[]
let completeTypes: { [key: string]: string[] }
/* eslint-enable @typescript-eslint/no-unused-vars */

// Types partagés
type BatchValue = {
  reporder: { [periode: string]: RepOrder }
  compact: { delete: { [dataType: string]: DataHash[] } }
}

type DataHash = string

type Periode = Date

type Siret = string

type RepOrder = {
  random_order: number
  periode: Periode
  siret: Siret
}

type Keys = Set<DataHash>
type State = { [key: string]: Keys }
