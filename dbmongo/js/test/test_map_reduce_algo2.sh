mkdir ./test_data_algo2
scp stockage:/home/centos/opensignauxfaibles_tests/reduce_test_data.json ./test_data_algo2/

# Tests here

echo "ISODate = (dateString) => new Date(dateString); NumberInt = (int) => int; testData = $(cat ./test_data_algo2/reduce_test_data.json)" > ./test_data_algo2/reduce_test_data.js
cat ../common/*.js >./test_data_algo2/jsFunctions.js
cat ../reduce.algo2/*.js >>./test_data_algo2/jsFunctions.js
jsc ./test_data_algo2/reduce_test_data.js ./test_data_algo2/jsFunctions.js ./test_map_reduce_algo2.js > ./test_data_algo2/stdout.log

# Clean up
# rm -rf ./test_data_algo2
