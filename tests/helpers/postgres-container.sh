#!/bin/bash

# This helper can start/stop a PostgreSQL container and send psql instructions to it.

COMMAND=$1
CONTAINER="sf-postgres"
IMAGE="postgres:17@sha256:fe3f571d128e8efadcd8b2fde0e2b73ebab6dbec33f6bfe69d98c682c7d8f7bd"

USER="testuser"
PASSWORD="testpass"
DATABASE="testdb"

# Add prefix to stderr stream
exec 2> >(sed 's/^/[DB] /' >&2)

case ${COMMAND} in
  stop)
    docker stop "${CONTAINER}" &>/dev/null
	  docker rm "${CONTAINER}" &>/dev/null
    exit ;;
  start)
    sudo docker run --name "${CONTAINER}" \
      --env "POSTGRES_DB=${DATABASE}" \
      --env "POSTGRES_USER=${USER}" \
      --env "POSTGRES_PASSWORD=${PASSWORD}" \
      --publish "${PORT:-5432}:5432" \
      --volume "$(pwd)/migrations:/migrations" \
      --detach --rm "${IMAGE}" &>/dev/null
    sleep 3
    sudo docker exec "${CONTAINER}" bash -c "for file in /migrations/*; do psql postgres://${USER}:${PASS}@localhost:5432/${DATABASE} -f \"\$file\"; done"
    exit ;;
  run)
    sudo docker exec -i "${CONTAINER}" psql "postgres://${USER}:${PASS}@localhost:5432/${DATABASE}"
    exit ;;
  ?)
    echo "error: ${COMMAND} is not a recognized command" ;;
esac

echo "usage: $0 stop|start|run"
