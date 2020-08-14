#!/bin/bash

# Test de bout en bout de l'API "reduce" à l'aide de données publiques.
# Inspiré de test-api-reduce-2.sh et algo2_tests.ts.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-api-reduce.output.json"
GOLDEN_FILE="tests/output-snapshots/test-api-reduce.golden.json"
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
tests/helpers/populate-from-objects.sh \
  | tests/helpers/mongodb-container.sh run

echo ""
echo "💎 Computing the Features collection thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/reduce 👉 $(http --print=b --ignore-stdin :5000/api/data/reduce algo=algo2 batch=1905)"

(tests/helpers/mongodb-container.sh run \
  > "${OUTPUT_FILE}" \
) <<< 'printjson(db.Features_TestData.find().toArray());'

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

echo ""
# Check if the --update flag was passed
if [[ "$*" == *--update* ]]
then
    echo "🖼  Updating golden master file..."
    cp "${OUTPUT_FILE}" "${GOLDEN_FILE}"
else
    # Diff between expected and actual output
    diff --brief "${GOLDEN_FILE}" "${OUTPUT_FILE}"
    echo "✅ No diff. The reduce API works as usual."
fi

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
