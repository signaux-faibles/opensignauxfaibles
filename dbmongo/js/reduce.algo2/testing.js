Object.assign = function(target) {
  if (target == null) {
    throw new TypeError("Cannot convert undefined or null to object")
  }

  target = Object(target)
  for (var index = 1; index < arguments.length; index++) {
    var source = arguments[index]
    if (source != null) {
      for (var key in source) {
        if (Object.prototype.hasOwnProperty.call(source, key)) {
          target[key] = source[key]
        }
      }
    }
  }
  return target
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
