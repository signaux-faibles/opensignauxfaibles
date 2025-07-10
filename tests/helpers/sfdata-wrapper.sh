#!/bin/bash

# This helper runs the "sfdata" command using provided MongoDB port,
# and prefixes its output, to make warnings stand out during tests.

# Restore config on exit (including if errors happen)
function restoreConfig {
  [ -f config.backup.toml ] && mv config.backup.toml config.toml
}
trap restoreConfig EXIT

# Add prefix to stderr stream
exec 2> >(sed 's/^/  [sfdata] /' >&2)

# Set MONGODB_PORT in config.toml
[ -f config.toml ] && mv config.toml config.backup.toml
cp config-sample.toml config.toml
perl -pi'' -e "s,27017,${MONGODB_PORT}," config.toml
perl -pi'' -e "s,path/export,${TMP_DIR}," config.toml


# Run the command
./sfdata "$@" || true # pass all arguments to sfdata

# (at the end, trap will restore config.toml)
