# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

result_public=$(jsc \
  --strict-file=helpers/fakes.js \
  --strict-file=helpers/fake_emit_for_public.js \
  --strict-file=data/objects.js \
  --strict-file=../common/!(*_test).js \
  --strict-file=../public/!(*_test).js \
  --strict-file=../public/_test.js \
  2>&1)
if [ "$result_public" != 'true' ]; then
  echo "$result_public"
  exit 1
fi
