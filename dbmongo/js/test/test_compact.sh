# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

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
