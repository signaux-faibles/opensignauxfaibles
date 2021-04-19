#!/bin/bash

set -e # will stop the script if any command fails with a non-zero exit code

# Extract a comma-separated list of global variables that are expected by TypeScript files passed as arguments.
grep -F --no-filename 'declare const' $@ \
  | cut -d' ' -f3 \
  | cut -d':' -f1 \
  | sort -u \
  | uniq \
  | paste -sd "," -
