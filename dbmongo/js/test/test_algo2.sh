result_add=$(jsc ../reduce.algo2/lookAhead.js ../reduce.algo2/add.js algo2/testing.js algo2/add_test.js)
if [ "$result_add" != 'true' ]; then
  exit 1
fi

result_lookAhead=$(jsc ../reduce.algo2/lookAhead.js algo2/testing.js algo2/lookAhead_test.js)
if [ "$result_lookAhead" != 'true' ]; then
  exit 1
fi

result_cibleApprentissage=$(jsc ../reduce.algo2/lookAhead.js ../reduce.algo2/cibleApprentissage.js algo2/testing.js algo2/cibleApprentissage_test.js)
if [ "$result_cibleApprentissage" != 'true' ]; then
  exit 1
fi

result_mapreduce=$(jsc ../reduce.algo2/*.js algo2/lib_algo2.js algo2/naf.js objects.js algo2/test_algo2.js)
if [ "$result_mapreduce" != 'true' ]; then
  exit 1
fi