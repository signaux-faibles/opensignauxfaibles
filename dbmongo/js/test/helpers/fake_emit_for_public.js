// from /dbmongo/js/test/public/lib_public.js

var pool = {}

function emit(key, value) {
  id = JSON.stringify(key)
  pool[id] = (pool[id] || []).concat([{key, value}])
}
