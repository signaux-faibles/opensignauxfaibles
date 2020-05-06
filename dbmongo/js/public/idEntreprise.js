function idEntreprise(idEtablissement) {
  return {
    scope: 'entreprise',
    key: idEtablissement.slice(0,9),
    batch: jsParams.actual_batch
  }
}