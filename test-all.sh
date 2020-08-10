#!/bin/bash

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
(cd ./dbmongo/js && npm run lint && npm test) 2>&1 | indent

heading "go test"
(cd ./dbmongo && go test ./...) 2>&1 | indent && \

heading "go generate"
(cd ./dbmongo/lib/engine && go generate .) 2>&1 | indent

heading "go build"
(killall dbmongo 2>/dev/null || true; cd ./dbmongo && go build && echo "📦 dbmongo/dbmongo") 2>&1 | indent

heading "test-api.sh"
./tests/test-api.sh 2>&1 | indent

heading "test-api-reduce.sh"
./tests/test-api-reduce.sh 2>&1 | indent

heading "test-api-reduce-2.sh"
./tests/test-api-reduce-2.sh 2>&1 | indent
