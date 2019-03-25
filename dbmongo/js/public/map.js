function map() {
  var value = f.flatten(this.value, actual_batch)
  
  if (this.value.scope=="etablissement") {
    vcrp = {}
    vcmde = {}

    vcrp.effectif = f.effectifs(value)
    vcmde.dernier_effectif = vcrp.effectif[vcrp.effectif.length - 1]
    // if (v.cotisation) {
    //   let output_cotisationsdettes = f.cotisationsdettes(v)
    //   v.cotisation = output_cotisation
    // }

    if (vcrp.effectif.length > 0) {
      emit({scope: ["bfc", "crp"], key: this.value.key, batch: actual_batch}, vcrp)
      emit({scope: ["bfc"], key: this.value.key, batch: actual_batch}, vcmde)
    }


  }
  // else if (this.value.scope=='entreprise') {
  //   emit({scope: "entreprise", key: this.value.key, batch: actual_batch}, v) 
  // }
}