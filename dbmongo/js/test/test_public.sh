# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

result_public=$(jsc \
  helpers/fakes.js \
  helpers/reducers.js \
  helpers/fake_emit_for_public.js \
  data/objects.js \
  ../common/!(*_test*).js \
  ../public/!(*_test*).js \
  ../public/_test.js \
  2>&1)
if [ "$result_public" != 'true' ]; then
  echo "$result_public"
  exit 1
fi
