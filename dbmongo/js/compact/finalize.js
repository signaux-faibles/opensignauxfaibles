function finalize(k, o) {
  // Pour tous les batchs qui ont été intégrés
  //    batches.reduce((m, batch) => {
  //       // Set des types que où l'on jette le passé 
  //        var deleteOld = new Set(completeTypes[batch])
  //        
  //        // on crée l'objet batch s'il n'existe pas
  //        // avec la mention compact["status"]= "false" si la mention n'existait pas
  //        o.batch[batch] = (o.batch[batch] || {})
  //        o.batch[batch].compact = (o.batch[batch].compact || {})
  //        o.batch[batch].compact["status"] = (o.batch[batch].compact["status"]||false)
  //
  //        // Pour tous les types possibles 
  //        types.forEach(type => {
  //            // on crée l'objet type
  //            // on rajoute la mention compact.delete qui est un objet vide
  //            o.batch[batch][type] = (o.batch[batch][type]||{})
  //            m[type] = (m[type] || new Set()) // mémoire m des objets déjà intégré dans le passé
  //            var keys = Object.keys(o.batch[batch][type])
  //            o.batch[batch].compact.delete = (o.batch[batch].compact.delete||{})
  //
  //             // si c'est un type où l'on jette le passé, on se débarasse des clés
  //            // qui ne sont pas du dernier batch (filter)
  //           // en les ajoutant  à compact.delete[type]
  //            if (deleteOld.has(type) && o.batch[batch].compact.status == false) {
  //                var discardKeys = [...m[type]].filter(key => !(new Set(keys).has(key)))
  //                o.batch[batch].compact.delete[type] = discardKeys;
  //            }
  //
  //            // on supprime finalement ces clés de la mémoire m
  //            if (deleteOld.has(type)) {
  //                o.batch[batch].compact.delete[type] = (o.batch[batch].compact.delete[type] || {})
  //                o.batch[batch].compact.delete[type].forEach(key => {
  //                    m[type].delete(key)
  //                })
  //            }
  //
  //            // on filtre les nouvelles clés qu'on a déjà en mémoire m (rien n'a changé)
  //            keys.filter(key => (m[type].has(key))).forEach(key => delete o.batch[batch][type][key])
  //            // ces nouvelles clés viennent compléter la mémoire m
  //            m[type] = new Set([...m[type]].concat(keys))
  //            // on supprime les types vides.
  //            if (Object.keys(o.batch[batch][type]).length == 0) {delete o.batch[batch][type]}
  //        })
  //
  //        //batch compacté : check
  //        o.batch[batch].compact = (o.batch[batch].compact||{})
  //        o.batch[batch].compact.status = true
  //        return m
  //    }, {})

  // 1er Filtrage
  // Utilisé dans Reduce Handler actuellement
  // Entreprise qui n'ont aucune information effectif 

  o.index = {"algo1":false,
    "algo2":false}


  if (o.scope == "entreprise") {
    o.index.algo1 = true
    o.index.algo2 = true
  } else  {
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
