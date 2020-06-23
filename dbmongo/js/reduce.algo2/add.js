function add(obj, output) {
    "use strict"
    Object.keys(output).forEach(function (periode) {
        if (periode in obj) {
            Object.assign(output[periode], obj[periode])
        } else {
            // throw new EvalError(
            //   "Attention, l'objet à fusionner ne possède pas les mêmes périodes que l'objet dans lequel il est fusionné"
            // )
        }
    })
}

exports.add = add
