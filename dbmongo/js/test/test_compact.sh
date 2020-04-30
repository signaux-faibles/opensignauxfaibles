# This file is run by dbmongo/js_test.go.

result_add=$(jsc \
  ../compact/currentState.js \
  ../compact/currentState_test.js)
if [ "$result_add" != 'true' ]; then
  echo "$result_add"
  exit 1
fi

result_add=$(jsc \
  ../compact/currentState.js \
  ../compact/reduce.js \
  ./helpers/testing.js \
  ../compact/reduce_test.js)
if [ "$result_add" != 'true' ]; then
  echo "$result_add"
  exit 1
fi

result_finalize=$(jsc \
  ../compact/currentState.js \
  ../compact/finalize.js \
  ../compact/complete_reporder.js \
  ./helpers/testing.js \
  ../compact/finalize_test.js)
if [ "$result_finalize" != 'true' ]; then
  echo "$result_finalize"
  exit 1
fi
