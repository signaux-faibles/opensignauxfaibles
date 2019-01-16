function naf(output_indexed, naf) {
  Object.keys(output_indexed).forEach(k =>{
    if (("code_ape" in output_indexed[k]) && (output_indexed[k].code_ape !== null)){
      output_indexed[k].code_naf = naf.n5to1[output_indexed[k].code_ape]
      output_indexed[k].libelle_naf = naf.n1[output_indexed[k].code_naf]
      output_indexed[k].libelle_ape2 = naf.n2[output_indexed[k].code_ape.substring(0,2)]
      output_indexed[k].libelle_ape3 = naf.n3[output_indexed[k].code_ape.substring(0,3)]
      output_indexed[k].libelle_ape4 = naf.n4[output_indexed[k].code_ape.substring(0,4)]
      output_indexed[k].libelle_ape5 = naf.n5[output_indexed[k].code_ape]
    }
  })
}