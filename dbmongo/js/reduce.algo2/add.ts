export function add<Added, Target>(
  obj: { [periode: string]: Added },
  output: { [periode: string]: Target }
): void {
  "use strict"
  Object.keys(output).forEach(function (periode) {
    if (periode in obj) {
      Object.assign(output[periode], obj[periode])
    }
  })
}
