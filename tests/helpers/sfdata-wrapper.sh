#!/bin/bash

# This helper can configure, start and stop the "sfdata" command.

COMMAND=$1

# Add prefix to stderr stream
exec 2> >(sed 's/^/[API] /' >&2)

case ${COMMAND} in
  stop)
    killall sfdata >/dev/null
    [ -f config.backup.toml ] && mv config.backup.toml config.toml
    exit ;;
  setup)
    [ -f config.toml ] && mv config.toml config.backup.toml
    cp config-sample.toml config.toml
    perl -pi'' -e "s,/foo/bar/data-raw,sample-data-raw," config.toml
    perl -pi'' -e "s,27017,${MONGODB_PORT}," config.toml
    exit ;;
  run)
    ./sfdata "${@:2}" || true # pass all arguments except the first one
    # TODO: call sfdata directly from tests, if possible
    exit ;;
  ?)
    echo "error: ${COMMAND} is not a recognized command" ;;
esac

echo "usage: $0 stop|setup|start"
