#!/bin/bash

# SKIP_ON_CI

# This golden-file-based test runner was designed to prevent
# regressions on the JS functions (common + algo2) used to compute the
# "Features" collection from the "RawData" collection.
# Usage: ./test_map_reduce_algo2.sh [--update]

# This file is run by dbmongo/js_test.go.

# Stop script and exit on error of any command
set -e
# enable exclusion of test files in wildcard
shopt -s extglob

# Download realistic data set
TMP_PATH="./test_data_algo2"
mkdir ${TMP_PATH}
# Clean up on exit
trap "{ rm -rf ${TMP_PATH}; echo \"Cleaned up temp directory\"; }" EXIT
scp stockage:/home/centos/opensignauxfaibles_tests/* ${TMP_PATH}/

# Prepare test data set
JSON_TEST_DATASET="$(cat ./test_data_algo2/reduce_test_data.json)"
echo "makeTestData = ({ ISODate, NumberInt }) => (${JSON_TEST_DATASET});" \
  > ${TMP_PATH}/reduce_test_data.js

# Run tests
jsc ${TMP_PATH}/reduce_test_data.js ../common/!(*_test).js ../reduce.algo2/*.js ./test_map_reduce_algo2.js \
  > ${TMP_PATH}/map_stdout.log

if [ "$1" == "--update" ]; then
  cp ${TMP_PATH}/map_stdout.log ${TMP_PATH}/map_golden.log
  scp ${TMP_PATH}/map_golden.log stockage:/home/centos/opensignauxfaibles_tests/
fi

# compare map_stdout.log with golden file, return non-zero exit code if any difference is found
DIFF=$(diff ${TMP_PATH}/map_stdout.log ${TMP_PATH}/map_golden.log)
if [ "${DIFF}" != "" ]; then
  echo "Test failed, because of diff: ${DIFF}"
  echo "If the diff is expected, update the golden file on server by running ./test_map_reduce_algo2.js --update"
  exit 1
fi

echo "✅ Test passed"
exit 0
