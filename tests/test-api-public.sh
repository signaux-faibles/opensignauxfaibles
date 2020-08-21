#!/bin/bash

# Test de bout en bout de l'API "/public". Inspiré de test-api-public.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-api-public.output.json"
GOLDEN_FILE="tests/output-snapshots/test-api-public.golden.json"
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
  | tests/helpers/mongodb-container.sh run >/dev/null

echo ""
echo "💎 Computing the Public collection thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/public 👉 $(http --print=b --ignore-stdin :5000/api/data/public batch=1905)"

(tests/helpers/mongodb-container.sh run \
  > "${OUTPUT_FILE}" \
) <<< 'printjson(db.Public.find().toArray());'

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
