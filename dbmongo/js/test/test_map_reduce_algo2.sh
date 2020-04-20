# This golden-file-based test runner was designed to prevent
# regressions on the JS functions (common + algo2) used to compute the
# "Features" collection from the "RawData" collection.

# Download realistic data set
mkdir ./test_data_algo2
scp stockage:/home/centos/opensignauxfaibles_tests/reduce_test_data.json ./test_data_algo2/

# Run tests
echo "makeTestData = ({ ISODate, NumberInt }) => ($(cat ./test_data_algo2/reduce_test_data.json));" > ./test_data_algo2/reduce_test_data.js
jsc ./test_data_algo2/reduce_test_data.js ../common/*.js ../reduce.algo2/*.js ./test_map_reduce_algo2.js > ./test_data_algo2/stdout.log
cat ./test_data_algo2/stdout.log

# TODO: compare stdout.log with golden file, return non-zero exit code if any difference is found

# Clean up
rm -rf ./test_data_algo2
