"use strict";

// from /dbmongo/js/test/public/lib_public.js

const pool = {}

function emit(key, value) {
  const id = JSON.stringify(key)
  pool[id] = (pool[id] || []).concat([{key, value}])
}
