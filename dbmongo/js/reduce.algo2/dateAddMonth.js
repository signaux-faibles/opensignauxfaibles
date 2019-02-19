function dateAddMonth(date, nbMonth) {
  var result = new Date(date.getTime())
  result.setUTCMonth(result.getUTCMonth() + nbMonth)
  return result
}