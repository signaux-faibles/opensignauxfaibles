#!/bin/bash

# SKIP_ON_CI

# This golden-file-based test runner was designed to prevent
# regressions on the JS functions (common + algo2) used to compute the
# "Features" collection from the "RawData" collection.
# Usage: ./test_algo2.sh [--update]

# This file is run by dbmongo/js_test.go.

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
jsc \
  ./helpers/fakes.js \
  ../common/!(*_test*).js \
  ${TMP_PATH}/reduce_test_data.js \
  ./data/naf.js \
  ../reduce.algo2/!(*_test*).js \
  ../reduce.algo2/map_test.js \
  2>&1 \
  > ${TMP_PATH}/map_stdout.log

# Run tests
jsc \
  ./helpers/fakes.js \
  ../common/!(*_test*).js \
  ${TMP_PATH}/reduce_test_data.js \
  ./data/naf.js \
  ../reduce.algo2/!(*_test*).js \
  ../reduce.algo2/finalize_test.js \
  2>&1 \
  > ${TMP_PATH}/finalize_stdout.log

if [ "$1" == "--update" ]; then
  cp ${TMP_PATH}/map_stdout.log ${TMP_PATH}/map_golden.log
  cp ${TMP_PATH}/finalize_stdout.log ${TMP_PATH}/finalize_golden.log
  scp ${TMP_PATH}/*_golden.log stockage:/home/centos/opensignauxfaibles_tests/
fi

# compare *_stdout.log with golden files, return` non-zero exit code if any difference is found
DIFF1=$(diff ${TMP_PATH}/map_golden.log ${TMP_PATH}/map_stdout.log 2>&1)
DIFF2=$(diff ${TMP_PATH}/finalize_golden.log ${TMP_PATH}/finalize_stdout.log 2>&1)
if [ "${DIFF1}${DIFF2}" != "" ]; then
  echo "Test failed, because of diff: ${DIFF1}${DIFF2}"
  echo "If this diff was expected, update the golden file on server by running ./test_algo2.sh --update"
  exit 1
fi
``
echo "✅ Test passed"
exit 0
