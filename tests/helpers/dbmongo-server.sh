#!/bin/bash

# This helper can configure, start and stop the "dbmongo" API server.

COMMAND=$1

# Add prefix to stderr stream
exec 2> >(sed 's/^/[API] /' >&2)

case ${COMMAND} in
  stop)
    killall dbmongo >/dev/null
    [ -f dbmongo/config.backup.toml ] && mv dbmongo/config.backup.toml dbmongo/config.toml
    exit ;;
  setup)
    [ -f dbmongo/config.toml ] && mv dbmongo/config.toml dbmongo/config.backup.toml
    cp dbmongo/config-sample.toml dbmongo/config.toml
    perl -pi'' -e "s,/foo/bar/data-raw,sample-data-raw," dbmongo/config.toml
    perl -pi'' -e "s,27017,${MONGODB_PORT}," dbmongo/config.toml
    exit ;;
  start)
    cd dbmongo
    bash -c "./dbmongo &>/dev/null &" # we run in a separate shell to hide the "terminated" message when the process is killed by trap
    sleep 2 # give some time for dbmongo to start
    exit ;;
  help) ;;
  ?) echo "error: ${COMMAND} is not a recognized command" ;;
esac

echo "usage: $0 stop|setup|start|help"
