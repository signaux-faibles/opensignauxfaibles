import "../globals.ts"
import { compactBatch } from "./compactBatch"
import * as f from "./currentState"

// Entrée: données d'entreprises venant de ImportedData, regroupées par entreprise ou établissement.
// Sortie: un objet fusionné par entreprise ou établissement, contenant les données historiques et les données importées, à destination de la collection RawData.
// Opérations: retrait des données doublons et application des corrections de données éventuelles.
export function reduce(
  key: SiretOrSiren,
  values: CompanyDataValues[] // chaque element contient plusieurs batches pour cette entreprise ou établissement
): CompanyDataValues {
  "use strict"

  // Tester si plusieurs batchs. Reduce complet uniquement si plusieurs
  // batchs. Sinon, juste fusion des attributs
  const auxBatchSet = new Set()
  const severalBatches = values.some((value) => {
    auxBatchSet.add(Object.keys(value.batch || {}))
    return auxBatchSet.size > 1
  })

  // Fusion batch par batch des types de données sans se préoccuper des doublons.
  const naivelyMergedCompanyData: CompanyDataValues = values.reduce(
    (m, value: CompanyDataValues) => {
      Object.keys(value.batch).forEach((batch) => {
        type DataType = keyof BatchValue
        m.batch[batch] = (Object.keys(value.batch[batch]) as DataType[]).reduce(
          (batchValues: BatchValue, type: DataType) => ({
            ...batchValues,
            [type]: value.batch[batch][type],
          }),
          m.batch[batch] || {}
        )
      })
      return m
    },
    { key: key, scope: values[0].scope, batch: {} }
  )

  // Cette fonction reduce() est appelée à deux moments:
  // 1. agregation par établissement d'objets ImportedData. Dans cet étape, on
  // ne travaille généralement que sur un seul batch.
  // 2. agregation de ces résultats au sein de RawData, en fusionnant avec les
  // données potentiellement présentes. Dans cette étape, on fusionne
  // généralement les données de plusieurs batches. (données historiques)

  if (!severalBatches) return naivelyMergedCompanyData

  //////////////////////////////////////////////////
  // ETAPES DE LA FUSION AVEC DONNÉES HISTORIQUES //
  //////////////////////////////////////////////////

  // 0. On calcule la memoire au moment du batch à modifier
  const memory_batches: BatchValue[] = Object.keys(
    naivelyMergedCompanyData.batch
  )
    .filter((batch) => batch < batchKey)
    .sort()
    .reduce((m: BatchValue[], batch: string) => {
      m.push(naivelyMergedCompanyData.batch[batch])
      return m
    }, [])

  // Memory conserve les données aplaties de tous les batches jusqu'à batchKey
  // puis sera enrichie au fur et à mesure du traitements des batches suivants.
  const memory = f.currentState(memory_batches)

  const reduced_value = naivelyMergedCompanyData

  // On itère sur chaque batch à partir de batchKey pour les compacter.
  // Il est possible qu'il y ait moins de batch en sortie que le nombre traité
  // dans la boucle, si ces batchs n'apportent aucune information nouvelle.
  batches
    .filter((batch) => batch >= batchKey)
    .forEach((batch) => {
      reduced_value.batch[batch] = reduced_value.batch[batch] || {}
      const currentBatch = reduced_value.batch[batch]
      const compactedBatch = compactBatch(currentBatch, memory, batch)
      if (Object.keys(compactedBatch).length === 0) {
        delete reduced_value.batch[batch]
      } else {
        reduced_value.batch[batch] = compactedBatch
      }
    })

  return reduced_value
}
