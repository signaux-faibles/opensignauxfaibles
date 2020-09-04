#!/bin/bash

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
DATABASE="mongodb://localhost:27018/test"

mongo "${DATABASE}" <<< "db.RawData.count()" # => 6929282 documents

echo ""
echo "💎 Validating data..."
AGGREG_PREPARATION='
    { $limit: 1000000 }, // on ne traite que les 1000000 premiers documents de RawData (TODO: à retirer)
    { $project: { _id: 1, batches: { $objectToArray: "$value.batch" } } }, // => { _id, batches: Array<{ k: BatchKey, v: BatchValues }> }
    { $unwind: { path: "$batches", preserveNullAndEmptyArrays: false } }, // => { _id, batches: { k: BatchKey, v: BatchValues } }
    { $project: { _id: 1, batchKey: "$batches.k", "dataPerHash": { $objectToArray: "$batches.v" } } }, // => { _id, batchKey, dataPerHash: Array<{ k: DataType, v: ParHash<Data> }> }
    { $unwind: { path: "$dataPerHash", preserveNullAndEmptyArrays: false } }, // => { _id, batchKey, dataPerHash: { k: DataType, v: ParHash<Data> } }
    { $project: { _id: 1, batchKey: 1, dataType: "$dataPerHash.k", "dataPerHash": { $objectToArray: "$dataPerHash.v" } } }, // => { _id, batchKey, dataType, dataPerHash: Array<{ k: Hash, v: Data }> }
    { $unwind: { path: "$dataPerHash", preserveNullAndEmptyArrays: false } }, // => { _id, batchKey, dataPerHash: { k: DataType, v: ParHash<Data> } }, // => { _id, batchKey, dataType, dataPerHash: { k: Hash, v: Data } }
    { $project: { _id: 1, batchKey: 1, dataType: 1, dataHash: "$dataPerHash.k", "dataObject": "$dataPerHash.v" } }, // => { _id, batchKey, dataType, dataHash, dataObject: Data }
'
cat > "dbvalid.js" << CONTENT
  print("invalid records:");
  db.RawData.aggregate([
    ${AGGREG_PREPARATION}
    {
            \$match: {
              dataType: "delai",
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
    }
  ]).forEach(printjson);
  print("done.");
CONTENT

cat "dbvalid.js"
time mongo "${DATABASE}" "dbvalid.js"
