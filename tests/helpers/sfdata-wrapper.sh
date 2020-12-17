#!/bin/bash

# This helper can run the "sfdata" command using provided MongoDB port.

COMMAND=$1

# Add prefix to stderr stream
exec 2> >(sed 's/^/[sfdata] /' >&2)

case ${COMMAND} in
  stop)
    killall sfdata &>/dev/null
    [ -f config.backup.toml ] && mv config.backup.toml config.toml
    exit ;;
  setup)
    [ -f config.toml ] && mv config.toml config.backup.toml
    cp config-sample.toml config.toml
    perl -pi'' -e "s,27017,${MONGODB_PORT}," config.toml # TODO: pass DB PORT to Viper as an environment variable => call sfdata directly from tests => delete sfdata-wrapper.sh
    exit ;;
  run)
    ./sfdata "${@:2}" || true # pass all arguments except the first one
    exit ;;
  ?)
    echo "error: ${COMMAND} is not a recognized command" ;;
esac

echo "usage: $0 stop|setup|run"
