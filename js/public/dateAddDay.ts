export function dateAddDay(date: Date, nbDays: number): Date {
  "use strict"
  const result = new Date(date.getTime())
  result.setDate(result.getDate() + nbDays)
  return result
}
