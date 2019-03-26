function map() {
  var v = f.flatten(this.value, actual_batch)

  
  if (this.value.scope=="etablissement") {
    emit({scope: "etablissement", key: this.value.key, batch: actual_batch}, v) 
  }
  else if (this.value.scope=='entreprise') {
    emit({scope: "entreprise", key: this.value.key, batch: actual_batch}, v) 
  }
}