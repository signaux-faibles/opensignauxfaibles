# This file is run by dbmongo/js_test.go.

set -e # will stop on any error
shopt -s extglob # enable exclusion of test files in wildcard

result_add=$(jsc \
  ../reduce.algo2/!(*_test).js \
  ../common/!(*_test).js \
  helpers/testing.js \
  algo2/lib_algo2.js \
  ../reduce.algo2/add_test.js \
  )
if [ "$result_add" != 'true' ]; then
  exit 1
fi

result_lookAhead=$(jsc \
  ../reduce.algo2/!(*_test).js \
  ../common/!(*_test).js \
  helpers/testing.js \
  algo2/lib_algo2.js \
  ../reduce.algo2/lookAhead_test.js \
  )
if [ "$result_lookAhead" != 'true' ]; then
  exit 1
fi

result_cibleApprentissage=$(jsc \
  ../reduce.algo2/!(*_test).js \
  ../common/!(*_test).js \
  helpers/testing.js \
  algo2/lib_algo2.js \
  ../reduce.algo2/cibleApprentissage_test.js \
  )
if [ "$result_cibleApprentissage" != 'true' ]; then
  exit 1
fi

# TODO pourquoi ce test est commenté ?
# result_mapreduce=$(jsc \
#   ../reduce.algo2/!(*_test).js \
#   ../common/!(*_test).js \
#   algo2/lib_algo2.js \
#   data/naf.js \
#   data/objects.js \
#   ./reduce.algo2/_test.js\
#   )
# if [ "$result_mapreduce" != 'true' ]; then
#   exit 1
# fi
