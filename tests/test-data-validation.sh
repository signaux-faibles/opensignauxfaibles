#!/bin/bash

# Test de la validation de données de MongoDB.
# Inspiré de test-api-reduce-2.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.
#
# To update golden files: `$ ./test-data-validation.sh --update`
# 
# These tests require the presence of private files => Make sure to:
# - run `$ git secret reveal` before running these tests;
# - run `$ git secret hide` (to encrypt changes) after updating.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
mkdir -p "${TMP_DIR}"

# Clean up on exit
function teardown {
    tests/helpers/mongodb-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
tests/helpers/mongodb-container.sh run << CONTENT
  db.RawData.insertMany($(cat tests/input-data/RawData.validation.json))
CONTENT

echo ""
echo "💎 Validating data..."
AGGREG_PREPARATION='
    { $project: { _id: 1, batches: { $objectToArray: "$value.batch" } } }, // => { _id, batches: Array<{ k: BatchKey, v: BatchValues }> }
    { $unwind: { path: "$batches", preserveNullAndEmptyArrays: false } }, // => { _id, batches: { k: BatchKey, v: BatchValues } }
    { $project: { _id: 1, batchKey: "$batches.k", "dataPerHash": { $objectToArray: "$batches.v" } } }, // => { _id, batchKey, dataPerHash: Array<{ k: DataType, v: ParHash<Data> }> }
    { $unwind: { path: "$dataPerHash", preserveNullAndEmptyArrays: false } }, // => { _id, batchKey, dataPerHash: { k: DataType, v: ParHash<Data> } }
    { $project: { _id: 1, batchKey: 1, dataType: "$dataPerHash.k", "dataPerHash": { $objectToArray: "$dataPerHash.v" } } }, // => { _id, batchKey, dataType, dataPerHash: Array<{ k: Hash, v: Data }> }
    { $unwind: { path: "$dataPerHash", preserveNullAndEmptyArrays: false } }, // => { _id, batchKey, dataPerHash: { k: DataType, v: ParHash<Data> } }, // => { _id, batchKey, dataType, dataPerHash: { k: Hash, v: Data } }
    { $project: { _id: 1, batchKey: 1, dataType: 1, dataHash: "$dataPerHash.k", "dataObject": "$dataPerHash.v" } }, // => { _id, batchKey, dataType, dataHash, dataObject: Data }
'
tests/helpers/mongodb-container.sh run << CONTENT
  printjson(db.RawData.aggregate([
    ${AGGREG_PREPARATION}
    {
      \$facet: {
        "valid": [
          {
            \$match: {
              \$jsonSchema: {
                bsonType: "object",
                properties: {
                  dataObject: $(cat dbmongo/validation/delai.schema.json)
                }
              }
            }
          },
          {
            \$project: {
              _id: 0,
              dataType: 1,
              dataHash: 1
            }
          }
        ],
        "invalid": [
          {
            \$match: {
              \$nor: [
                {
                  \$jsonSchema: {
                    bsonType: "object",
                    properties: {
                      dataObject: $(cat dbmongo/validation/delai.schema.json)
                    }
                  }
                }
              ]
            }
          },
          {
            \$project: {
              _id: 0,
              dataType: 1,
              dataHash: 1
            }
          }
        ]
      }
    },
  ]).toArray()[0])
CONTENT

# Display JS errors logged by MongoDB, if any
tests/helpers/mongodb-container.sh exceptions || true

# tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
