# This file is run by dbmongo/js_test.go.

set -e

result_public=$(jsc ../public/*.js ../common/*.js public/lib_public.js objects.js public/test_public.js)
if [ "$result_public" != 'true' ]; then
  echo "$result_public"
  exit 1
fi
