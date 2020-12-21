#!/bin/bash

# Test de bout en bout du CLI et de la documentation de ses commandes.
# Ce script doit être exécuté depuis la racine du projet. Ex: par test-all.sh.

# Setup
FLAGS="$*" # the script will update the golden file if "--update" flag was provided as 1st argument
TMP_DIR="tests/tmp-test-execution-files"
OUTPUT_FILE="${TMP_DIR}/test.output.txt"
GOLDEN_FILE="tests/output-snapshots/test-cli.golden.txt"
mkdir -p "${TMP_DIR}"

echo "" > "${OUTPUT_FILE}"

function test {
  CMD="./$1"
  echo "- ${CMD}"
  echo "$ ${CMD}"  >> "${OUTPUT_FILE}"
  echo "$(${CMD})" >> "${OUTPUT_FILE}"
  echo "---"       >> "${OUTPUT_FILE}"
}

# run test cases
test "sfdata"                           2>> "${OUTPUT_FILE}"
test "sfdata unknown_command"           2>> "${OUTPUT_FILE}"
test "sfdata --help"                    2>> "${OUTPUT_FILE}"
test "sfdata purge"                     2>> "${OUTPUT_FILE}"
test "sfdata purge --help"              2>> "${OUTPUT_FILE}"
test "sfdata check --help"              2>> "${OUTPUT_FILE}"
test "sfdata pruneEntities --help"      2>> "${OUTPUT_FILE}"
test "sfdata import --help"             2>> "${OUTPUT_FILE}"
test "sfdata validate --help"           2>> "${OUTPUT_FILE}"
test "sfdata compact --help"            2>> "${OUTPUT_FILE}"
test "sfdata reduce --help"             2>> "${OUTPUT_FILE}"
test "sfdata public --help"             2>> "${OUTPUT_FILE}"
test "sfdata etablissements --help"     2>> "${OUTPUT_FILE}"
test "sfdata entreprises --help"        2>> "${OUTPUT_FILE}"
test "sfdata purgeNotCompacted --help"  2>> "${OUTPUT_FILE}"

set -e # will stop the script if any command fails with a non-zero exit code
tests/helpers/diff-or-update-golden-master.sh "${FLAGS}" "${GOLDEN_FILE}" "${OUTPUT_FILE}"

rm -rf "${TMP_DIR}"
