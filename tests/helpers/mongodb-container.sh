#!/bin/bash

# This helper can start/stop a MongoDB container and send mongo shell instructions to it.

COMMAND=$1
CONTAINER="sf-mongodb"
IMAGE="mongo:4.2@sha256:1c2243a5e21884ffa532ca9d20c221b170d7b40774c235619f98e2f6eaec520a"
DATABASE="signauxfaibles"

# Add prefix to stderr stream
exec 2> >(sed 's/^/[DB] /' >&2)

case ${COMMAND} in
  stop)
    sudo docker stop "${CONTAINER}" &>/dev/null
    exit ;;
  start)
    sudo docker run --name "${CONTAINER}" --publish "${PORT}:27017" --detach --rm "${IMAGE}" >/dev/null
    exit ;;
  run)
    sudo docker exec -i "${CONTAINER}" dd status=none of="file.js" # upload the js commands to a file in the container, because it's much faster than piping directly to mongo shell
    sudo docker exec -i "${CONTAINER}" mongo --quiet "${DATABASE}" "file.js"
    exit ;;
  exceptions)
    sudo docker logs "${CONTAINER}" | grep --color=always "uncaught exception"
    exit ;;
  ?)
    echo "error: ${COMMAND} is not a recognized command" ;;
esac

echo "usage: $0 stop|start|run|exceptions"
