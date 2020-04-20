# mkdir ./test_data_algo2
# scp stockage:/home/centos/opensignauxfaibles_tests/reduce_test_data.json ./test_data_algo2/

# Tests here
cat ../reduce.algo2/*.js >./test_data_algo2/jsFunctions.js
cat ../common/*.js >>./test_data_algo2/jsFunctions.js
jsc ./test_data_algo2/jsFunctions.js ./test_map_reduce_algo2.js

# Clean up
# rm -rf ./test_data_algo2
