# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

result_add=$(jsc \
  helpers/testing.js \
  helpers/fakes.js \
  helpers/algo2_reducer.js \
  ../common/!(*_test).js \
  ../reduce.algo2/!(*_test).js \
  ../reduce.algo2/add_test.js \
  2>&1)
if [ "$result_add" != 'true' ]; then
  echo "$result_add"
  exit 1
fi

result_lookAhead=$(jsc \
  helpers/testing.js \
  helpers/fakes.js \
  helpers/algo2_reducer.js \
  ../common/!(*_test).js \
  ../reduce.algo2/!(*_test).js \
  ../reduce.algo2/lookAhead_test.js \
  2>&1)
if [ "$result_lookAhead" != 'true' ]; then
  echo "$result_lookAhead"
  exit 1
fi

result_cibleApprentissage=$(jsc \
  helpers/testing.js \
  helpers/fakes.js \
  helpers/algo2_reducer.js \
  ../common/!(*_test).js \
  ../reduce.algo2/!(*_test).js \
  ../reduce.algo2/cibleApprentissage_test.js \
  2>&1)
if [ "$result_cibleApprentissage" != 'true' ]; then
  echo "$result_cibleApprentissage"
  exit 1
fi

# TODO pourquoi ce test est commenté ?
# result_mapreduce=$(jsc \
#   ../reduce.algo2/!(*_test).js \
#   ../common/!(*_test).js \
#   helpers/fakes.js \
#   helpers/algo2_reducer.js \
#   data/naf.js \
#   data/objects.js \
#   ./reduce.algo2/_test.js\
#   2>&1)
# if [ "$result_mapreduce" != 'true' ]; then
#   echo "$result_mapreduce"
#   exit 1
# fi
