#!/bin/bash

# Usage:
# $ git secret reveal                 # pour déchiffrer les données utilisées par les tests (golden files, etc...)
# $ ./test-all.sh                     # pour éxecuter tous les tests
# $ ./test-all.sh --update-snapshots  # pour éxecuter tous les tests et mettre à jour les snapshots des tests go + "ava" + les golden files des tests de bout en bout
# $ git secret changes                # pour visualiser les modifications éventuellement apportées aux golden files
# $ git secret hide                   # pour chiffrer les golden files suite à leur modification

function heading {
  echo ""
  echo "–––––"
  echo "$1"
  echo "–––––"
}

function indent {
  sed 's/^/  /'
}

# Optional tests (cannot stop the script)

heading "pick specified node.js version"
(. ~/.nvm/nvm.sh; cd ./js && nvm use) 2>&1 | indent

# Mandatory tests (can stop the script)

set -e # will stop the script if any command fails with a non-zero exit code
set -o pipefail # ... even for tests which pipe their output to indent

heading "npm install"
(cd ./js && npm install) 2>&1 | indent

heading "npm test"
(cd ./js && npm run lint && npm test -- $@) 2>&1 | indent

heading "go test"
if [[ "$*" == *--update* ]]
then
    (go test ./... -test.count=1 -update) 2>&1 | indent
else
    (go test ./... -test.count=1) 2>&1 | indent
fi

heading "go generate"
(cd ./lib/engine && go generate .) 2>&1 | indent

heading "go build"
(killall sfdata 2>/dev/null || true; go build -o "sfdata" && echo "📦 sfdata") 2>&1 | indent

heading "test-cli.sh"
./tests/test-cli.sh $@ 2>&1 | indent

heading "test-prune-entities.sh"
./tests/test-prune-entities.sh $@ 2>&1 | indent

heading "test.sh"
./tests/test.sh $@ 2>&1 | indent

heading "test-validate.sh"
./tests/test-validate.sh $@ 2>&1 | indent

heading "test-check.sh"
./tests/test-check.sh $@ 2>&1 | indent

heading "test-import.sh"
./tests/test-import.sh $@ 2>&1 | indent

heading "test-compact.sh"
./tests/test-compact.sh $@ 2>&1 | indent

heading "test-compact-failure.sh"
./tests/test-compact-failure.sh $@ 2>&1 | indent

heading "test-public.sh"
./tests/test-public.sh $@ 2>&1 | indent

heading "test-reduce.sh"
./tests/test-reduce.sh $@ 2>&1 | indent

heading "test-reduce-2.sh"
./tests/test-reduce-2.sh $@ 2>&1 | indent

heading "test-purge-batch.sh"
./tests/test-purge-batch.sh $@ 2>&1 | indent

heading "test-export.sh"
./tests/test-export.sh $@ 2>&1 | indent

# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "ℹ️  Golden master files were updated => you may have to run: $ git secret hide" # to re-encrypt the golden master files
fi
