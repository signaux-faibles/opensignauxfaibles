# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

result_raison_sociale=$(jsc \
  helpers/testing.js \
  helpers/fakes.js \
  ../common/!(*_test).js \
  ../common/raison_sociale_test.js \
  2>&1)
if [ "$result_raison_sociale" != 'true' ]; then
  echo "$result_raison_sociale"
  exit 1
fi
