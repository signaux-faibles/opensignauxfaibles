function idEntreprise(idEtablissement) {
  "use strict";
  return {
    scope: 'entreprise',
    key: idEtablissement.slice(0,9),
    batch: actual_batch
  }
}