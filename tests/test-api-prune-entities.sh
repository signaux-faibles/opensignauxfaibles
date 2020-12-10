#!/bin/bash

# Test de bout en bout de l'API "/data/pruneEntities". Inspiré de test-api-purge-batch.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
TMP_DIR="tests/tmp-test-execution-files"
FILTER_FILE="${TMP_DIR}/test-api-prune-entities.filter.csv"
OUTPUT_FILE="${TMP_DIR}/test-api-prune-entities.output.json"
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
echo "333333333" > "${FILTER_FILE}"
echo "111111111" >> "${FILTER_FILE}"
tests/helpers/mongodb-container.sh run >/dev/null << CONTENT
  db.RawData.insertMany([
    {
      _id: "111111111",
      value: { key: "111111111", scope: "entreprise", batch: { "2007": {} } },
    },
    {
      _id: "11111111100000",
      value: { key: "11111111100000", scope: "etablissement", batch: { "2007": {} } },
    },
    {
      _id: "222222222",
      value: { key: "222222222", scope: "entreprise", batch: { "2007": {} } },
    },
    {
      _id: "22222222200000",
      value: { key: "22222222200000", scope: "etablissement", batch: { "2007": {} } },
    },
    {
      _id: "333333333",
      value: { key: "333333333", scope: "entreprise", batch: { "2007": {} } },
    },
  ]);
  db.Admin.insertOne({
    _id: { key: "2010", type: "batch" },
    files: {
      filter: [ "/../../${FILTER_FILE}" ],
    },
  });
CONTENT

echo ""
echo "💎 Test: count and prune entities from RawData..."
tests/helpers/dbmongo-server.sh start
API_RESULT=$(http --print=b --ignore-stdin :5000/api/data/pruneEntities batch=2010)
echo "- POST /api/data/pruneEntities 👉 ${API_RESULT}"

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
    "found 2 entities to prune": ${API_RESULT}.count === 2,
    "222222222 was not pruned yet": db.RawData.find({_id: "222222222"}).count() === 1,
    "22222222200000 was not pruned yet": db.RawData.find({_id: "22222222200000"}).count() === 1,
    "Journal has 1 entry": db.Journal.count() === 1,
    "Journal reports PurgeBatch": report.reportType === "PruneEntities",
    "Journal report has date": !!report.date === true,
    "Journal report has start date": !!report.startDate === true,
  }).forEach(([ testName, testRes ]) => print(testName, ':', testRes));
CONTENT

echo "- POST /api/data/pruneEntities delete=true 👉 $(http --print=b --ignore-stdin :5000/api/data/pruneEntities batch=2010 delete:=true)"

(tests/helpers/mongodb-container.sh run \
  | reportFailedTests \
) << CONTENT
  Object.entries({
    "333333333 was not pruned": db.RawData.find({_id: "333333333"}).count() === 1,
    "111111111 was not pruned": db.RawData.find({_id: "111111111"}).count() === 1,
    "11111111100000 was not pruned": db.RawData.find({_id: "11111111100000"}).count() === 1,
    "222222222 was pruned": db.RawData.find({_id: "222222222"}).count() === 0,
    "22222222200000 was pruned": db.RawData.find({_id: "22222222200000"}).count() === 0,
  }).forEach(([ testName, testRes ]) => print(testName, ':', testRes));
CONTENT

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
