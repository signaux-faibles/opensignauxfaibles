function outputs (v, serie_periode) {
  var output_array = serie_periode.map(function (e) {
    return {
      "siret": v.siret,
      "periode": e,
      "effectif": null,
      "date_effectif": null,
      "apart_heures_consommees": 0,
      "apart_motif_recours": 0,
      "debit_array": [],
      "etat_proc_collective": "in_bonis"
    }
  });
    
  var output_indexed = output_array.reduce(function (periode, val) {
      periode[val.periode.getTime()] = val
      return periode
  }, {})

  return [output_array, output_indexed]
}
