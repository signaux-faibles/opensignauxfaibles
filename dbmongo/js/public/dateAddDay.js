function dateAddDay(date, nbMonth) {
  "use strict";
  var result = new Date(date.getTime())
  result.setDate( result.getDate() + nbMonth );
  return result
}