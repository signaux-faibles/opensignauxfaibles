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

heading "pick specified node.js version" && (. ~/.nvm/nvmc.sh; cd ./dbmongo/js && nvm use) | indent

# Mandatory tests (can stop the script)

set -e # will stop the script if any command fails with a non-zero exit code
set -o pipefail # ... even for tests which pipe their output to indent

heading "npm install" && (cd ./dbmongo/js && npm install) | indent
heading "npm test" && (cd ./dbmongo/js && npm run lint && npm test) | indent
heading "go test" && (cd ./dbmongo && go test ./...) | indent && \
heading "go generate" && (cd ./dbmongo/lib/engine && go generate .) | indent
heading "go build" && (killall dbmongo 2>/dev/null; cd ./dbmongo && go build) | indent
heading "test-api.sh" && ./tests/test-api.sh | indent
heading "test-api-reduce.sh" && ./tests/test-api-reduce.sh | indent
heading "test-api-reduce-2.sh" && ./tests/test-api-reduce-2.sh | indent
