// import { f } from "./functions"
// import {
//   EntréeApConso,
//   EntréeApDemande,
//   EntréeCompte,
//   EntréeDelai,
//   EntréeDéfaillances,
//   EntréeDiane,
//   EntréeEllisphere,
//   EntréePaydex,
//   EntréeSirene,
//   EntréeSireneEntreprise,
// } from "../GeneratedTypes"
import {
  CompanyDataValues,
  BatchKey,
  SortieRedressementUrssaf2203,
} from "../RawDataTypes"
// import { SortieDebit } from "./debits"
// import { Bdf } from "./bdf"

export type SortieMap = SortieRedressementUrssaf2203

// Paramètres globaux utilisés par "public"
declare const actual_batch: BatchKey

// Types de données en entrée et sortie
export type Input = { _id: unknown; value: CompanyDataValues }
export type OutKey = string
export type OutValue = Partial<SortieMap>
declare function emit(key: string, value: OutValue): void

export function map(this: Input): void {
  emit(this.value.key, {
    partPatronaleAncienne: 0,
    partOuvriereAncienne: 0,
    partPatronaleRecente: 0,
    partOuvriereRecente: 0,
  })
}
