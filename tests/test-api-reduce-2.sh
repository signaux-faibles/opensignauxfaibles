#!/bin/bash

# Test de bout en bout de l'API "reduce" à l'aide de données réalistes.
# Inspiré de test-api.sh et finalize_test.js.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.
#
# To update golden files: `$ ./test-api-reduce-2.sh --update`
# 
# These tests require the presence of private files => Make sure to:
# - run `$ git secret reveal` before running these tests;
# - run `$ git secret hide` (to encrypt changes) after updating.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
GOLDEN_FILE="tests/output-snapshots/reduce-Features.golden.json"
TMP_DIR="tests/tmp-test-execution-files"
mkdir -p "${TMP_DIR}"

# Clean up on exit
function teardown {
    tests/helpers/dbmongo-server.sh stop || true # keep tearing down, even if "No matching processes belonging to you were found"
    tests/helpers/mongodb-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start

MONGODB_PORT="27016" tests/helpers/dbmongo-server.sh setup

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
cat > "${TMP_DIR}/db_popul.js" << CONTENTS
  db.Admin.remove({})
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

  db.Features_TestData.remove({})

  db.RawData.remove({})
  db.RawData.insertMany(
CONTENTS
cat >> "${TMP_DIR}/db_popul.js" < tests/input-data/RawData.sample.json
echo ")" >> "${TMP_DIR}/db_popul.js"

tests/helpers/mongodb-container.sh run < "${TMP_DIR}/db_popul.js" >/dev/null

echo ""
echo "💎 Computing the Features collection thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/reduce 👉 $(http --print=b --ignore-stdin :5000/api/data/reduce algo=algo2 batch=2002_1)"

(tests/helpers/mongodb-container.sh run \
  | tests/helpers/remove-random_order.sh \
  > test-api-2.output.json \
) <<< 'db.Features_TestData.find().toArray();'

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file..."
    cp test-api-2.output.json "${GOLDEN_FILE}"
    echo "ℹ️  Updated ${GOLDEN_FILE} => run: $ git secret hide" # to re-encrypt the golden master file, after having updated it
else
    # Diff between expected and actual output
    diff --brief "${GOLDEN_FILE}" test-api-2.output.json
    echo "✅ No diff. The reduce API works as usual."
fi
echo ""
rm test-api-2.output.json
rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
