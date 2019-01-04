function reduce(key, values) {
    return values.reduce((m, value) => {
        Object.keys((value.batch||{})).forEach(batch => {
            Object.keys(value.batch[batch]).forEach(type => {
                m.batch = (m.batch||{})
                m.batch[batch] = (m.batch[batch] || {})
                m.batch[batch][type] = (m.batch[batch][type] || {})
                Object.assign(m.batch[batch][type],value.batch[batch][type])
            })
        })

        return m
    },{"key": key, "type": this.value.scope  })

}