#!/bin/bash

# Test de bout en bout de la commande "reduce" à l'aide de données réalistes.
# Inspiré de test.sh et finalize_test.js.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.
#
# To update golden files: `$ ./test-reduce-2.sh --update`
# 
# These tests require the presence of private files => Make sure to:
# - run `$ git secret reveal` before running these tests;
# - run `$ git secret hide` (to encrypt changes) after updating.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/reduce-Features.output.json"
GOLDEN_FILE="tests/output-snapshots/reduce-Features.golden.json"
mkdir -p "${TMP_DIR}"

# Clean up on exit
function teardown {
    tests/helpers/mongodb-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start
export MONGODB_PORT="27016" # for tests/helpers/sfdata-wrapper.sh

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
tests/helpers/mongodb-container.sh run > /dev/null << CONTENTS
  db.Admin.insertOne({
    "_id" : {
        "key" : "2002_1",
        "type" : "batch"
    },
    "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2016-01-01T00:00:00.000+0000"),
        "date_fin_effectif" : ISODate("2016-03-01T00:00:00.000+0000")
    },
    "name" : "TestData"
  })

  db.RawData.insertMany(
    $(cat tests/input-data/RawData.sample.json)
  )
CONTENTS

echo ""
echo "💎 Computing the Features collection..."
echo "- sfdata reduce 👉 $(tests/helpers/sfdata-wrapper.sh reduce --until-batch=2002_1)"

(tests/helpers/mongodb-container.sh run \
  | tests/helpers/remove-random_order.sh \
  > "${OUTPUT_FILE}" \
) <<< 'printjson(db.Features_TestData.find().toArray());'

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
