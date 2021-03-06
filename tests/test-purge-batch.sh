#!/bin/bash

# Test de bout en bout de la commande "purge". Inspiré de test-public.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
TMP_DIR="tests/tmp-test-execution-files"
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
tests/helpers/mongodb-container.sh run >/dev/null << CONTENT
  $(tests/helpers/populate-from-objects.sh)
  db.Admin.insertOne({ '_id' : { 'key' : '1901', 'type' : 'batch' } })
CONTENT

echo ""
echo "💎 Test: purge batch 1901 from RawData..."
echo "- sfdata batch/purge 👉 $(tests/helpers/sfdata-wrapper.sh purge --since-batch=1901 --i-understand-what-im-doing)"

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

# Print test results from stdin. Fails on any "false" result.
# Expected format for each line: "<test label> : <true|false>"
function reportFailedTests {
  while IFS='$\n' read -r line; do
    echo "  - $line" | (grep --color=always " : false") || true # display failed test
    echo "  - $line" | grep " : true" # display passing test, and make the test function fail otherwise
  done
}

(tests/helpers/mongodb-container.sh run \
  | reportFailedTests \
) << CONTENT
  const report = db.Journal.find().toArray().pop() || {};
  Object.entries({
    "1901 was purged": db.RawData.find({"value.batch.1901": {"\$exists": true}}).count() === 0,
    "1812 was not purged": db.RawData.find({"value.batch.1812": {"\$exists": true}}).count() > 0,
    "Journal has 1 entry": db.Journal.count() === 1,
    "Journal reports PurgeBatch": report.reportType === "PurgeBatch",
    "Journal report has date": !!report.date === true,
    "Journal report has start date": !!report.startDate === true,
  }).forEach(([ testName, testRes ]) => print(testName, ':', testRes));
CONTENT

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
