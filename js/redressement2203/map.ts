/* eslint-disable @typescript-eslint/no-unused-vars */
import { f } from "./functions"
import {
  CompanyDataValues,
  SortieRedressementUrssaf2203,
} from "../RawDataTypes"
import { EntréeDebit } from "../GeneratedTypes"

export type SortieMap = SortieRedressementUrssaf2203

declare const dateStr: string

// Types de données en entrée et sortie
export type Input = { _id: unknown; value: CompanyDataValues }
export type OutKey = string
export type OutValue = Partial<SortieMap>

class SommesDettes {
  public partOuvriere: number
  public partPatronale: number

  constructor() {
    this.partOuvriere = 0
    this.partPatronale = 0
  }
}
declare function emit(key: string, value: OutValue): void

function recupererValeursUniquesEcartsNegatifs(debits: EntréeDebit[]) {
  const ecartsNegatifs = debits.map((debit) => debit.numero_ecart_negatif)
  return [...new Set(ecartsNegatifs)]
}

function recupererDetteParType(debits: EntréeDebit[]): SommesDettes {
  const ecartsNegatifs = recupererValeursUniquesEcartsNegatifs(debits)
  let mostRecentBatch: EntréeDebit
  const sommesDettes: SommesDettes = new SommesDettes()
  for (const _en of ecartsNegatifs) {
    mostRecentBatch = debits.reduce((a, b) =>
      a.periode.start > b.periode.start ? a : b
    )
    sommesDettes.partOuvriere += mostRecentBatch.part_ouvriere
    sommesDettes.partPatronale += mostRecentBatch.part_patronale
  }
  return sommesDettes
}

export function map(this: Input): void {
  const testDate = new Date(dateStr)

  const values = f.flatten(this.value, "2203")
  const beforeBatches = []
  const afterBatches = []

  if (values.debit) {
    for (const debit of Object.values(values.debit)) {
      debit.periode.start > testDate
        ? afterBatches.push(debit)
        : beforeBatches.push(debit)
    }
  }

  const dettesAnciennesParECN: SommesDettes = recupererDetteParType(
    beforeBatches
  )
  const dettesAnciennesDebutParECN: SommesDettes = recupererDetteParType(
    beforeBatches.filter((b) => b.date_traitement <= testDate)
  )

  const dettesRecentesParECN: SommesDettes = recupererDetteParType(afterBatches)

  emit(this.value.key, {
    partPatronaleAncienne: dettesAnciennesParECN.partPatronale,
    partOuvriereAncienne: dettesAnciennesParECN.partOuvriere,
    partPatronaleRecente: dettesRecentesParECN.partPatronale,
    partOuvriereRecente: dettesRecentesParECN.partOuvriere,
    partOuvriereAncienneDebut: dettesAnciennesDebutParECN.partOuvriere,
    partPatronaleAncienneDebut: dettesAnciennesDebutParECN.partPatronale,
  })
}
