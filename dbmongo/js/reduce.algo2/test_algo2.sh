#!/bin/bash

# Usage: ./test_algo2.sh [--update]

npx ava ./test_algo2_tests.ts

# TODO:
# if [ "$1" == "--update" ]; then
#   cp ${TMP_PATH}/map_stdout.log ${TMP_PATH}/map_golden.log
#   cp ${TMP_PATH}/finalize_stdout.log ${TMP_PATH}/finalize_golden.log
#   scp ${TMP_PATH}/*_golden.log stockage:/home/centos/opensignauxfaibles_tests/
# fi
