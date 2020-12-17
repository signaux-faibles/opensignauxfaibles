#!/bin/bash

# Test de bout en bout de la commande "reduce" à l'aide de données publiques.
# Inspiré de test-reduce-2.sh et algo2_tests.ts.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-api-reduce.output.json"
GOLDEN_FILE="tests/output-snapshots/test-api-reduce.golden.json"
mkdir -p "${TMP_DIR}"

# Clean up on exit
function teardown {
    tests/helpers/sfdata-wrapper.sh stop || true # keep tearing down, even if "No matching processes belonging to you were found"
    tests/helpers/mongodb-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start

MONGODB_PORT="27016" tests/helpers/sfdata-wrapper.sh setup

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
tests/helpers/populate-from-objects.sh \
  | tests/helpers/mongodb-container.sh run

# We create a collection with dummy data which should not remain after the execution of Reduce
echo "db.Features_TestData.insertOne({a:1})" | tests/helpers/mongodb-container.sh run

echo ""
echo "💎 Computing the Features collection..."
RESULT=$(tests/helpers/sfdata-wrapper.sh run reduce --until-batch=1905)
echo "- sfdata reduce 👉 ${RESULT}"

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

print("// Documents from db.Features_TestData:");
printjson(db.Features_TestData.find().toArray());

print("// Response body from sfdata reduce:");
CONTENT

echo "${RESULT}" >> "${OUTPUT_FILE}"

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
