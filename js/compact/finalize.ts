import { f } from "./functions"
import {
  CompanyDataValues,
  CompanyDataValuesWithFlags,
  Siret,
  Siren,
} from "../RawDataTypes"

// finalize permet de:
// - indiquer les établissements à inclure dans les calculs de variables
// (processus reduce.algo2)
// - intégrer les reporder pour permettre la reproductibilité de
// l'échantillonnage pour l'entraînement du modèle.
export function finalize(
  k: Siret | Siren,
  companyDataValues: CompanyDataValues
): CompanyDataValuesWithFlags {
  "use strict"

  let o: CompanyDataValuesWithFlags = {
    ...companyDataValues,
    index: { algo2: false },
  }

  if (o.scope === "entreprise") {
    o.index.algo2 = true
  } else {
    // Est-ce que l'un des batchs a un effectif ?
    const batches = Object.keys(o.batch)
    batches.some((batch) => {
      const hasEffectif = Object.keys(o.batch[batch]?.effectif || {}).length > 0
      o.index.algo2 = hasEffectif
      return hasEffectif
    })
    // Complete reporder if missing
    o = f.complete_reporder(k, o)
  }
  return o
}
