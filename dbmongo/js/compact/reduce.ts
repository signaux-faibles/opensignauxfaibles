import {
  CompanyDataValues,
  BatchValue,
  BatchKey,
  SiretOrSiren,
} from "../RawDataTypes"
import { compactBatch } from "./compactBatch"
import { currentState } from "./currentState"

// Paramètres globaux utilisés par "compact"
declare const batches: BatchKey[]
declare const fromBatchKey: BatchKey

// Entrée: données d'entreprises venant de ImportedData, regroupées par entreprise ou établissement.
// Sortie: un objet fusionné par entreprise ou établissement, contenant les données historiques et les données importées, à destination de la collection RawData.
// Opérations: retrait des données doublons et application des corrections de données éventuelles.
export function reduce(
  key: SiretOrSiren,
  values: CompanyDataValues[] // chaque element contient plusieurs batches pour cette entreprise ou établissement
): CompanyDataValues {
  "use strict"

  /* DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO */ const f = {
    compactBatch, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    currentState, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

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
    { key, scope: values[0].scope, batch: {} }
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
  const memoryBatches: BatchValue[] = Object.keys(
    naivelyMergedCompanyData.batch
  )
    .filter((batch) => batch < fromBatchKey)
    .sort()
    .reduce((m: BatchValue[], batch: string) => {
      m.push(naivelyMergedCompanyData.batch[batch])
      return m
    }, [])

  // Memory conserve les données aplaties de tous les batches jusqu'à fromBatchKey
  // puis sera enrichie au fur et à mesure du traitement des batches suivants.
  const memory = f.currentState(memoryBatches)

  const reducedValue: CompanyDataValues = {
    key: naivelyMergedCompanyData.key,
    scope: naivelyMergedCompanyData.scope,
    batch: {},
  }

  // Copie telle quelle des batches jusqu'à fromBatchKey.
  Object.keys(naivelyMergedCompanyData.batch)
    .filter((batch) => batch < fromBatchKey)
    .forEach((batch) => {
      reducedValue.batch[batch] = naivelyMergedCompanyData.batch[batch]
    })

  // On itère sur chaque batch à partir de fromBatchKey pour les compacter.
  // Il est possible qu'il y ait moins de batch en sortie que le nombre traité
  // dans la boucle, si ces batchs n'apportent aucune information nouvelle.
  batches
    .filter((batch) => batch >= fromBatchKey)
    .forEach((batch) => {
      const currentBatch = naivelyMergedCompanyData.batch[batch]
      const compactedBatch = f.compactBatch(currentBatch, memory, batch)
      if (Object.keys(compactedBatch).length > 0) {
        reducedValue.batch[batch] = compactedBatch
      }
    })

  return reducedValue
}
