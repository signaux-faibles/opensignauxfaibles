#!/bin/bash

set -e # will stop the script if any command fails with a non-zero exit code

# Generate GeneratedTypes.d.ts from validation/*/schema.json files.
./generate-types.sh

# Run typescript transpiler, to generate .js files from .ts files.
$(npm bin)/tsc --p "tsconfig-transpilation.json"

# Exclude JavaScript keywords that are not supported by MongoDB.
perl -pi'' -e 's/^const .*$//g;' -e 's/^export //;' -e 's/^import .*$//g' ./**/*.js
# Note: We use perl because sed adds an empty line at the end of every js file,
# which was adding changes to git's staging, while debugging failing tests.

# Fails if any JavaScript file references a symbol that is not included in the list of expected globals.
function checkJS {
  FILES="$1"
  GLOBALS="$2"
  $(npm bin)/eslint --no-eslintrc \
    --parser-options=ecmaVersion:6 --env es6 \
    --rule "no-undef:2" --quiet \
    --ignore-pattern "functions.js" \
    --global "${GLOBALS}" \
    "${FILES}"
}

# Check that all functions and constants referenced in JavaScript files will be
# made available by MongoDB or the map-reduce command sent by sfdata.
checkJS "compact/*.js"      "f,emit,$(./get-globals.sh 'compact/*.ts')"
checkJS "public/*.js"       "f,emit,$(./get-globals.sh 'public/*.ts')"
checkJS "purgeBatch/*.js"   "f,emit,$(./get-globals.sh 'purgeBatch/*.ts')"
checkJS "reduce.algo2/*.js" "f,emit,print,bsonsize,$(./get-globals.sh 'reduce.algo2/*.ts')"

# TODO: effacer ces lignes -- utilisées pour référence pendant le développement de PR #345
# ./get-globals.sh 'compact/*.ts'      # => batches,completeTypes,fromBatchKey,serie_periode ✅
# ./get-globals.sh 'public/*.ts'       # => actual_batch,date_fin,serie_periode ✅
# ./get-globals.sh 'reduce.algo2/*.ts' # => actual_batch,date_fin,includes,naf,offset_effectif,serie_periode
# ./get-globals.sh 'common/*.ts'       # => 
# ./get-globals.sh 'purgeBatch/*.ts'   # => fromBatchKey
