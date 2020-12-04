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
(. ~/.nvm/nvm.sh; cd ./dbmongo/js && nvm use) 2>&1 | indent

# Mandatory tests (can stop the script)

set -e # will stop the script if any command fails with a non-zero exit code
set -o pipefail # ... even for tests which pipe their output to indent

heading "npm install"
(cd ./dbmongo/js && npm install) 2>&1 | indent

heading "npm test"
(cd ./dbmongo/js && npm run lint && npm test -- $@) 2>&1 | indent

heading "go test"
if [[ "$*" == *--update* ]]
then
    (cd ./dbmongo && go test ./... -test.count=1 -update) 2>&1 | indent
else
    (cd ./dbmongo && go test ./... -test.count=1) 2>&1 | indent
fi

heading "go generate"
(cd ./dbmongo/lib/engine && go generate .) 2>&1 | indent

heading "go build"
(killall dbmongo 2>/dev/null || true; cd ./dbmongo && go build && echo "📦 dbmongo/dbmongo") 2>&1 | indent

heading "test-api-prune-entities.sh"
./tests/test-api-prune-entities.sh $@ 2>&1 | indent

heading "test-api.sh"
./tests/test-api.sh $@ 2>&1 | indent

heading "test-api-validate.sh"
./tests/test-api-validate.sh $@ 2>&1 | indent

heading "test-api-check.sh"
./tests/test-api-check.sh $@ 2>&1 | indent

heading "test-api-import.sh"
./tests/test-api-import.sh $@ 2>&1 | indent

heading "test-api-compact.sh"
./tests/test-api-compact.sh $@ 2>&1 | indent

heading "test-api-compact-failure.sh"
./tests/test-api-compact-failure.sh $@ 2>&1 | indent

heading "test-api-public.sh"
./tests/test-api-public.sh $@ 2>&1 | indent

heading "test-api-reduce.sh"
./tests/test-api-reduce.sh $@ 2>&1 | indent

heading "test-api-reduce-2.sh"
./tests/test-api-reduce-2.sh $@ 2>&1 | indent

heading "test-api-purge-batch.sh"
./tests/test-api-purge-batch.sh $@ 2>&1 | indent

heading "test-api-export.sh"
./tests/test-api-export.sh $@ 2>&1 | indent

heading "test-api-swagger.sh"
./tests/test-api-swagger.sh $@ 2>&1 | indent

# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "ℹ️  Golden master files were updated => you may have to run: $ git secret hide" # to re-encrypt the golden master files
fi
