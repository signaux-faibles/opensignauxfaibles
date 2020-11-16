/* Generated by bundle-schemas.go - DO NOT EDIT */

package engine

var validationSchemas = map[string]string{
"bdf.schema.json": `{
  "title": "EntréeBdf",
  "description": "Note: CE SCHEMA EST INCOMPLET POUR L'INSTANT. Cf https://github.com/signaux-faibles/opensignauxfaibles/pull/143",
  "bsonType": "object",
  "required": ["siren"],
  "properties": {
    "siren": {
      "bsonType": "string",
      "pattern": "^[0-9]{9}$"
    }
  },
  "additionalProperties": false
}
`,
"delai.schema.json": `{
  "title": "EntréeDelai",
  "description": "Note: CE SCHEMA EST INCOMPLET POUR L'INSTANT. Cf https://github.com/signaux-faibles/opensignauxfaibles/pull/143",
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
      "description": "Doit valoir 1 ou plus:",
      "minimum": 1
    },
    "montant_echeancier": {
      "bsonType": "number",
      "description": "Doit valoir plus que 0 euros:",
      "minimum": 0.01
    }
  },
  "additionalProperties": false
}
`,
"detect_invalid_entries.pipeline.json": `[
  { "$project": { "_id": 1, "batches": { "$objectToArray": "$value.batch" } } },
  { "$unwind": { "path": "$batches", "preserveNullAndEmptyArrays": false } },
  {
    "$project": {
      "_id": 1,
      "batchKey": "$batches.k",
      "dataPerHash": { "$objectToArray": "$batches.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": "$dataPerHash.k",
      "dataPerHash": "$dataPerHash.v"
    }
  },
  { "$match": { "dataPerHash": { "$not": { "$type": 3 } } } }
]
`,
"flatten_data_entries.pipeline.json": `[
  { "$project": { "_id": 1, "batches": { "$objectToArray": "$value.batch" } } },
  { "$unwind": { "path": "$batches", "preserveNullAndEmptyArrays": false } },
  {
    "$project": {
      "_id": 1,
      "batchKey": "$batches.k",
      "dataPerHash": { "$objectToArray": "$batches.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": "$dataPerHash.k",
      "dataPerHash": { "$objectToArray": "$dataPerHash.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": 1,
      "dataHash": "$dataPerHash.k",
      "dataObject": "$dataPerHash.v"
    }
  }
]
`,
}
