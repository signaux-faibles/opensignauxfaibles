# This file is run by dbmongo/js_test.go.

set -e # stop on errors
shopt -s extglob # enable exclusion of test files in wildcard

result_raison_sociale=$(jsc ../common/!(*_test).js helpers/testing.js ../common/raison_sociale_test.js)
if [ "$result_raison_sociale" != 'true' ]; then
  echo "$result_raison_sociale"
  exit 1
fi

