function detteFiscale (diane){
  if  (("dette_fiscale_et_sociale" in diane) && (diane["dette_fiscale_et_sociale"] !== null) &&
      ("valeur_ajoutee" in diane) && (diane["valeur_ajoutee"] !== null) &&
      (diane["valeur_ajoutee"] != 0)){
    return diane["dette_fiscale_et_sociale"]/ diane["valeur_ajoutee"] * 100
  } else {
    return null
  }
}