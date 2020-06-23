function outputs(v, serie_periode) {
    "use strict"
    var output_array = serie_periode.map(function (e) {
        return {
            siret: v.key,
            periode: e,
            effectif: null,
            etat_proc_collective: "in_bonis",
            interessante_urssaf: true,
            outcome: false,
        }
    })

    var output_indexed = output_array.reduce(function (periodes, val) {
        periodes[val.periode.getTime()] = val
        return periodes
    }, {})

    return [output_array, output_indexed]
}

exports.outputs = outputs
