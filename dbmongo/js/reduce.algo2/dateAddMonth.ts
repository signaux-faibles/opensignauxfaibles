export function dateAddMonth(date: Date, nbMonth: number): Date {
  "use strict"
  const result = new Date(date.getTime())
  result.setUTCMonth(result.getUTCMonth() + nbMonth)
  return result
}
