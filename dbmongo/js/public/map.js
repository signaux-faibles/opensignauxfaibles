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
    vcmde.compte = f.compte(value.compte)
    vcmde.procol = f.dealWithProcols(value.altares, "altares",  null).concat(f.dealWithProcols(value.procol, "procol",  null))
    if (vcmde.procol.length >= 1){
      vcmde.last_procol = vcmde.procol[vcmde.procol.length - 1]
    } else {
      vcmde.last_procol = {"etat": "in_bonis"}
    }
    emit(actual_batch + "_etablissement_" + this.value.key, vcmde)
    emit(actual_batch + "_entreprise_" + this.value.key.slice(0,9), {sirets: [this.value.key]})
  }
  else if (this.value.scope == "entreprise") {
    let v = {}
    let diane = f.diane(value.diane)
    let bdf = f.bdf(value.bdf)
    let sirene_ul = (value.sirene_ul || {})[Object.keys(value.sirene_ul || {})[0] || ""]
    let crp = value.crp
    if (diane.length > 0) {
      v.diane = diane
    }
    if (bdf.length > 0) {
      v.bdf = bdf
    }
    if (sirene_ul) {
      v.sirene_ul = sirene_ul
    }
    if (crp) {
      v.crp = crp
    }
    if (Object.keys(v) != []) {
      emit(actual_batch + "_entreprise_" + this.value.key, v)
    }
  }
}
