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
    {"_id":{"key":"2011_0_urssaf","type":"batch"},"complete_types":["cotisation","debit","delai","procol","effectif","effectif_ent"],"files":{"admin_urssaf":["/2011/2011_0_urssaf/b74a8b78a26bfddcf45f2b175a069838"],"ccsf":["/2011/2011_0_urssaf/b8e54e397c20b64c179954f8a3b8f266"],"cotisation":["/2011/2011_0_urssaf/f15a90dc85b2114bf9d6fba2de985af4"],"debit":["/2011/2011_0_urssaf/3aef1ea58de1f89794abcaf814581a3f"],"delai":["/2011/2011_0_urssaf/72a63e0fe39ac8b81bb399d55ce9d76b"],"effectif":["/2011/2011_0_urssaf/d91272eeae631db500813c248c0fdf84"],"effectif_ent":["/2011/2011_0_urssaf/da88c6d026396ef551997a4a5c1bd446"],"filter":["/2011/2011_0_urssaf/filter_siren_2011_urssaf.csv"],"procol":["/2011/2011_0_urssaf/c4093239485db05bd0a4814a91ff3432"]},"param":{"date_debut":{"$date":{"$numberLong":"1388534400000"}},"date_fin":{"$date":{"$numberLong":"1604188800000"}},"date_fin_effectif":{"$date":{"$numberLong":"1598918400000"}}}},
    {"_id":{"key":"2011_1_sirene","type":"batch"},"complete_types":["sirene","sirene_ul"],"files":{"filter":["/2011/2011_1_sirene/filter_siren_2011_urssaf.csv"],"sirene":["/2011/2011_1_sirene/StockEtablissement_utf8_geo.csv"],"sirene_ul":["/2011/2011_1_sirene/sireneUL.csv"]},"param":{"date_debut":{"$date":{"$numberLong":"1388534400000"}},"date_fin":{"$date":{"$numberLong":"1604188800000"}},"date_fin_effectif":{"$date":{"$numberLong":"1598918400000"}}}},
    {"_id":{"key":"2011_2_activitepartielle","type":"batch"},"complete_types":["apconso","apdemande"],"files":{"apconso":["/2011/2011_2_activitepartielle/296a8fbc94c3857a22d31f3929e413dd"],"apdemande":["/2011/2011_2_activitepartielle/e69ed899dbaa70d3995986905905aa0b"],"filter":["/2011/filter_siren_2011_urssaf.csv"]},"param":{"date_debut":{"$date":{"$numberLong":"1388534400000"}},"date_fin":{"$date":{"$numberLong":"1604188800000"}},"date_fin_effectif":{"$date":{"$numberLong":"1598918400000"}}}},
  ])

  db.ImportedData.insertMany([
    {
      "_id": {
        "$oid": "5fc4fe82ef6c0e7f34db3925"
      },
      "value": {
        "scope": "entreprise",
        "key": "000000000",
        "batch": {
          "2011_0_urssaf": {
            "effectif_ent": {}
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

# RAWDATA_ERRORS_FILE=dbmongo/$(http --print=b --ignore-stdin :5000/api/data/validate collection=RawData | tr -d '"')
# echo "- POST /api/data/validate RawData 👉 ${RAWDATA_ERRORS_FILE}"
# diff <(echo -n '') <(zcat < "${RAWDATA_ERRORS_FILE}")

# IMPORTEDDATA_ERRORS_FILE=dbmongo/$(http --print=b --ignore-stdin :5000/api/data/validate collection=ImportedData | tr -d '"')
# echo "- POST /api/data/validate ImportedData 👉 ${IMPORTEDDATA_ERRORS_FILE}"
# grep --quiet '{"_id":"5f9192703029a1f7d4b1773b","batchKey":"2009","dataPerHash":{},"dataType":"cotisation"}' <(zcat < "${IMPORTEDDATA_ERRORS_FILE}") # we expect an invalid data entry to be listed

echo "- POST /api/data/compact should not fail"
RESULT=$(http --print=b --ignore-stdin :5000/api/data/compact fromBatchKey=2011_0_urssaf)
echo "${RESULT}"
echo "${RESULT}" | grep "ok"

echo "✅ OK"

# rm "${RAWDATA_ERRORS_FILE}" "${IMPORTEDDATA_ERRORS_FILE}"
# Now, the "trap" commands will clean up the rest.
