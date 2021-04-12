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

# Extract a comma-separated list of global variables that are expected by TypeScript files passed as arguments.
function getGlobals {
  grep -F --no-filename 'declare const' $@ \
    | cut -d' ' -f3 \
    | cut -d':' -f1 \
    | sort -u \
    | uniq \
    | paste -sd "," -
}

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
checkJS "compact/*.js" "f,emit,$(getGlobals 'compact/*.ts')"
checkJS "public/*.js" "f,emit,$(getGlobals 'public/*.ts')"
checkJS "reduce.algo2/*.js" "f,print,emit,bsonsize,$(getGlobals 'reduce.algo2/*.ts')"

# TODO: effacer ces lignes -- utilisées pour référence pendant le développement de PR #345
# getGlobals 'compact/*.ts' # => batches,completeTypes,fromBatchKey,serie_periode
# getGlobals 'public/*.ts' # => actual_batch,date_fin,serie_periode
# getGlobals 'reduce.algo2/*.ts' # => actual_batch,date_fin,includes,naf,offset_effectif,serie_periode
# getGlobals 'common/*.ts' => 
