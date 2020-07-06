# This file is run by dbmongo/js_test.go.

shopt -s extglob # enable exclusion of test files in wildcard

result_mapreduce=$(jsc \
  helpers/fakes.js \
  ../reduce.algo2/!(*_test*).js \
  ../common/!(*_test*).js \
  helpers/fake_emit_for_algo2.js \
  data/naf.js \
  data/objects.js \
  ../reduce.algo2/_test.js\
  2>&1)
if [ "$result_mapreduce" != 'true' ]; then
  echo "map_reduce_test: $result_mapreduce"
  exit 1
fi
