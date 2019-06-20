function dateAddDay(date, nbMonth) {
  var result = new Date(date.getTime())
  result.setDate( result.getDate() + nbMonth );
  return result
}