function populateNafAndApe(output_indexed, naf) {
  "use strict";
  Object.keys(output_indexed).forEach(k =>{
    if (("code_ape" in output_indexed[k]) && (output_indexed[k].code_ape !== null)){
      var code_ape = output_indexed[k].code_ape
      output_indexed[k].code_naf = naf.n5to1[code_ape]
      output_indexed[k].libelle_naf = naf.n1[output_indexed[k].code_naf]
      output_indexed[k].code_ape_niveau2 = code_ape.substring(0,2)
      output_indexed[k].code_ape_niveau3 = code_ape.substring(0,3)
      output_indexed[k].code_ape_niveau4 = code_ape.substring(0,4)
      output_indexed[k].libelle_ape2 = naf.n2[output_indexed[k].code_ape_niveau2]
      output_indexed[k].libelle_ape3 = naf.n3[output_indexed[k].code_ape_niveau3]
      output_indexed[k].libelle_ape4 = naf.n4[output_indexed[k].code_ape_niveau4]
      output_indexed[k].libelle_ape5 = naf.n5[code_ape]
    }
  })
}

exports.populateNafAndApe = populateNafAndApe
