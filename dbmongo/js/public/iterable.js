function iterable(dict) {
  "use strict";
  try {
    return Object.keys(dict).map(h => {
      return dict[h]
    })
  } catch(error) {
    return []
  }
}

exports.iterable = iterable
