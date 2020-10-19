#!/bin/bash

# Test de bout en bout de l'API "/batch/purge". Inspiré de test-api-public.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-api-purge-batch.output.json"
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
tests/helpers/mongodb-container.sh run >/dev/null << CONTENT
  $(tests/helpers/populate-from-objects.sh)
  db.Admin.insertOne({ '_id' : { 'key' : '1901', 'type' : 'batch' } })
CONTENT

echo ""
echo "💎 Test: purge batch 1901 from RawData..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/batch/purge 👉 $(http --print=b --ignore-stdin :5000/api/data/batch/purge fromBatch=1901 IUnderstandWhatImDoing:=true)"

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

(tests/helpers/mongodb-container.sh run \
  > "${OUTPUT_FILE}" \
) <<< 'printjson({
    "1901 was purged": db.RawData.find({"value.batch.1901": {"$exists": true}}).count() === 0,
    "1812 was not purged": db.RawData.find({"value.batch.1812": {"$exists": true}}).count() > 0,
  });'

cat "${OUTPUT_FILE}"

grep --quiet '{ "1901 was purged" : true, "1812 was not purged" : true }' "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
