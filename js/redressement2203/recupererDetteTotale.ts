/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { EntréeDebit } from "../GeneratedTypes"
import { SommesDettes } from "../RawDataTypes"
import { f } from "./functions"

export function recupererDetteTotale(debits: EntréeDebit[]): SommesDettes {
  const ecartsNegatifs = f.recupererValeursUniquesEcartsNegatifs(debits)
  const sommesDettes: SommesDettes = {
    partOuvriere: 0,
    partPatronale: 0,
  }
  for (const en of ecartsNegatifs) {
    const debitsECN = debits.filter((d) => f.cleEcartNegatif(d) === en)
    const mostRecentBatch = debitsECN.reduce((a, b) =>
      a.numero_historique > b.numero_historique ? a : b
    )
    sommesDettes.partOuvriere += mostRecentBatch.part_ouvriere
    sommesDettes.partPatronale += mostRecentBatch.part_patronale
  }
  return sommesDettes
}
