/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { EntréeDebit } from "../GeneratedTypes"
import { SommesDettes } from "../RawDataTypes"
import { f } from "./functions"

export function recupererDetteTotale(debits: EntréeDebit[]): SommesDettes {
  const ecartsNegatifs = f.recupererValeursUniquesEcartsNegatifs(debits)
  // let mostRecentBatch: EntréeDebit
  const sommesDettes: SommesDettes = {
    partOuvriere: 0,
    partPatronale: 0,
  }
  for (const en of ecartsNegatifs) {
    const debitsECN = debits.filter((d) => d.numero_ecart_negatif === en)
    const mostRecentBatch = debitsECN.reduce((a, b) =>
      a.date_traitement > b.date_traitement ? a : b
    )
    sommesDettes.partOuvriere += mostRecentBatch.part_ouvriere
    sommesDettes.partPatronale += mostRecentBatch.part_patronale
    // if (debitsECN.length > 0) {
    //   const mostRecentBatch = debitsECN.sort(
    //     (a, b) => b.date_traitement.getTime() - a.date_traitement.getTime()
    //   )

    // if (mostRecentBatch.length > 0) {
    //   const latestBatch = mostRecentBatch[0]!
    //   sommesDettes.partOuvriere += latestBatch.part_ouvriere
    //   sommesDettes.partPatronale += latestBatch.part_patronale
    // }
    // }
  }
  return sommesDettes
}
