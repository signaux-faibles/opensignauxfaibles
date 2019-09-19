function map() {
  var value = f.flatten(this.value, actual_batch)

  if (this.value.scope=="etablissement") {
    // var vcrp = {}
    let vcmde = {}

    vcmde.effectif = f.effectifs(value)
    vcmde.dernier_effectif = vcmde.effectif[vcmde.effectif.length - 1]
    vcmde.sirene = f.sirene(f.iterable(value.sirene))
    vcmde.cotisation = f.cotisations(value.cotisation)
    vcmde.debit = f.debits(value.debit)
    vcmde.apconso = f.apconso(value.apconso)
    vcmde.apdemande = f.apconso(value.apdemande)
    vcmde.idEntreprise = f.idEntreprise(this._id)
    vcmde.delai = f.delai(value.delai)
    vcmde.procol = f.dealWithProcols(value.altares, "altares",  null).concat(f.dealWithProcols(value.procol, "procol",  null)) 

    if (vcmde.procol.length >= 1){
      vcmde.last_procol = vcmde.procol[vcmde.procol.length - 1]
    } else {
      vcmde.last_procol = {"etat": "in_bonis"}
    }
    //if (vcmde.effectif.length > 0) {
    // emit({scope: ["bfc", "crp"], key: this.value.key, batch: actual_batch}, vcrp)
    emit({scope: "etablissement", key: this.value.key, batch: actual_batch}, vcmde)
    emit({scope: "entreprise", key: this.value.key.slice(0,9), batch: actual_batch}, {sirets: [this.value.key]})
  }
  else if (this.value.scope == "entreprise") {
    let v = {
      diane: f.diane(value.diane),
      bdf: f.bdf(value.bdf),
      sirene_ul: (value.sirene_ul || {})[Object.keys(value.sirene_ul || {})[0] || ""],
    }
    emit({scope: "entreprise", key: this.value.key, batch: actual_batch}, v)
  }
}
