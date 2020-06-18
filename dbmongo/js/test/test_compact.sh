# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

result_currentState=$(jsc \
  ../compact/currentState.js \
  ../compact/currentState_test.js\
  2>&1)
if [ "$result_currentState" != 'true' ]; then
  echo "result_currentState: $result_currentState"
  exit 1
fi

result_reduce=$(jsc \
  ./helpers/fakes.js \
  ./helpers/testing.js \
  ../common/!(*_test*).js \
  ../compact/currentState.js \
  ../compact/reduce.js \
  ../compact/reduce_test.js\
  2>&1)
if [ "$result_reduce" != 'true' ]; then
  echo "result_reduce: $result_reduce"
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
