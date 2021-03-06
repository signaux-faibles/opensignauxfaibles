#!/bin/bash

# Test de bout en bout de la commande "compact".
# Test de non regression pour https://github.com/signaux-faibles/opensignauxfaibles/issues/248.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

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

tests/helpers/mongodb-container.sh run << CONTENTS
  db.Admin.insertMany([
    {"_id":{"key":"2011_0_urssaf","type":"batch"}},
    {"_id":{"key":"2011_1_sirene","type":"batch"}},
  ])

  db.ImportedData.insertMany([
    {
      "value": {
        "scope": "entreprise",
        "key": "000000000",
        "batch": {
          "2011_0_urssaf": {
          }
        }
      }
    }
  ])

  db.RawData.insertMany([
    {
      "_id": "000000000",
      "value": {
        "key": "000000000",
        "scope": "entreprise",
        "batch": {
    			"1910_6": {}
        }
      }
    },
  ])
CONTENTS

echo ""
echo "💎 Compacting RawData..."

VALIDATION_REPORT=$(tests/helpers/sfdata-wrapper.sh validate --collection=RawData)
echo "- sfdata validate RawData"
diff <(echo '') <(echo "${VALIDATION_REPORT}") # no validation errors detected in RawData

VALIDATION_REPORT=$(tests/helpers/sfdata-wrapper.sh validate --collection=ImportedData)
echo "- sfdata validate ImportedData"
diff <(echo '') <(echo "${VALIDATION_REPORT}") # no validation errors detected in ImportedData

echo "- sfdata compact should not fail"
RESULT=$(tests/helpers/sfdata-wrapper.sh compact --since-batch=2011_0_urssaf)
echo "${RESULT}" | grep --quiet "ok"

echo "✅ OK"

# Now, the "trap" commands will clean up tmp files.
