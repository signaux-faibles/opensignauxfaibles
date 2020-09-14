#!/bin/bash

set -e # will stop the script if any command fails with a non-zero exit code

# We use perl because sed adds an empty line at the end of every js file,
# which was adding changes to git's staging, while debugging failing tests.
perl -pi'' -e 's/^const .*$//g' ./**/*.js
perl -pi'' -e 's/^export //' ./**/*.js
perl -pi'' -e 's/^import .*$//g' ./**/*.js

function getGlobals {
  grep -F --no-filename 'declare const' $@ \
    | cut -d' ' -f3 \
    | cut -d':' -f1 \
    | sort -u \
    | uniq \
    | paste -sd "," -
}

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

checkJS "compact/*.js" "f,emit,$(getGlobals 'compact/*.ts')"
checkJS "public/*.js" "f,emit,$(getGlobals 'public/*.ts')"
checkJS "reduce.algo2/*.js" "f,print,emit,bsonsize,$(getGlobals 'reduce.algo2/*.ts')"
