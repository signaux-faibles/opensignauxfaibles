function compte(v, periodes) {
  "use strict"
  const output_compte = {}

  //  var offset_compte = 3
  Object.keys(v.compte).forEach((hash) => {
    const periode = v.compte[hash].periode.getTime()

    output_compte[periode] = output_compte[periode] || {}
    output_compte[periode].compte_urssaf = v.compte[hash].numero_compte
  })

  return output_compte
}
