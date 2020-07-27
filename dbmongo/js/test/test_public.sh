# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

result_public=$(jsc \
  helpers/test_env_for_public.js \
  helpers/reducers.js \
  data/objects.js \
  ../common/!(*_test*).js \
  ../public/!(*_test*).js \
  ../public/_test.js \
  2>&1)
if [ "$result_public" != 'true' ]; then
  echo "$result_public"
  exit 1
fi
