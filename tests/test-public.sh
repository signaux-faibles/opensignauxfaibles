#!/bin/bash

# Test de bout en bout de la commande "public". Inspiré de test-public.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-public.output.json"
GOLDEN_FILE="tests/output-snapshots/test-public.golden.json"
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
tests/helpers/populate-from-objects.sh \
  | tests/helpers/mongodb-container.sh run >/dev/null

echo ""
echo "💎 Computing the Public collection..."
RESULT=$(tests/helpers/sfdata-wrapper.sh public --until-batch=1905)
echo "- sfdata public 👉 ${RESULT}"

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

print("// Documents from db.Public:");
printjson(db.Public.find().toArray());

print("// Response body from sfdata public:");
CONTENT

echo "${RESULT}" >> "${OUTPUT_FILE}"

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
