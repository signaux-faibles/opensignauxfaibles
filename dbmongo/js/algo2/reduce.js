function reduce(_, values) {
    return values.reduce((val, accu) => {
        return Object.assign(accu, val)
    }, {})
}