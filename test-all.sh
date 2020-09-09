#!/bin/bash

# Usage:
# $ git secret reveal                 # pour déchiffrer les données utilisées par les tests (golden files, etc...)
# $ ./test-all.sh                     # pour éxecuter tous les tests
# $ ./test-all.sh --update-snapshots  # pour éxecuter tous les tests et mettre à jour les snapshots des tests "ava"
# $ ./test-all.sh --update            # pour éxecuter tous les tests et mettre à jour les golden files des tests de bout en bout
# $ git secret changes                # pour visualiser les modifications éventuellement apportées aux golden files
# $ git secret hide                   # pour chiffrer les golden files suite à leur modification

FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument

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
(cd ./dbmongo/js && npm run lint && npm test -- "${FLAGS}") 2>&1 | indent

heading "go test"
(cd ./dbmongo && go test ./...) 2>&1 | indent && \

heading "go generate"
(cd ./dbmongo/lib/engine && go generate .) 2>&1 | indent

heading "go build"
(killall dbmongo 2>/dev/null || true; cd ./dbmongo && go build && echo "📦 dbmongo/dbmongo") 2>&1 | indent

heading "test-api.sh"
./tests/test-api.sh "${FLAGS}" 2>&1 | indent

heading "test-api-validate.sh"
./tests/test-api.sh "${FLAGS}" 2>&1 | indent

heading "test-api-public.sh"
./tests/test-api-public.sh "${FLAGS}" 2>&1 | indent

heading "test-api-reduce.sh"
./tests/test-api-reduce.sh "${FLAGS}" 2>&1 | indent

heading "test-api-reduce-2.sh"
./tests/test-api-reduce-2.sh "${FLAGS}" 2>&1 | indent

heading "test-api-export.sh"
./tests/test-api-export.sh "${FLAGS}" 2>&1 | indent

heading "test-api-swagger.sh"
./tests/test-api-swagger.sh "${FLAGS}" 2>&1 | indent

# Check if the --update flag was passed
if [[ "${FLAGS}" == *--update* ]]
then
    echo "ℹ️  Golden master files were updated => you may have to run: $ git secret hide" # to re-encrypt the golden master files
fi
