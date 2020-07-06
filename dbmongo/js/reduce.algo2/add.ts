export function add(
  obj: { [periode: string]: unknown },
  output: { [periode: string]: unknown }
): void {
  "use strict"
  Object.keys(output).forEach(function (periode) {
    if (periode in obj) {
      Object.assign(output[periode], obj[periode])
    }
  })
}
