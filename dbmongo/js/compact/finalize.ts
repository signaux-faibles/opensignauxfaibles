import "../globals.ts"

// finalize permet de:
// - indiquer les établissements à inclure dans les calculs de variables
// (processus reduce.algo2)
// - intégrer les reporder pour permettre la reproductibilité de
// l'échantillonnage pour l'entraînement du modèle.
export function finalize(k: Siret, o: RawDataValues): RawDataValues {
  "use strict"
  o.index = { algo1: false, algo2: false }

  if (o.scope === "entreprise") {
    o.index.algo1 = true
    o.index.algo2 = true
  } else {
    // Est-ce que l'un des batchs a un effectif ?
    const batches = Object.keys(o.batch)
    batches.some((batch) => {
      const hasEffectif = Object.keys(o.batch[batch].effectif || {}).length > 0
      o.index.algo1 = hasEffectif
      o.index.algo2 = hasEffectif
      return hasEffectif
    })
    // Complete reporder if missing
    // TODO: do not complete if all indexes are false.
    o = f.complete_reporder(k, o)
  }
  return o
}
