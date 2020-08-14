#!/bin/bash

# This helper can start/stop a MongoDB container and send mongo shell instructions to it.

COMMAND=$1
CONTAINER="sf-mongodb"
IMAGE="mongo:4.2@sha256:1c2243a5e21884ffa532ca9d20c221b170d7b40774c235619f98e2f6eaec520a"
DATABASE="signauxfaibles"

case ${COMMAND} in
  stop) sudo docker stop ${CONTAINER} &>/dev/null; exit ;;
  start) sudo docker run --name ${CONTAINER} --publish ${PORT}:27017 --detach --rm ${IMAGE} >/dev/null; exit ;;
  run) sudo docker exec -i ${CONTAINER} mongo --quiet ${DATABASE}; exit ;;
  help) ;;
  ?) echo "error: ${COMMAND} is not a recognized command" ;;
esac

echo "usage: $0 stop|start|run|help"
