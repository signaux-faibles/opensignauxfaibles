result_raison_sociale=$(jsc ../common/*.js testing.js common/test_raison_sociale.js)
if [ "$result_raison_sociale" != 'true' ]; then
  exit 1
fi

