# This file is run by dbmongo/js_test.go.

set -e

result_raison_sociale=$(jsc ../common/*.js helpers/testing.js ../common/raison_sociale_test.js)
if [ "$result_raison_sociale" != 'true' ]; then
  echo "$result_raison_sociale"
  exit 1
fi

