#!/bin/bash

# Test de bout en bout de l'API "compact".
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

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

tests/helpers/mongodb-container.sh run << CONTENTS
  db.Admin.insertMany([
    {
      "_id" : {
        "key" : "2009",
        "type" : "batch"
      },
      "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2016-01-01T00:00:00.000+0000"),
        "date_fin_effectif" : ISODate("2016-03-01T00:00:00.000+0000")
      },
      "name" : "TestData"
    },
    {
      "_id" : {
        "key" : "2010",
        "type" : "batch"
      },
      "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2016-01-01T00:00:00.000+0000"),
        "date_fin_effectif" : ISODate("2016-03-01T00:00:00.000+0000")
      },
      "name" : "TestData"
    }
  ])

  db.ImportedData.insertMany([
    {
      "_id": "import1",
      "value": {
        "key": "01234567891011",
        "scope": "etablissement",
        "batch": {
          "2009": {
            "cotisation": undefined
          }
        },
        "index": {
          "algo2": true
        }
      }
    }
  ])

  db.RawData.insertMany([
    { 
      "_id" : "01234567891011", 
      "value" : {
        "key" : "01234567891011", 
        "scope" : "etablissement", 
        "batch" : {}
      }
    }
  ])
CONTENTS

echo ""
echo "💎 Compacting RawData thru dbmongo API..."
tests/helpers/dbmongo-server.sh start
echo "- POST /api/data/compact => diff:"

diff <(echo "ok") <(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=2009) # will fail with "TypeError: can't convert undefined to object"

# Now, the "trap" commands will clean up the rest.
