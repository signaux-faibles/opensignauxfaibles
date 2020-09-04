  printjson(db.RawData.aggregate([
    
    { $project: { _id: 1, batches: { $objectToArray: "$value.batch" } } }, // => { _id, batches: Array<{ k: BatchKey, v: BatchValues }> }
    { $unwind: { path: "$batches", preserveNullAndEmptyArrays: false } }, // => { _id, batches: { k: BatchKey, v: BatchValues } }
    { $project: { _id: 1, batchKey: "$batches.k", "dataPerHash": { $objectToArray: "$batches.v" } } }, // => { _id, batchKey, dataPerHash: Array<{ k: DataType, v: ParHash<Data> }> }
    { $unwind: { path: "$dataPerHash", preserveNullAndEmptyArrays: false } }, // => { _id, batchKey, dataPerHash: { k: DataType, v: ParHash<Data> } }
    { $project: { _id: 1, batchKey: 1, dataType: "$dataPerHash.k", "dataPerHash": { $objectToArray: "$dataPerHash.v" } } }, // => { _id, batchKey, dataType, dataPerHash: Array<{ k: Hash, v: Data }> }
    { $unwind: { path: "$dataPerHash", preserveNullAndEmptyArrays: false } }, // => { _id, batchKey, dataPerHash: { k: DataType, v: ParHash<Data> } }, // => { _id, batchKey, dataType, dataPerHash: { k: Hash, v: Data } }
    { $project: { _id: 1, batchKey: 1, dataType: 1, dataHash: "$dataPerHash.k", "dataObject": "$dataPerHash.v" } }, // => { _id, batchKey, dataType, dataHash, dataObject: Data }

    {
      $facet: {
        "valid": [
          {
            $match: {
              dataType: "delai",
              $jsonSchema: {
                bsonType: "object",
                properties: {
                  dataObject: {
  "bsonType": "object",
  "required": [
    "date_creation",
    "date_echeance",
    "duree_delai",
    "montant_echeancier"
  ],
  "properties": {
    "date_creation": { "bsonType": "date" },
    "date_echeance": { "bsonType": "date" },
    "duree_delai": {
      "bsonType": "number",
      "minimum": 1,
      "description": "doit valoir 1 ou plus"
    },
    "montant_echeancier": {
      "bsonType": "number",
      "minimum": 0.01,
      "description": "doit valoir plus que 0 euros"
    }
  }
}
                }
              }
            }
          },
          {
            $project: {
              _id: 0,
              dataType: 1,
              dataHash: 1
            }
          }
        ],
        "invalid": [
          {
            $match: {
              dataType: "delai",
              $nor: [
                {
                  $jsonSchema: {
                    bsonType: "object",
                    properties: {
                      dataObject: {
  "bsonType": "object",
  "required": [
    "date_creation",
    "date_echeance",
    "duree_delai",
    "montant_echeancier"
  ],
  "properties": {
    "date_creation": { "bsonType": "date" },
    "date_echeance": { "bsonType": "date" },
    "duree_delai": {
      "bsonType": "number",
      "minimum": 1,
      "description": "doit valoir 1 ou plus"
    },
    "montant_echeancier": {
      "bsonType": "number",
      "minimum": 0.01,
      "description": "doit valoir plus que 0 euros"
    }
  }
}
                    }
                  }
                }
              ]
            }
          },
          {
            $project: {
              _id: 0,
              dataType: 1,
              dataHash: 1
            }
          }
        ]
      }
    },
  ]).toArray()[0])
