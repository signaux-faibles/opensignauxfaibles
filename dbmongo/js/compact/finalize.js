function finalize(k, o) {
  "use strict";
  o.index = {"algo1":false,
    "algo2":false}


  if (o.scope == "entreprise") {
    o.index.algo1 = true
    o.index.algo2 = true
  } else {
    // Est-ce que l'un des batchs a un effectif ?
    var batches = Object.keys(o.batch)
    batches.some(batch => {
      let hasEffectif = Object.keys(o.batch[batch].effectif || {}).length > 0
      o.index.algo1 = hasEffectif
      o.index.algo2 = hasEffectif
      return (hasEffectif)
    })
    // Complete reporder if missing
    o = f.complete_reporder(k, o)
  }
  return(o)
}
