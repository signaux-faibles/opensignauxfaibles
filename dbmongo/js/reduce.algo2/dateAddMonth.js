function dateAddMonth(date, nbMonth) {
  "use strict";
  var result = new Date(date.getTime())
  result.setUTCMonth(result.getUTCMonth() + nbMonth)
  return result
}