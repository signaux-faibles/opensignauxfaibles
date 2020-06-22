export function add<Added, Target>(
  obj: { [periode: string]: Added },
  output: { [periode: string]: Target },
  array?: Target[]
): [(Added & Target)[], { [periode: string]: Added & Target }] {
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
  return [
    array as (Added & Target)[],
    output as { [periode: string]: Added & Target },
  ]
}
