#!/bin/bash

# Test de bout en bout des commandes "entreprise" et "etablissement".
# Inspiré de test.sh.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

tests/helpers/mongodb-container.sh stop

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
COLOR_YELLOW='\033[1;33m'
COLOR_DEFAULT='\033[0m'
ETAB_GOLDEN_FILE="tests/output-snapshots/test-export-etablissements.golden.json"
ENTR_GOLDEN_FILE="tests/output-snapshots/test-export-entreprises.golden.json"
TMP_DIR="tests/tmp-test-execution-files"
mkdir -p "${TMP_DIR}"

# Clean up on exit
function teardown {
    echo -e "${COLOR_DEFAULT}"
    tests/helpers/sfdata-wrapper.sh stop || true # keep tearing down, even if "No matching processes belonging to you were found"
    tests/helpers/mongodb-container.sh stop
}
trap teardown EXIT

PORT="27016" tests/helpers/mongodb-container.sh start

MONGODB_PORT="27016" tests/helpers/sfdata-wrapper.sh setup

echo ""
echo "📝 Inserting test data..."
sleep 1 # give some time for MongoDB to start
tests/helpers/mongodb-container.sh run > /dev/null << CONTENTS
  db.Admin.insertOne({
    "_id" : {
        "key" : "2002_1",
        "type" : "batch"
    },
    "param" : {
        "date_debut" : ISODate("2014-01-01T00:00:00.000+0000"),
        "date_fin" : ISODate("2014-03-01T00:00:00.000+0000")
    }
  })

  // The random order of documents is intentional, to make sure that the output is correctly sorted no matter what
  db.ImportedData.insertMany([
    {
        "_id": "etab2",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "etablissement",
            "key": "01234567891012"
        }
    },
    {
        "_id": "etab21",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "etablissement",
            "key": "21234567891011"
        }
    },
    {
        "_id": "entr2",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "entreprise",
            "key": "212345678"
        }
    },
    {
        "_id": "entr1",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "entreprise",
            "key": "012345678"
        }
    },
    {
        "_id": "etab1",
        "value": {
            "batch": {
                "2002_1": {}
            },
            "scope": "etablissement",
            "key": "01234567891011"
        }
    },
  ])

  db.Scores.insertMany([
    {
        "_id": "score1",
        "siret" : "01234567891011",
        "periode" : "2014-01-01",
        "score" : 0.97,
        "batch" : "2002_1",
        "timestamp" : ISODate("2014-01-01T00:00:00.000+0000"),
        "algo" : "algo_avec_urssaf",
        "alert" : "Alerte seuil F1"
    },
    {
        "_id": "score2",
        "siret" : "01234567891011",
        "periode" : "2014-02-01",
        "score" : 0.98,
        "batch" : "2002_1",
        "timestamp" : ISODate("2014-02-01T00:00:00.000+0000"),
        "algo" : "algo_avec_urssaf",
        "alert" : "Alerte seuil F1"
    },
  ])
CONTENTS

echo ""
echo "💎 Computing the Public collection..."
echo "- sfdata compact 👉 $(tests/helpers/sfdata-wrapper.sh run compact --since-batch=2002_1)"
echo "- sfdata public 👉 $(tests/helpers/sfdata-wrapper.sh run public --until-batch=2002_1 --key=.........)" # TODO: we specify a placeholder value as key, so that PublicOne() is run instead of Public(), so the data is generated for etablissements that don't have effectif values, and therefore are outside of the "algo2" scope.

echo ""
echo "🚚 Export enterprise data..."
# This step is required only if key was provided when calling sfdata public
RENAME_RESULT=$(tests/helpers/mongodb-container.sh run <<< 'db.Public_debug.renameCollection("Public");')
echo "- rename 'Public_debug' collection to 'Public' 👉 ${RENAME_RESULT}"
# Make sure that the export only relies on Score and Public collections => clear collections that were populated for/by other endpoints
CLEAN_RESULT=$(tests/helpers/mongodb-container.sh run <<< 'db.Admin.drop(); db.ImportedData.drop(); db.RawData.drop();')
echo "- drop other db collections 👉 ${CLEAN_RESULT}"

function stopIfFailed {
    if [[ "$1" == *failed* ]]
    then
        exit 1
    fi
}

# Parameter validation
RESULT=$(tests/helpers/sfdata-wrapper.sh run etablissements --key="invalid" | (grep "key doit être un numéro SIREN" || echo -e "${COLOR_YELLOW}failed${COLOR_DEFAULT}"))
echo "- sfdata etablissements with invalid key 👉 ${RESULT}"
stopIfFailed "${RESULT}"
RESULT=$(tests/helpers/sfdata-wrapper.sh run entreprises --key="invalid" | (grep "key doit être un numéro SIREN" || echo -e "${COLOR_YELLOW}failed${COLOR_DEFAULT}"))
echo "- sfdata entreprises with invalid key 👉 ${RESULT}"
stopIfFailed "${RESULT}"

# sfdata etablissements with key=212345678 should return just one match
RESULTS=$(tests/helpers/sfdata-wrapper.sh run etablissements --key="212345678")
MATCH=$(echo "${RESULTS}" | grep --quiet "etablissement_21234567891011" && echo "found etablissement_21234567891011" || echo -e "${COLOR_YELLOW}failed${COLOR_DEFAULT}")
COUNT=$(echo "${RESULTS}" | wc -l)
echo "- sfdata etablissements with key=212345678 👉 ${MATCH}, ${COUNT} result(s)"
stopIfFailed "${MATCH}"
if [[ "${COUNT}" -ne "1" ]]
then
    exit 1
fi

# sfdata entreprises with key=212345678 should return just one match
RESULTS=$(tests/helpers/sfdata-wrapper.sh run entreprises --key="212345678")
MATCH=$(echo "${RESULTS}" | grep --quiet "entreprise_212345678" && echo "found entreprise_212345678" || echo -e "${COLOR_YELLOW}failed${COLOR_DEFAULT}")
COUNT=$(echo "${RESULTS}" | wc -l)
echo "- sfdata entreprises with key=212345678 👉 ${MATCH}, ${COUNT} result(s)"
stopIfFailed "${MATCH}"
if [[ "${COUNT}" -ne "1" ]]
then
    exit 1
fi

# Export enterprise data
RESULTS=$(tests/helpers/sfdata-wrapper.sh run etablissements)
ETABLISSEMENTS_FILE="${TMP_DIR}/etablissements.json"
echo "${RESULTS}" > "${ETABLISSEMENTS_FILE}"
echo "- sfdata etablissements 👉 ${ETABLISSEMENTS_FILE}"
RESULTS=$(tests/helpers/sfdata-wrapper.sh run entreprises)
ENTREPRISES_FILE="${TMP_DIR}/entreprises.json"
echo "${RESULTS}" > "${ENTREPRISES_FILE}"
echo "- sfdata entreprises 👉 ${ENTREPRISES_FILE}"

tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${ETAB_GOLDEN_FILE}" "${ETABLISSEMENTS_FILE}"
tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${ENTR_GOLDEN_FILE}" "${ENTREPRISES_FILE}"

rm "${ETABLISSEMENTS_FILE}" "${ENTREPRISES_FILE}"
rm -rf "${TMP_DIR}"
# Now, the "trap" commands will clean up the rest.
