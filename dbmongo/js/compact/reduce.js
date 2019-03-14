function reduce(key, values) {

  //if (key == "30493863200011") {var deleteme = true} else {var deleteme = false}

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

  // Pour tous les batchs qui ont été intégrés
  batches.reduce((m, batch) => {
    //if (deleteme) {
      //print("-------- BATCH --------------------------------") // deleteme
      //print(batch) } // deleteme
    // Set des types que où l'on jette le passé pour le batch courant
    if (!reduced_value.batch[batch]) { 
      return m 
      //if (deleteme) { print("Batch skipped: not interesting")}
    }
    var deleteOld = new Set(completeTypes[batch])

    // on s'en fiche du compact.status ? 

    // Faut-il se soucier de la disparition de stocks passés?
    var stock_types = completeTypes[batch].filter(type => (m[type] || new Set()).size > 0)
    //if (deleteme) print("Les types complets ", completeTypes[batch])
    //if (deleteme) print("Les types complets avec un mémoire ", stock_types)
    // Les données qui ont bougé dans le batch en cours
    var new_types =  Object.keys(reduced_value.batch[batch])
    // On dédoublonne au besoin
    var all_interesting_types = [...new Set([...stock_types, ...new_types])]

    //if (deleteme) print("Types intéressants à explorer: ")
    //if (deleteme) print(all_interesting_types)

    // Pour tous les types intéressants
    all_interesting_types.forEach( type => {
      //if (deleteme) { print("ooo ", type, " ooo") }

      if (type == "compact") return 
      // on crée l'objet type en mémoire s'il n'existe pas
      m[type] = m[type] || new Set()
      // clés pour ce batch et ce type
      var keys = Object.keys(reduced_value.batch[batch][type] || {})

      //if (deleteme) print(" Clés envisagées : ", keys) 
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
          //if (deleteme) print("Type stock non encore compacté")
          reduced_value.batch[batch].compact.delete = reduced_value.batch[batch].compact.delete || {}
          var discardKeys = [...m[type]].filter(key => !(new Set(keys).has(key)))
          //if (deleteme) print(" Marquons les clés qu'on ne retrouve pas dans le dernier batch")
          //if (deleteme) print(discardKeys)
          reduced_value.batch[batch].compact.delete[type] = discardKeys;
        }
        // le cas échéant, on supprime ces clés de la mémoire m
        reduced_value.batch[batch].compact.delete[type] = (reduced_value.batch[batch].compact.delete[type] || [] )
        //if (deleteme) print(" Ces clés sont à supprimer de la mémoire:")
        //if (deleteme) print(reduced_value.batch[batch].compact.delete[type])
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
      //if (deleteme) print("Mes clés déjà connues")
      //if (deleteme) print(keys.filter(key => (m[type].has(key))))
      keys.filter(key => (m[type].has(key))).forEach(key => delete reduced_value.batch[batch][type][key])
      m[type] = new Set([...m[type]].concat(keys))
      //if (deleteme) print("''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''")
      //if (deleteme) print("Etat de la mémoire:",type)
      //if (deleteme) print([...m[type]])
      // on supprime les types vides.
      if (reduced_value.batch[batch][type] && Object.keys(reduced_value.batch[batch][type]).length == 0) {
        delete reduced_value.batch[batch][type]
      }
    })
    return m
  }, {})

  return(reduced_value)
}
