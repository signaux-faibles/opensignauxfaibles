function compte (v, periodes) {
  let output_compte = {}


  //  var offset_compte = 3
  Object.keys(v.compte).forEach(hash =>{
    var periode = compte[hash].periode.getTime()

    output_compte[periode] =  output_compte[periode] || {}
    output_compte[periode].compte_urssaf =  compte[hash].numero_compte
  })

  return output_compte
}
