function delais (v, output_indexed) {
  Object.keys(v.delai).map(
    function (hash) {
      var delai = v.delai[hash]
      var date_creation = new Date(Date.UTC(delai.date_creation.getUTCFullYear(), delai.date_creation.getUTCMonth(), 1, 0, 0, 0, 0))
      var date_echeance = new Date(Date.UTC(delai.date_echeance.getUTCFullYear(), delai.date_echeance.getUTCMonth(), 1, 0, 0, 0, 0))
      var pastYearTimes = generatePeriodSerie(date_creation, date_echeance).map(function (date) { return date.getTime() })
      pastYearTimes.map(
        function(time){
          if (time in output_indexed) {
            var remaining_months = (date_echeance.getUTCMonth() - new Date(time).getUTCMonth()) +
              12*(date_echeance.getUTCFullYear() - new Date(time).getUTCFullYear()) 
            output_indexed[time].delai = remaining_months
            output_indexed[time].duree_delai = delai.duree_delai
            output_indexed[time].montant_echeancier = delai.montant_echeancier

            if (delai.duree_delai > 0){
              output_indexed[time].ratio_dette_delai = (output_indexed[time].montant_part_patronale + 
                output_indexed[time].montant_part_ouvriere - delai.montant_echeancier * remaining_months * 30 / 
                (delai.duree_delai)) /
                delai.montant_echeancier
            }
          }
        }
      )
    }
  ) 
}