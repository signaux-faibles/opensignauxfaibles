function map () {

  let v = f.flatten(this.value, actual_batch)

  if (v.scope == "etablissement") {
    let o = f.outputs(v, serie_periode)
    let output_array = o[0]
    let output_indexed = o[1]

    if (v.debit) { f.debits(v) }
    if (v.effectif) {f.effectifs(v, output_array, output_indexed)}
    if (v.apconso && v.apdemande) {f.apart(v, output_indexed)}
    if (v.delai) {f.delais(v, output_indexed)}
    if (v.altares) {f.defaillances(v, output_indexed)}
    if (v.cotisations && v.debits) {f.cotisationsdettes(v, output_array, output_indexed)}
    if (v.ccsf) {f.ccsf(v, output_array)}
    if (v.sirene) {f.sirene(v, output_array)}

    output_array.forEach(val => {
      emit(    
        {'siret': this._id,
          'scope': this.value.scope,
          'batch': actual_batch,
          'algo': 'algo2',
          'periode': val.periode},
          val
      )
    })
  }
  
  
}