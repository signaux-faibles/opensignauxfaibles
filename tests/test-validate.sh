#!/bin/bash

# Test de bout en bout de la commande "validate".
# Inspiré de test-export.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
INPUT_FILE="tests/input-data/RawData.validation.json"
GOLDEN_FILE="tests/output-snapshots/test-validate.golden.json"
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/output.json"
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
tests/helpers/mongodb-container.sh run << CONTENT
  db.RawData.insertMany($(cat ${INPUT_FILE}))
CONTENT

echo ""
echo "💎 Testing sfdata..."
VALIDATION_REPORT=$(tests/helpers/sfdata-wrapper.sh validate --collection=RawData)
echo "- sfdata validate"

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

print("// Result from sfdata validate:");
CONTENT

echo "${VALIDATION_REPORT}" >> "${OUTPUT_FILE}"

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
