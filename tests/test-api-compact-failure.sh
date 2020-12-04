#!/bin/bash

# Test de bout en bout de l'API "compact".
# Test de non regression pour https://github.com/signaux-faibles/opensignauxfaibles/issues/248.
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
echo "💎 Compacting RawData thru dbmongo API..."
tests/helpers/dbmongo-server.sh start

RAWDATA_ERRORS_FILE=dbmongo/$(http --print=b --ignore-stdin :5000/api/data/validate collection=RawData | tr -d '"')
echo "- POST /api/data/validate RawData 👉 ${RAWDATA_ERRORS_FILE}"
diff <(echo -n '') <(zcat < "${RAWDATA_ERRORS_FILE}") # no validation errors detected in RawData
rm "${RAWDATA_ERRORS_FILE}"

IMPORTEDDATA_ERRORS_FILE=dbmongo/$(http --print=b --ignore-stdin :5000/api/data/validate collection=ImportedData | tr -d '"')
echo "- POST /api/data/validate ImportedData 👉 ${IMPORTEDDATA_ERRORS_FILE}"
diff <(echo -n '') <(zcat < "${IMPORTEDDATA_ERRORS_FILE}") # no validation errors detected in ImportedData
rm "${IMPORTEDDATA_ERRORS_FILE}"

echo "- POST /api/data/compact should not fail"
RESULT=$(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=2011_0_urssaf)
echo "${RESULT}"
echo "${RESULT}" | grep --quiet "ok"

echo "✅ OK"

# Now, the "trap" commands will clean up tmp files.
