#!/bin/bash

set -e # will stop the script if any command fails with a non-zero exit code

# We use perl because sed adds an empty line at the end of every js file,
# which was adding changes to git's staging, while debugging failing tests.
perl -pi'' -e 's/^const .*$//g' ./**/*.js
perl -pi'' -e 's/^export //' ./**/*.js
perl -pi'' -e 's/^import .*$//g' ./**/*.js

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

# Check that JS files only call functions through the f global variable.
checkJS "compact/*.js" "f,emit,fromBatchKey,batches,serie_periode,completeTypes"
checkJS "public/*.js" "f,emit,fromBatchKey,batches,serie_periode,completeTypes,date_fin,actual_batch"
checkJS "reduce.algo2/*.js" "f,print,emit,bsonsize,fromBatchKey,batches,serie_periode,completeTypes,date_fin,actual_batch,offset_effectif,includes,naf"

# TODO: extract globals from code
