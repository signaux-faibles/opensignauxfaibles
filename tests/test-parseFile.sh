#!/bin/bash

# Test de bout en bout de la commande "parseFile".
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

set -e # will stop the script if any command fails with a non-zero exit code

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test-parseFile.output.txt"
GOLDEN_FILE="tests/output-snapshots/test-parseFile.golden.txt"
mkdir -p "${TMP_DIR}"

echo ""
echo "💎 Parsing data..."
echo "- sfdata parseFile ..."
NO_DB=1 tests/helpers/sfdata-wrapper.sh parseFile --parser "diane" --file "lib/diane/testData/dianeTestData.txt" > "${OUTPUT_FILE}"
tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
