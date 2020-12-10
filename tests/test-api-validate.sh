#!/bin/bash

# Test de bout en bout de POST /api/data/validate.
# Inspiré de test-api-export.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
INPUT_FILE="tests/input-data/RawData.validation.json"
GOLDEN_FILE="tests/output-snapshots/test-api-validate.golden.json"
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/output.json"
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
tests/helpers/mongodb-container.sh run << CONTENT
  db.RawData.insertMany($(cat ${INPUT_FILE}))
CONTENT

echo ""
echo "💎 Testing the dbmongo API..."
tests/helpers/dbmongo-server.sh start
API_RESULT=$(http --print=b --ignore-stdin :5000/api/data/validate collection=RawData)
echo "- POST /api/data/validate 👉 ${API_RESULT}"

(tests/helpers/mongodb-container.sh run \
  > "${OUTPUT_FILE}" \
) << CONTENT
print("// db.Journal:");
const report = db.Journal.find().toArray().pop() || {};
printjson({
  count: db.Journal.count(),
  reportType: report.reportType,
  hasDate: !!report.date,
  hasStartDate: !!report.startDate,
});

print("// Result from /api/data/validate:");
CONTENT

OUTPUT_GZ_FILE=dbmongo/$(echo ${API_RESULT} | tr -d '"')
zcat < "${OUTPUT_GZ_FILE}" >> "${OUTPUT_FILE}"

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm "${OUTPUT_GZ_FILE}"
rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
