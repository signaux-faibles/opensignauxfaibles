"use strict";

function compare(a, b) {
  if (Object.keys(a).length != Object.keys(b).length){
    return false
  }
  var equal = Object.keys(a).every(function(k) {
    return(JSON.stringify(a[k]) == JSON.stringify(b[k]))
  })
  return equal
}

function compareIgnoreRandom(a, b) {
  if (Object.keys(a).length != Object.keys(b).length){
    return false
  }
  var equal = Object.keys(a).every(function(k) {
    return(k == "random_order" || // Ignore random numbers
      JSON.stringify(a[k]) == JSON.stringify(b[k]) || //Compare exact match
      compareIgnoreRandom(a[k], b[k])) //Recursive call to compare if any random_order
  })
  return equal
}
