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

heading "go generate" && (cd dbmongo/lib/engine && go generate .) | indent && \
heading "npm test" && (cd dbmongo/js && npm run lint && npm test) | indent && \
heading "go test" && (cd dbmongo && go test ./...) | indent && \
heading "go build" && (killall dbmongo; cd dbmongo && go build) | indent && \
heading "test-api.sh" && ./test-api.sh | indent && \
heading "test-api-reduce.sh" && ./test-api-reduce.sh | indent && \
heading "test-api-2.sh" && ./test-api-2.sh | indent
