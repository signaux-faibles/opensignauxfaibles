#!/bin/bash

# Usage:
# $ ./test-all.sh           # pour éxecuter tous les tests
# $ ./test-all.sh --update  # pour éxecuter tous les tests et mettre à jour les snapshots des tests go + les golden files des tests de bout en bout

function heading {
  echo ""
  echo "–––––"
  echo "$1"
  echo "–––––"
}

function indent {
  sed 's/^/  /'
}

# Mandatory tests (can stop the script)

set -e # will stop the script if any command fails with a non-zero exit code
set -o pipefail # ... even for tests which pipe their output to indent

heading "make build"
(killall sfdata 2>/dev/null || true; make build && echo "📦 sfdata") 2>&1 | indent

if [[ "$*" == *--update* ]]
then
    heading "Update tests"
    (go test ./... -test.count=1) 2>&1 | indent

    heading "Update golden files"
    (go test -test.count=1 \
      ./lib/filter \
      ./lib/parsing/sirene \
      ./lib/parsing/sirene_ul \
      ./lib/parsing/sirene_histo \
      ./lib/parsing/urssaf \
      -update) 2>&1 | indent

    heading "Update e2e tests"
    (go test -tags=e2e -test.count=1 . -update) 2>&1 | indent
else
    heading "go test"
    (go test ./... -test.count=1) 2>&1 | indent

    heading "go test e2e"
    (go test ./... -tags=e2e -test.count=1) 2>&1 | indent
fi
