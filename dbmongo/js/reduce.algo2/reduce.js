function reduce(key, values) {
    "use strict"
    return values.reduce((val, accu) => {
        return Object.assign(accu, val)
    }, {})
}

exports.reduce = reduce
