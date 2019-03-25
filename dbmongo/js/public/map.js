function map() {
  var value = f.flatten(this.value, actual_batch)
  
  if (this.value.scope=="etablissement") {
    vcrp = {}
    vcmde = {}

    vcmde.effectif = f.effectifs(value)
    vcmde.dernier_effectif = vcmde.effectif[vcmde.effectif.length - 1]
    vcmde.sirene = f.sirene(f.iterable(value.sirene))
    // if (v.cotisation) {
    //   let output_cotisationsdettes = f.cotisationsdettes(v)
    //   v.cotisation = output_cotisation
    // }

    if (vcmde.effectif.length > 0) {
      // emit({scope: ["bfc", "crp"], key: this.value.key, batch: actual_batch}, vcrp)
      emit({scope: "etablissement", key: this.value.key, batch: actual_batch}, vcmde)
    }


  }
  else if (this.value.scope=='entreprise') {
     v = {
       diane: f.diane(value.diane),
       bdf: f.bdf(value.bdf)
     }
     emit({scope: "entreprise", key: this.value.key, batch: actual_batch}, v) 
  }
}