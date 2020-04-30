# This file is run by dbmongo/js_test.go.

result_add=$(jsc \
  ../compact/currentState.js \
  ../compact/currentState_test.js\
  2>&1)
if [ "$result_add" != 'true' ]; then
  echo "$result_add"
  exit 1
fi

result_add=$(jsc \
  ./helpers/testing.js \
  ../compact/currentState.js \
  ../compact/reduce.js \
  ../compact/reduce_test.js\
  2>&1)
if [ "$result_add" != 'true' ]; then
  echo "$result_add"
  exit 1
fi

result_finalize=$(jsc \
  ./helpers/testing.js \
  ../compact/currentState.js \
  ../compact/complete_reporder.js \
  ../compact/finalize.js \
  ../compact/finalize_test.js\
  2>&1)
if [ "$result_finalize" != 'true' ]; then
  echo "$result_finalize"
  exit 1
fi
