#!/bin/bash

set -e # will stop the script if any command fails with a non-zero exit code

# We use perl because sed adds an empty line at the end of every js file,
# which was adding changes to git's staging, while debugging failing tests.
perl -pi'' -e 's/^const .*$//g' ./**/*.js
perl -pi'' -e 's/^export //' ./**/*.js
perl -pi'' -e 's/^import .*$//g' ./**/*.js

# Check that JS files only call functions through the f global variable.
GLOBALS="f,emit,fromBatchKey,batches,serie_periode,completeTypes" && $(npm bin)/eslint --no-eslintrc --parser-options=ecmaVersion:6 --env es6 --global "${GLOBALS}" --rule "no-undef:2" --quiet --ignore-pattern functions.js compact/*.js
GLOBALS="f,emit,fromBatchKey,batches,serie_periode,completeTypes,date_fin,actual_batch" && $(npm bin)/eslint --no-eslintrc --parser-options=ecmaVersion:6 --env es6 --global "${GLOBALS}" --rule "no-undef:2" --quiet --ignore-pattern functions.js public/*.js
GLOBALS="f,print,emit,bsonsize,fromBatchKey,batches,serie_periode,completeTypes,date_fin,actual_batch,offset_effectif,includes,naf" && $(npm bin)/eslint --no-eslintrc --parser-options=ecmaVersion:6 --env es6 --global "${GLOBALS}" --rule "no-undef:2" --quiet --ignore-pattern functions.js reduce.algo2/*.js

# TODO: extract a function, to reduce code duplication
# TODO: extract globals from code
