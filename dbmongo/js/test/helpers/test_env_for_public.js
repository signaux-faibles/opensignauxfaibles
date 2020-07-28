"use strict"

// globals

const f = this
const exports = { f }
const ISODate = (date) => new Date(date)

// from /dbmongo/js/test/public/lib_public.js

const pool = {}

function emit(key, value) {
  const id = JSON.stringify(key)
  pool[id] = (pool[id] || []).concat([{key, value}])
}
