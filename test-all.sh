#!/bin/bash

set -o pipefail # so any failing test can stop the script despite piping to indent

function heading {
  echo ""
  echo "–––––"
  echo "$1"
  echo "–––––"
}

function indent {
  sed 's/^/  /'
}

heading "pick specified node.js version (optional)" && (. ~/.nvm/nvm.sh && cd ./dbmongo/js && nvm use) | indent; \
heading "npm install" && (cd ./dbmongo/js && npm install) | indent && \
heading "npm test" && (cd ./dbmongo/js && npm run lint && npm test) | indent && \
heading "go test" && (cd ./dbmongo && go test ./...) | indent && \
heading "go generate" && (cd ./dbmongo/lib/engine && go generate .) | indent && \
heading "go build" && (killall dbmongo >/dev/null; cd ./dbmongo && go build) | indent && \
heading "test-api.sh" && ./tests/test-api.sh | indent && \
heading "test-api-reduce.sh" && ./tests/test-api-reduce.sh | indent && \
heading "test-api-reduce-2.sh" && ./tests/test-api-reduce-2.sh | indent
