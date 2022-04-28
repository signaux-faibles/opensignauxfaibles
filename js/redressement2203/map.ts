import { f } from "./functions"
import {
  CompanyDataValues,
  SommesDettes,
  SortieRedressementUrssaf2203,
} from "../RawDataTypes"

export type SortieMap = SortieRedressementUrssaf2203

declare const dateStr: string

// Types de données en entrée et sortie
export type Input = { _id: unknown; value: CompanyDataValues }
export type OutKey = string
export type OutValue = Partial<SortieMap>

declare function emit(key: string, value: OutValue): void

export function map(this: Input): void {
  const testDate = new Date(dateStr)

  const values = f.flatten(this.value, "2203")
  const beforeBatches = [] // TODO : renommer les variables
  const afterBatches = []

  if (values.debit) {
    for (const debit of Object.values(values.debit)) {
      debit.periode.start > testDate
        ? afterBatches.push(debit)
        : beforeBatches.push(debit)
    }
  }
  const dettesAnciennesParECN: SommesDettes = f.recupererDetteTotale(
    beforeBatches
  )

  const dettesAnciennesDebutParECN: SommesDettes = f.recupererDetteTotale(
    beforeBatches.filter((b) => b.date_traitement <= testDate)
  )

  const dettesRecentesParECN: SommesDettes = f.recupererDetteTotale(
    afterBatches
  )

  emit(this.value.key, {
    partPatronaleAncienne: dettesAnciennesParECN.partPatronale,
    partOuvriereAncienne: dettesAnciennesParECN.partOuvriere,
    partPatronaleRecente: dettesRecentesParECN.partPatronale,
    partOuvriereRecente: dettesRecentesParECN.partOuvriere,
    partOuvriereAncienneDebut: dettesAnciennesDebutParECN.partOuvriere,
    partPatronaleAncienneDebut: dettesAnciennesDebutParECN.partPatronale,
  })
}
