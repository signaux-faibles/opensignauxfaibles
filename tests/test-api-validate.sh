#!/bin/bash

# Test de bout en bout de POST /api/data/validate.
# Inspiré de test-api-export.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
COLOR_YELLOW='\033[1;33m'
COLOR_DEFAULT='\033[0m'
INPUT_FILE="tests/input-data/RawData.validation.json"
GOLDEN_FILE="tests/output-snapshots/test-data-validation.golden.json"
TMP_DIR="tests/tmp-test-execution-files"
mkdir -p "${TMP_DIR}"

# Clean up on exit
function teardown {
    echo -e "${COLOR_DEFAULT}"
    tests/helpers/dbmongo-server.sh stop || true # keep tearing down, even if "No matching processes belonging to you were found"
    tests/helpers/mongodb-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start

MONGODB_PORT="27016" tests/helpers/dbmongo-server.sh setup

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
tests/helpers/mongodb-container.sh run << CONTENT
  db.RawData.insertMany($(cat ${INPUT_FILE}))
CONTENT

echo ""
echo "💎 Computing the Public collection thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
OUTPUT_FILE=dbmongo/$(http --print=b --ignore-stdin :5000/api/data/validate collection=RawData | tr -d '"')
echo "- POST /api/data/validate 👉 ${OUTPUT_FILE}"

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm "${OUTPUT_FILE}"
rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
