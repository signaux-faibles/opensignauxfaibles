function reduce(key, values) {

  // Tester si plusieurs batchs. Reduce complet uniquement si plusieurs
  // batchs. Sinon, juste fusion des attributs
  let auxBatchSet = new Set()

  let severalBatches = values.some(value => {
    auxBatchSet.add(Object.keys(value.batch || {}))
    return auxBatchSet.size > 1
  })

  //fusion des attributs dans values
  let reduced_value = values.reduce((m, value) => {
    Object.keys((value.batch||{})).forEach(batch => {
      m.batch = (m.batch||{})
      m.batch[batch] = (m.batch[batch] || {})
      Object.keys(value.batch[batch]).forEach(type => {
        m.batch[batch][type] = (m.batch[batch][type] || {})
        Object.assign(m.batch[batch][type], value.batch[batch][type])
      })
    })
    return m
  }, {"key": key, "scope": values[0].scope  })

  if (!severalBatches) return(reduced_value)

  ///////////////////////////////////
  ///// ETAPES //////////////////////
  ///////////////////////////////////
  // Uniquement si severalBatches
  // 0. On calcule la memoire au moment du batch à modifier
  var memory_batches = Object.keys(reduced_value.batch).filter( batch =>
    batch < batchKey
  ).sort().reduce((m, batch) => {
    m.push(reduced_value.batch[batch])
    return(m)
  },[])

  var memory = f.currentState(memory_batches)

  // Pour tous les batchs à modifier, c'est-à-dire le batch ajouté et tous les
  // suivants.
  var modified_batches = batches.filter( batch =>
    batch >= batchKey
  )

  modified_batches.forEach(batch => {

    reduced_value.batch[batch] = reduced_value.batch[batch] || {}

    // Les types où il y  a potentiellement des suppressions
    var stock_types = completeTypes[batch].filter(type => (memory[type] || new Set()).size > 0)
    // Les types qui ont bougé dans le batch en cours
    var new_types =  Object.keys(reduced_value.batch[batch])
    // On dedoublonne au besoin
    var all_interesting_types = [...new Set([...stock_types, ...new_types])]

    // Filtrage selon les types effectivement importés
    if (types.length > 0){
      stock_types = stock_types.filter(type => types.includes(type))
      new_types = new_types.filter(type => types.includes(type))
      all_interesting_types = all_interesting_types.filter(type => types.includes(type))
    }

    // 1. On recupère les cles ajoutes et les cles supprimes
    // -----------------------------------------------------

    var hashToDelete = {}
    var hashToAdd = {}

    all_interesting_types.forEach(type => {
      // Le type compact gère les clés supprimées
      if (type == "compact") {
        if (reduced_value.batch[batch].compact.delete){
          Object.keys(reduced_value.batch[batch].compact.delete).forEach(delete_type => {
            reduced_value.batch[batch].compact.delete[delete_type].forEach(hash => {
              hashToDelete[delete_type] = hashToDelete[delete_type] || new Set()
              hashToDelete[delete_type].add(hash)
            })
          })
        }
      } else {
        Object.keys(reduced_value.batch[batch][type] || {}).forEach(hash => {
          hashToAdd[type] = hashToAdd[type] || new Set()
          hashToAdd[type].add(hash)
        })
      }
    })

    //
    // 2. On ajoute aux cles supprimees les types stocks de la memoire.
    // ----------------------------------------------------------------

    stock_types.forEach(type => {
      hashToDelete[type] = new Set([...(hashToDelete[type] || new Set()) ,
        ...memory[type]])
    })


    Object.keys(hashToDelete).forEach(type => {

      // 3.a Pour chaque cle supprimee: est-ce qu'elle est bien dans la
      // memoire ? sinon on la retire de la liste des clés supprimées (pas de
      // maj memoire)
      // -----------------------------------------------------------------------------------------------------------------
      hashToDelete[type] = new Set([...hashToDelete[type]].filter( hash => {
        return((memory[type] || new Set()).has(hash))
      }))


      // 3.b Est-ce qu'elle a ete egalement ajoutee en même temps que
      // supprimée ? (par exemple remplacement d'un stock complet à
      // l'identique) Dans ce cas là, on retire cette clé des valeurs ajoutées
      // et supprimées
      // i.e. on herite de la memoire. (pas de maj de la memoire)
      // ------------------------------------------------------------------------------

      hashToDelete[type] = new Set([...hashToDelete[type]].filter( hash => {
        let also_added = (hashToAdd[type] || new Set()).has(hash)
        if (also_added) {
          hashToAdd[type].delete(hash)
        }
        return(!also_added)
      }))

      // 3.c On retire les cles restantes de la memoire.
      // --------------------------------------------------
      hashToDelete[type].forEach( hash => {
        memory[type].delete(hash)
      })

    })

    Object.keys(hashToAdd).forEach(type => {

      // 4.a Pour chaque cle ajoutee: est-ce qu'elle est dans la memoire ? Si oui on filtre cette cle
      // i.e. on herite de la memoire. (pas de maj de la memoire)
      // ---------------------------------------------------------------------------------------------
      hashToAdd[type] = new Set([...hashToAdd[type]].filter( hash => {
        return(!(memory[type] || new Set()).has(hash))
      }))


      // 4.b Pour chaque cle ajoutee restante: on ajoute à la memoire.
      // -------------------------------------------------------------

      hashToAdd[type].forEach( hash => {
        memory[type] = memory[type] || new Set()
        memory[type].add(hash)
      })
    })


    // 5. On met à jour reduced_value
    // -------------------------------
    stock_types.forEach(type => {
      if (hashToDelete[type]) {
        reduced_value.batch[batch].compact = reduced_value.batch[batch].compact  || {}
        reduced_value.batch[batch].compact.delete = reduced_value.batch[batch].compact.delete  || {}
        reduced_value.batch[batch].compact.delete[type] = [...hashToDelete[type]]
      }
    })


    new_types.forEach(type => {
      if (hashToAdd[type]) {
        reduced_value.batch[batch][type] = Object.keys(reduced_value.batch[batch][type] || {}).filter( hash => {
          return(hashToAdd[type].has(hash))
        }).reduce( (m, hash) => {
          m[hash] = reduced_value.batch[batch][type][hash]
          return(m)
        }, {})
      }
    })

    // 6. nettoyage
    // ------------

    if (reduced_value.batch[batch]){
      //types vides
      Object.keys(reduced_value.batch[batch]).forEach( type => {
        if (Object.keys(reduced_value.batch[batch][type]).length == 0){
          delete reduced_value.batch[batch][type]
        }
      })
      //hash à supprimer vides (compact.delete)
      if (reduced_value.batch[batch].compact && reduced_value.batch[batch].compact.delete) {
        Object.keys(reduced_value.batch[batch].compact.delete).forEach( type => {
          if (reduced_value.batch[batch].compact.delete[type].length == 0){
            delete reduced_value.batch[batch].compact.delete[type]
          }
        })
        if (Object.keys(reduced_value.batch[batch].compact.delete).length == 0 ) {
          delete reduced_value.batch[batch].compact
        }
      }
      //batchs vides
      if (Object.keys(reduced_value.batch[batch]).length == 0 ) {
        delete reduced_value.batch[batch]
      }
    }
  })

  return(reduced_value)
}
