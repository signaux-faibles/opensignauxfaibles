function idEntreprise(idEtablissement) {
  return {
    scope: 'entreprise',
    key: idEtablissement.slice(0,9),
    batch: actual_batch
  }
}