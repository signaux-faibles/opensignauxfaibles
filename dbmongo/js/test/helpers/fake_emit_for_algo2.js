"use strict";

// from /dbmongo/js/test/algo2/lib_algo2.js

const pool = {}

function emit(key, value) {
  const id = key.siren + key.batch + key.periode.getTime()
  pool[id] = (pool[id] || []).concat([{key, value}])
}
