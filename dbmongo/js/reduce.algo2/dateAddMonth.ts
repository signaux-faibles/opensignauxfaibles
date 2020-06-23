function dateAddMonth(date, nbMonth) {
  "use strict"
  const result = new Date(date.getTime())
  result.setUTCMonth(result.getUTCMonth() + nbMonth)
  return result
}

exports.dateAddMonth = dateAddMonth
