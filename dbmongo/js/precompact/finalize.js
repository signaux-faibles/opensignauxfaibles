function finalize(k, o) {
    
  o.index = {"algo1":false,
    "algo2":false}


  if (o.scope == "entreprise") {
    o.index.algo1 = true
    o.index.algo2 = true
  } else {
    // Est-ce que l'un des batchs a un effectif ? 
    Object.keys(o.batch).some(batch => {
      let hasEffectif = Object.keys(o.batch[batch].effectif || {}).length > 0 
      o.index.algo1 = hasEffectif 
      o.index.algo2 = hasEffectif
      return (hasEffectif)
    })
  }
  return o
}
