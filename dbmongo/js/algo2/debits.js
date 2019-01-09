function debits (v) {
  // relier les débits
  var ecn = Object.keys(v.debit).reduce((m, h) => {
      var d = [h, v.debit[h]]
      var start = d[1].periode.start
      var end = d[1].periode.end
      var num_ecn = d[1].numero_ecart_negatif
      var compte = d[1].numero_compte
      var key = start + "-" + end + "-" + num_ecn + "-" + compte
      m[key] = (m[key] || []).concat([{
          "hash": d[0],
          "numero_historique": d[1].numero_historique,
          "date_traitement": d[1].date_traitement
      }]) 
      return m
  }, {})

  Object.keys(ecn).forEach(i => {
      ecn[i].sort(f.compareDebit)
      var l = ecn[i].length
      ecn[i].forEach((e, idx) => {
          if (idx <= l - 2) {
              v.debit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
          }
      })
  })

}