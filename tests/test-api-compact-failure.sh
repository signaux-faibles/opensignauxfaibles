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
        "key" : "2008",
        "type" : "batch"
      }
    }
  ])

  db.ImportedData.insertMany([
    {
      "value": {
        "key": "01234567891011",
        "scope": "etablissement",
        "batch": {
          "2009": {
            "cotisation": undefined
          }
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
        "batch" : { "2007": {} }
      }
    }
  ])
CONTENTS

echo ""
echo "💎 Compacting RawData thru dbmongo API..."
tests/helpers/dbmongo-server.sh start

OUTPUT_GZ_FILE=dbmongo/$(http --print=b --ignore-stdin :5000/api/data/validate collection=RawData | tr -d '"')
echo "- POST /api/data/validate RawData 👉 ${OUTPUT_GZ_FILE}"
diff <(echo -n '') <(zcat < "${OUTPUT_GZ_FILE}")

OUTPUT_GZ_FILE=dbmongo/$(http --print=b --ignore-stdin :5000/api/data/validate collection=ImportedData | tr -d '"')
echo "- POST /api/data/validate ImportedData 👉 ${OUTPUT_GZ_FILE}, contents:"
diff <(echo '(invalid data entry)') <(zcat < "${OUTPUT_GZ_FILE}") # we expect an invalid data entry to be listed

echo "- POST /api/data/compact => diff:"
diff <(echo -n '"ok"') <(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=2008) # will fail with "TypeError: can't convert undefined to object"
echo "✅ No diff => OK"

# Now, the "trap" commands will clean up the rest.
