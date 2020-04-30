# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

# TODO: use `2>&1` instead of `set -e`, in all sh tests
result_public=$(jsc ../public/!(*_test).js ../common/!(*_test).js helpers/fakes.js objects.js ../public/_test.js 2>&1)
if [ "$result_public" != 'true' ]; then
  echo "$result_public"
  exit 1
fi
