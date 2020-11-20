import { f } from "./functions"
import {
  CompanyDataValues,
  CompanyDataValuesWithFlags,
  SiretOrSiren,
} from "../RawDataTypes"

// finalize permet de:
// - indiquer les établissements à inclure dans les calculs de variables
// (processus reduce.algo2)
// - intégrer les reporder pour permettre la reproductibilité de
// l'échantillonnage pour l'entraînement du modèle.
export function finalize(
  k: SiretOrSiren,
  companyDataValues: CompanyDataValues
): CompanyDataValuesWithFlags {
  "use strict"

  let o: CompanyDataValuesWithFlags = {
    ...companyDataValues,
    index: { algo1: false, algo2: false },
  }

  if (o.scope === "entreprise") {
    o.index.algo1 = true
    o.index.algo2 = true
  } else {
    // Est-ce que l'un des batchs a un effectif ?
    const batches = Object.keys(o.batch)
    if (batches.some((batch) => hasEffectif(o, batch))) {
      o.index.algo1 = true
      o.index.algo2 = true
      // Complete reporder if missing
      o = f.complete_reporder(k, o)
    }
    // do not complete if all indexes are false.
  }
  return o
}

const hasEffectif = (o: CompanyDataValuesWithFlags, batch: string) =>
  Object.keys(o.batch[batch].effectif || {}).length > 0
