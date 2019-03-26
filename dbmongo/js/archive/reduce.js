function reduce(key, values) {
  //fusion des objets dans values
  let reduced_value = values.reduce((m, value) => {
    Object.keys((value.batch||{})).forEach(batch => {
      m.batch = (m.batch||{})
      m.batch[batch] = (m.batch[batch] || {})
      Object.keys(value.batch[batch]).forEach(type => {
        m.batch[batch][type] = (m.batch[batch][type] || {})
        Object.assign(m.batch[batch][type],value.batch[batch][type])
      })
    })
    return m
  }, {"key": key, "scope": values[0].scope  })

  ///////////////////////////////////
  ///// ETAPES //////////////////////
  ///////////////////////////////////
  // 0. On indique le premier batch modifié, tous les suivants le seront aussi.
  // 0bis. On calcule la mémoire au moment du batch à modifier
  // Pour tous les batchs à modifier:
  // 1. Pour le batch en cours, on regarde les clés ajoutées, les clés supprimées
  // 2. On ajoute aux clés supprimées les types stocks de la mémoire. 
  // 3.a Pour chaque clé supprimée: est-ce qu'elle est bien dans la mémoire ? sinon on la retire (pas de maj mémoire)
  // 3.b Est-ce qu'elle a été également ajoutée ? Dans ce cas là, on retire les deux 
  // i.e. on hérite de la mémoire. (pas de maj de la mémoire)
  // 3.c On retire les clés restantes de la mémoire. 
  // 4.a Pour chaque clé ajoutée: est-ce qu'elle est dans la mémoire ? Si oui on filtre cette clé
  // i.e. on hérite de la mémoire. (pas de maj de la mémoire)
  // 4.b Pour chaque clé ajoutée restante: on ajoute à la mémoire. 
  // 
  // Pour tous les batchs qui ont été intégrés (dans l'ordre alphabétique)
  //TODO gérer les suppressions si le batch ne contient aucune clé !
  batches.reduce((m, batch) => {
    //if (!reduced_value.batch[batch]) { 
    //  return m 
    //} // NE FONCTIONNE PAS: les types complets peuvent imposer des suppressions
    reduced_value.batch[batch] = reduced_value.batch[batch] || {}

    
    // Set des types où l'on jette le passé pour le batch courant
    // Y a-t-il potentiellement une modification de stocks passés?
    var deleteOld = new Set(completeTypes[batch])
    var stock_types = completeTypes[batch].filter(type => (m[type] || new Set()).size > 0)
    // Les données qui ont bougé dans le batch en cours
    var new_types =  Object.keys(reduced_value.batch[batch])
    // On dédoublonne au besoin
    var all_interesting_types = [...new Set([...stock_types, ...new_types])]


    // Pour tous les types intéressants
    all_interesting_types.forEach( type => {

      if (type == "compact") return 
      // on crée l'objet type en mémoire s'il n'existe pas
      m[type] = m[type] || new Set()
      // clés pour ce batch et ce type
      var keys = Object.keys(reduced_value.batch[batch][type] || {})

      /////////////////////////////////////////////////////////////////////////
      // ETAPE: supprimer les anciennes valeurs pour les types concernés //////
      /////////////////////////////////////////////////////////////////////////
      // si c'est un type où l'on jette le passé, on se débarasse des clés
      // qui ne sont pas du dernier batch (filter)
      // en les ajoutant  à compact.delete[type]
      if (deleteOld.has(type) ) {
        // on veut garder en mémoire si on a déjà calculé les clés à supprimer
        // TODO c'est cette étape qui empêche d'importer un batch dans le passé !
        reduced_value.batch[batch].compact = reduced_value.batch[batch].compact || {}
        reduced_value.batch[batch].compact.status = reduced_value.batch[batch].compact.status || false
        // Si déjà compacté, on passe notre chemin, sinon on marque les clés à supprimer
        if (reduced_value.batch[batch].compact.status == false) {
          reduced_value.batch[batch].compact.delete = reduced_value.batch[batch].compact.delete || {}
          var discardKeys = [...m[type]].filter(key => !(new Set(keys).has(key)))
          reduced_value.batch[batch].compact.delete[type] = discardKeys;
        }
        // le cas échéant, on supprime ces clés de la mémoire m
        reduced_value.batch[batch].compact.delete[type] = (reduced_value.batch[batch].compact.delete[type] || [] )
        reduced_value.batch[batch].compact.delete[type].forEach(key => {
          m[type].delete(key)
        })
      }
      /////////////////////////////////////////////////////////////////////////
      // ETAPE: ajouter les clés au batch qui sont neuves par rapport aux /////
      // clés connues                                                     /////
      /////////////////////////////////////////////////////////////////////////
      // on filtre les nouvelles clés qu'on connait déjà d'un batch précédent 
      // Les autres viennent compléter la mémoire m
      keys.filter(key => (m[type].has(key))).forEach(key => delete reduced_value.batch[batch][type][key])
      m[type] = new Set([...m[type]].concat(keys))
      // on supprime les types vides.
      if (reduced_value.batch[batch][type] && Object.keys(reduced_value.batch[batch][type]).length == 0) {
        delete reduced_value.batch[batch][type]
      }
    })
    if (reduced_value.batch[batch] && Object.keys(reduced_value.batch[batch]).length == 0 ) {
      delete reduced_value.batch[batch]
    }
    return m
  }, {})

  return(reduced_value)
}
