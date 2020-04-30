# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

# TODO: use `2>&1` instead of `set -e`, in all sh tests
result_public=$(jsc \
  helpers/fakes.js \
  ../common/!(*_test).js \
  objects.js \
  ../public/!(*_test).js \
  ../public/_test.js \
  2>&1)
if [ "$result_public" != 'true' ]; then
  echo "$result_public"
  exit 1
fi
