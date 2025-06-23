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

heading "go generate"
(go generate ./...) 2>&1 | indent

heading "make build"
(killall sfdata 2>/dev/null || true; make build && echo "📦 sfdata") 2>&1 | indent

heading "go test"
if [[ "$*" == *--update* ]]
then
    (go test ./... -test.count=1 -update) 2>&1 | indent
else
    (go test ./... -test.count=1) 2>&1 | indent
fi

heading "test-cli.sh"
./tests/test-cli.sh $@ 2>&1 | indent


heading "test-validate.sh"
./tests/test-validate.sh $@ 2>&1 | indent

heading "test-check.sh"
./tests/test-check.sh $@ 2>&1 | indent

heading "test-import.sh"
./tests/test-import.sh $@ 2>&1 | indent

heading "test-parseFile.sh"
./tests/test-parseFile.sh $@ 2>&1 | indent

# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "ℹ️  Golden master files were updated"
fi
