#!/bin/bash

# This golden-file-based test runner was designed to prevent
# regressions on the JS functions (common + algo2) used to compute the
# "Features" collection from the "RawData" collection.

# Download realistic data set
TMP_PATH="./test_data_algo2"
mkdir ${TMP_PATH}
scp stockage:/home/centos/opensignauxfaibles_tests/* ${TMP_PATH}/

# Prepare test data set
JSON_TEST_DATASET="$(cat ./test_data_algo2/reduce_test_data.json)"
echo "makeTestData = ({ ISODate, NumberInt }) => (${JSON_TEST_DATASET});" \
  > ${TMP_PATH}/reduce_test_data.js

# Run tests
jsc ${TMP_PATH}/reduce_test_data.js ../common/*.js ../reduce.algo2/*.js ./test_map_reduce_algo2.js \
  > ${TMP_PATH}/map_stdout.log

if [ "$1" == "--update" ]; then
  cp ${TMP_PATH}/map_stdout.log ${TMP_PATH}/map_golden.log
  scp ${TMP_PATH}/map_golden.log stockage:/home/centos/opensignauxfaibles_tests/
fi

# compare map_stdout.log with golden file, return non-zero exit code if any difference is found
DIFF=$(diff ${TMP_PATH}/map_stdout.log ${TMP_PATH}/map_golden.log)
if [ "${DIFF}" != "" ]; then
  echo "Test failed, because of diff: ${DIFF}"
  exit 1
fi

# Clean up
# TODO: make sure to delete the file in the case of error.
rm -rf ${TMP_PATH}
