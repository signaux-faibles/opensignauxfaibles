f = {
  lookAhead
}

function compare(a, b) {
  if (Object.keys(a).length != Object.keys(b).length){
    return false
  }
  var equal = Object.keys(a).every(function(k) {
    return(JSON.stringify(a[k]) == JSON.stringify(b[k]))
  })
  return equal
}
